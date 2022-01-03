package main

import (
	"log"
	"math/rand"
	"net"
	"os"
	"os/signal"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	"github.com/S-Masakatsu/hello-grpc/bake-pancake/rpc"
	"github.com/S-Masakatsu/hello-grpc/bake-pancake/server/handler"
)

func main() {
	rand.Seed(time.Now().UnixNano())

	addr := ":50051"
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()
	rpc.RegisterPancakeBakerServiceServer(
		s, handler.NewBakerHandler(),
	)
	reflection.Register(s)

	go func() {
		log.Printf("start gRPC server listen: %s", addr)
		s.Serve(lis)
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit
	log.Println("stopping gRPC server...")
	s.GracefulStop()
}
