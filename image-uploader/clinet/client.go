package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"google.golang.org/grpc"

	"github.com/S-Masakatsu/hello-grpc/image-uploader/rpc"
)

func main() {
	ctx := context.Background()
	path := os.Args[1]

	cli := NewClient(ctx)
	res, err := cli.Upload(ctx, path)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(res)
}

type ImageUpload interface {
	Upload(context.Context, string) (*rpc.ImageUploadResponse, error)
}

type ClientImpl struct {
	client rpc.ImageUploadServiceClient
}

func NewClient(ctx context.Context) ImageUpload {
	cli, err := NewImageUploadClient(ctx, "localhost", 50051, 3*time.Second)
	if err != nil {
		log.Fatal(err)
	}
	return &ClientImpl{
		client: cli,
	}
}

func (i *ClientImpl) Upload(ctx context.Context, path string) (*rpc.ImageUploadResponse, error) {
	filename := strings.Replace(filepath.Base(path), filepath.Ext(path), "", -1)
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	buf := make([]byte, 1024)

	// 最初のリクエスト
	fileMeta := &rpc.ImageUploadRequest_FileMeta{
		Filename: filename,
	}
	stream, err := i.client.Upload(ctx)
	if err != nil {
		return nil, err
	}

	stream.Send(&rpc.ImageUploadRequest{File: &rpc.ImageUploadRequest_FileMeta_{
		FileMeta: fileMeta,
	}})

	// 画像のアップロード
	for {
		if _, err = f.Read(buf); err == io.EOF {
			break
		} else if err != nil {
			return nil, err
		}

		err = stream.Send(&rpc.ImageUploadRequest{File: &rpc.ImageUploadRequest_Data{
			Data: buf,
		}})
		if err != nil {
			return nil, err
		}
	}
	return stream.CloseAndRecv()
}

func NewImageUploadClient(ctx context.Context, host string, port int, d time.Duration) (rpc.ImageUploadServiceClient, error) {
	ctx, cancel := context.WithTimeout(ctx, d)
	defer cancel()

	opts := []grpc.DialOption{
		grpc.WithInsecure(),
	}
	conn, err := grpc.DialContext(ctx, fmt.Sprintf("%s:%d", host, port), opts...)
	if err != nil {
		return nil, err
	}
	return rpc.NewImageUploadServiceClient(conn), nil
}
