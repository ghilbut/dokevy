package main

import (
	"context"
	"flag"
	"fmt"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/types/known/emptypb"
	"log"
	"net"

	// external packages
	"google.golang.org/grpc"

	// project packages
	pb "github.com/ghilbut/dokevy/grpc/v1"
)

type SystemServiceServer struct {
	pb.UnimplementedSystemServiceServer
}

func (s *SystemServiceServer) Ping(context.Context, *emptypb.Empty) (*pb.Pong, error) {
	return &pb.Pong{Pong: "OK"}, nil
}

func main() {
	fmt.Println("Hello, World")

	flag.Parse()
	lis, err := net.Listen("tcp", fmt.Sprintf("localhost:50051"))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	creds := insecure.NewCredentials()
	opts := []grpc.ServerOption{grpc.Creds(creds)}
	grpcServer := grpc.NewServer(opts...)
	systemServer := SystemServiceServer{}
	pb.RegisterSystemServiceServer(grpcServer, &systemServer)
	grpcServer.Serve(lis)
}
