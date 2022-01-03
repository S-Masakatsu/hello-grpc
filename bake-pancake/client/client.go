package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"google.golang.org/grpc"

	"github.com/S-Masakatsu/hello-grpc/bake-pancake/rpc"
)

func main() {
	ctx := context.Background()
	client, err := NewBakePancakeClient(ctx, "localhost", 50051, 5*time.Second)
	if err != nil {
		log.Fatal(err)
	}

	var input string
	for {
		fmt.Print("input: ")
		fmt.Scan(&input)

		switch input {
		case "report":
			res, err := client.Report(ctx, &rpc.ReportRequest{})
			if err != nil {
				log.Fatal(err)
			}
			fmt.Println(res.Report.BakeCounts)
		case "exit":
			os.Exit(0)
		default:
			n, err := strconv.Atoi(input)
			if err != nil {
				fmt.Println(err)
				continue
			}
			req := &rpc.BakeRequest{Menu: rpc.Pancake_Menu(n)}
			res, err := client.Bake(ctx, req)
			if err != nil {
				log.Fatal(err)
			}
			fmt.Println(res)
		}
	}
}

func NewBakePancakeClient(ctx context.Context, host string, port int, d time.Duration) (rpc.PancakeBakerServiceClient, error) {
	ctx, cancel := context.WithTimeout(ctx, d)
	defer cancel()

	opts := []grpc.DialOption{
		grpc.WithInsecure(),
	}
	conn, err := grpc.DialContext(ctx, fmt.Sprintf("%s:%d", host, port), opts...)
	if err != nil {
		return nil, err
	}
	return rpc.NewPancakeBakerServiceClient(conn), nil
}
