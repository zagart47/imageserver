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
	"io"
	"log"
	"net"
	"os"
)

func Start() {
	listener1, err := net.Listen("tcp", "0.0.0.0:80")
	if err != nil {
		log.Fatal(err)
	}
	defer listener1.Close()
	grpcServer1 := grpc.NewServer(grpc.MaxConcurrentStreams(10))
	pb.RegisterFileServiceServer(grpcServer1, &Server{})
	defer grpcServer1.GracefulStop()

	listener2, err := net.Listen("tcp", "0.0.0.0:81")
	if err != nil {
		log.Fatal(err)
	}
	defer listener2.Close()
	grpcServer2 := grpc.NewServer(grpc.MaxConcurrentStreams(100))
	pb.RegisterListServiceServer(grpcServer2, &ListServer{})
	defer grpcServer2.GracefulStop()
	go func() {
		log.Fatal(grpcServer1.Serve(listener1))
	}()
	log.Fatal(grpcServer2.Serve(listener2))

}

type Server struct {
	pb.UnimplementedFileServiceServer
}

type ListServer struct {
	pb.UnimplementedListServiceServer
}

func (s Server) Download(request *pb.DownloadRequest, stream pb.FileService_DownloadServer) error {
	md, ok := metadata.FromIncomingContext(stream.Context())
	if !ok {
		return fmt.Errorf("md incoming error")
	}
	fileName := md.Get("filename")[0]
	f := storage.NewFile(fileName)
	repo := db.NewSQLiteRepository()
	if err := repo.CheckFileName(f.Name); errors.Is(err, db.ErrFileNotFound) {
		return err
	}
	file, err := os.Open(f.Path)
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			log.Println(err.Error())
		}
	}(file)
	if err != nil {
		return err
	}
	buff := make([]byte, 1024*1024)

	for {
		n, err := file.Read(buff)
		if err == io.EOF {
			break
		}
		if err != nil {
			return fmt.Errorf("buffer reading error (%s)", err.Error())
		}
		req := &pb.DownloadResponse{Fragment: buff[:n]}
		if err := stream.Send(req); err != nil {
			log.Fatal(err.Error())
		}
	}
	return nil
}

func (ls ListServer) GetFiles(context.Context, *pb.GetFilesRequest) (*pb.GetFilesResponse, error) {
	result, err := db.DownloadFileList()
	return &pb.GetFilesResponse{Info: result}, err
}

func (s Server) Upload(stream pb.FileService_UploadServer) error {
	md, ok := metadata.FromIncomingContext(stream.Context())
	if !ok {
		return fmt.Errorf("md incoming error")
	}
	fileName := md.Get("filename")[0]
	f := storage.NewFile(fileName)

	repo := db.NewSQLiteRepository()

	for {
		req, err := stream.Recv()
		if err == io.EOF {
			if err := os.WriteFile(f.Path, f.Buffer.Bytes(), 0644); err != nil {
				return err
			}
			if err := repo.CheckFileName(f.Name); errors.Is(err, db.ErrFileFound) {
				if err := repo.Update(f.Name); err != nil {
					return err
				}
			} else if errors.Is(err, db.ErrFileNotFound) {
				if err := repo.Create(f.Name); err != nil {
					return err
				}
			}
			return stream.SendAndClose(&pb.UploadResponse{})
		}
		f.Buffer.Write(req.GetFragment())
		if err != nil {
			return status.Error(codes.Internal, err.Error())
		}
	}
}
