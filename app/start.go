package app

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"imageserver/db"
	"imageserver/model"
	pb "imageserver/pkg/proto"
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

func (s Server) Download(request *pb.DownloadRequest, server pb.FileService_DownloadServer) error {
	name := request.GetFilename()
	bufferSize := 64 * 1024
	file, err := os.Open("files/" + name)
	if err != nil {
		fmt.Println(err)
		return err
	}
	defer file.Close()
	resp := &pb.DownloadResponse{Filename: name}
	err = server.Send(resp)
	if err != nil {
		return err
	}
	buff := make([]byte, bufferSize)
	for {
		bytesRead, err := file.Read(buff)
		if err != nil {
			if err != io.EOF {
				fmt.Println(err)
			}
			break
		}
		resp = &pb.DownloadResponse{
			Fragment: buff[:bytesRead],
		}
		err = server.Send(resp)
		if err != nil {
			return err
		}
	}
	return nil
}

func (ls ListServer) GetFiles(context.Context, *pb.GetFilesRequest) (*pb.GetFilesResponse, error) {
	fileRepository := db.NewSQLiteRepository()
	all, err := fileRepository.All()
	if err != nil {
		return nil, err
	}
	return &pb.GetFilesResponse{Info: all}, nil
}

func (s Server) Upload(stream pb.FileService_UploadServer) error {
	md, ok := metadata.FromIncomingContext(stream.Context())
	if !ok {
		return fmt.Errorf("md incoming error")
	}
	fileName := md.Get("filename")[0]
	f := model.NewFile(fileName)

	repo := db.NewSQLiteRepository()

	for {
		req, err := stream.Recv()
		if err == io.EOF {
			err := os.WriteFile(f.Path, f.Buffer.Bytes(), 0644)
			if err != nil {
				return err
			}
			err = repo.CheckFileName(f.Name)
			if err != nil {
				return err
			}
			return stream.SendAndClose(&pb.UploadResponse{Name: f.Name})
		}
		f.Buffer.Write(req.GetFragment())
		if err != nil {
			return status.Error(codes.Internal, err.Error())
		}
	}
}
