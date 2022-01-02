package main

import (
	"context"
	"log"
	"os"

	"google.golang.org/grpc"

	"github.com/S-Masakatsu/hello-grpc/rpc"
)

func main() {
	addr := "localhost:50051"
	ctx := context.Background()

	conn, err := grpc.DialContext(ctx, addr, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	cli := rpc.NewGreeterClient(conn)

	name := os.Args[1]

	res, err := cli.SayHello(ctx, &rpc.HelloRequest{Name: name})
	if err != nil {
		log.Fatalf("could not greet: %v", err)
	}
	log.Printf("Greeting: %s", res.Message)
}
