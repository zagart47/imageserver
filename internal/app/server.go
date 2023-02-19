package app

import (
	"context"
	"errors"
	"github.com/joho/godotenv"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"imageserver/internal/myerror"
	"imageserver/internal/repository"
	"imageserver/internal/storage"
	"imageserver/internal/table"
	pb "imageserver/proto"
	"io"
	"log"
	"net"
	"os"
)

func Run() {
	err := godotenv.Load()
	fileHost := os.Getenv("fileHost")
	fileListener, err := net.Listen("tcp", fileHost)
	if err != nil {
		log.Fatal(err)
	}
	defer func(fileListener net.Listener) {
		if err := fileListener.Close(); err != nil {
			log.Fatal(err)
		}
	}(fileListener)
	fileServer := grpc.NewServer(grpc.MaxConcurrentStreams(10))
	pb.RegisterFileServiceServer(fileServer, &Server{})
	defer fileServer.GracefulStop()

	listHost := os.Getenv("listHost")
	listListener, err := net.Listen("tcp", listHost)
	if err != nil {
		log.Fatal(err)
	}
	defer func(listListener net.Listener) {
		if err := listListener.Close(); err != nil {
			log.Fatal(err)
		}
	}(listListener)
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

var repo = repository.NewSQLiteRepository(repository.DB)

func (s Server) Download(request *pb.DownloadRequest, stream pb.FileService_DownloadServer) error {
	md, ok := metadata.FromIncomingContext(stream.Context())
	if !ok {
		return myerror.Err.MetaData
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
	buffer := make([]byte, 64*1024)

	for {
		n, err := open.Read(buffer)
		if err == io.EOF {
			break
		}
		if err != nil {
			return myerror.Err.Buffer
		}
		resp := &pb.DownloadResponse{Fragment: buffer[:n]}
		if err := stream.Send(resp); err != nil {
			return status.Error(codes.Internal, err.Error())
		}
	}
	return nil
}

func (ls ListServer) GetFiles(context.Context, *pb.GetFilesRequest) (*pb.GetFilesResponse, error) {
	fl, err := repo.ShowAllRecords()
	if err != nil {
		return nil, err
	}
	result := table.MakeTable(fl)
	return &pb.GetFilesResponse{Info: result}, err
}

func (s Server) Upload(stream pb.FileService_UploadServer) error {
	md, ok := metadata.FromIncomingContext(stream.Context())
	if !ok {
		return myerror.Err.MetaData
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
			} else if errors.Is(err, myerror.Err.FileNotFound) {
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
