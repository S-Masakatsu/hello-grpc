package main

import (
	"context"
	"log"
	"net"

	"google.golang.org/grpc"

	"github.com/S-Masakatsu/hello-grpc/rpc"
)

type server struct{}

func (s *server) SayHello(ctx context.Context, in *rpc.HelloRequest) (*rpc.HelloResponse, error) {
	log.Printf("Received: %s", in.Name)
	return &rpc.HelloResponse{Message: "Hello, " + in.Name}, nil
}

func main() {
	addr := ":50051"
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	rpc.RegisterGreeterServer(s, &server{})

	log.Printf("gRPC server listenting on " + addr)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to server: %v", err)
	}
}
