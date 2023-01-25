package app

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"imageserver/db"
	"imageserver/internal"
	"imageserver/model"
	pb "imageserver/pkg/proto"
	"io"
	"io/ioutil"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
)

func Start() {
	listener, err := net.Listen("tcp", "localhost:12223")
	if err != nil {
		log.Fatal(err)
	}
	defer listener.Close()
	sm := &internal.ServerMiddleware{}
	grpcServer := grpc.NewServer(grpc.UnaryInterceptor(sm.Interceptor))
	signalCh := make(chan os.Signal, 1)
	signal.Notify(signalCh, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-signalCh
		grpcServer.GracefulStop()
	}()
	pb.RegisterFileServiceServer(grpcServer, &Server{})
	log.Fatal(grpcServer.Serve(listener))
}

type Server struct {
	pb.UnimplementedFileServiceServer
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
			log.Println("error while sending chunk:", err)
			return err
		}
	}
	return nil
}

func (s Server) GetFiles(context.Context, *pb.GetFilesRequest) (*pb.GetFilesResponse, error) {
	fileRepository := db.NewSQLiteRepository(db.DB)
	all := fileRepository.All()
	return &pb.GetFilesResponse{Info: all}, nil
}

func (s Server) Upload(stream pb.FileService_UploadServer) error {
	fileRepository := db.NewSQLiteRepository(db.DB)
	nameData, _ := stream.Recv()
	name := model.File{FileName: nameData.GetFilename()}
	for {
		req, err := stream.Recv()
		if err == io.EOF {
			ioutil.WriteFile("files/"+nameData.GetFilename(), req.GetFragment(), 0644)
			err := fileRepository.CheckFileName(name.FileName)
			if err != nil {
				return err
			}
			return stream.SendAndClose(&pb.UploadResponse{Name: name.FileName})

		}
		if err != nil {
			return status.Error(codes.Internal, err.Error())
		}

	}

}
