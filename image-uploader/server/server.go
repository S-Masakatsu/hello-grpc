package main

import (
	"fmt"
	"log"
	"math/rand"
	"net"
	"os"
	"os/signal"
	"time"

	"github.com/S-Masakatsu/hello-grpc/image-uploader/rpc"
	"github.com/S-Masakatsu/hello-grpc/image-uploader/server/handler"
	grpc_zap "github.com/grpc-ecosystem/go-grpc-middleware/logging/zap"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	rand.Seed(time.Now().UnixNano())

	addr := ":50051"
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	logger, err := zap.NewProduction()
	if err != nil {
		log.Fatal(err)
	}
	grpc_zap.ReplaceGrpcLogger(logger)

	s := grpc.NewServer(
		grpc.UnaryInterceptor(
			grpc_zap.UnaryServerInterceptor(logger),
		),
	)

	rpc.RegisterImageUploadServiceServer(s, handler.NewImageUploadHandler())
	reflection.Register(s)

	go func() {
		fmt.Printf("start gRPC server addr: %s", addr)
		s.Serve(lis)
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit
	log.Println("stoping gRPC server...")
	s.GracefulStop()
}
