package app

import (
	"context"
	"errors"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"imageserver/db"
	pb "imageserver/pkg/proto"
	"imageserver/storage"
	"imageserver/table"
	"io"
	"log"
	"net"
	"os"
)

func Start() {
	fileListener, err := net.Listen("tcp", "0.0.0.0:80")
	if err != nil {
		log.Fatal(err)
	}
	defer fileListener.Close()
	fileServer := grpc.NewServer(grpc.MaxConcurrentStreams(10))
	pb.RegisterFileServiceServer(fileServer, &Server{})
	defer fileServer.GracefulStop()

	listListener, err := net.Listen("tcp", "0.0.0.0:81")
	if err != nil {
		log.Fatal(err)
	}
	defer listListener.Close()
	listServer := grpc.NewServer(grpc.MaxConcurrentStreams(100))
	pb.RegisterListServiceServer(listServer, &ListServer{})
	defer listServer.GracefulStop()
	go func() {
		log.Fatal(fileServer.Serve(fileListener))
	}()
	log.Fatal(listServer.Serve(listListener))

}

type Server struct {
	pb.UnimplementedFileServiceServer
}

type ListServer struct {
	pb.UnimplementedListServiceServer
}

var repo = db.NewSQLiteRepository(db.DB)

func (s Server) Download(request *pb.DownloadRequest, stream pb.FileService_DownloadServer) error {
	md, ok := metadata.FromIncomingContext(stream.Context())
	if !ok {
		return fmt.Errorf("md incoming error")
	}
	mdFileName := md.Get("filename")[0]
	file := storage.NewFile(mdFileName)
	if err := repo.CheckFileName(file.Name); err != nil {
		return err
	}
	open, err := os.Open(file.Path)
	if err != nil {
		return err
	}
	defer func(file *os.File) error {
		if err := file.Close(); err != nil {
			return err
		}
		return nil
	}(open)
	buffer := make([]byte, 1024*1024)

	for {
		n, err := open.Read(buffer)
		if err == io.EOF {
			break
		}
		if err != nil {
			return fmt.Errorf("buffer reading error (%s)", err.Error())
		}
		resp := &pb.DownloadResponse{Fragment: buffer[:n]}
		if err := stream.Send(resp); err != nil {
			return err
		}
	}
	return nil
}

func (ls ListServer) GetFiles(context.Context, *pb.GetFilesRequest) (*pb.GetFilesResponse, error) {
	fl, err := repo.DownloadFileList()
	result := table.MakeTable(fl)
	return &pb.GetFilesResponse{Info: result}, err
}

func (s Server) Upload(stream pb.FileService_UploadServer) error {
	md, ok := metadata.FromIncomingContext(stream.Context())
	if !ok {
		return fmt.Errorf("md incoming error")
	}
	file := storage.NewFile(md.Get("filename")[0])

	for {
		req, err := stream.Recv()
		if err == io.EOF {
			if err := os.WriteFile(file.Path, file.Buffer.Bytes(), 0644); err != nil {
				return err
			}
			if err := repo.CheckFileName(file.Name); err == nil {
				if err := repo.Update(file.Name); err != nil {
					return err
				}
			} else if errors.Is(err, db.ErrFileNotFound) {
				if err := repo.Create(file.Name); err != nil {
					return err
				}
			} else if err != nil {
				return err
			}
			return stream.SendAndClose(&pb.UploadResponse{})
		}
		file.Buffer.Write(req.GetFragment())
		if err != nil {
			return status.Error(codes.Internal, err.Error())
		}
	}
}
