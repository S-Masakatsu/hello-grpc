package handler

import (
	"bytes"
	"io"
	"log"
	"net/http"
	"sync"

	"github.com/google/uuid"

	"github.com/S-Masakatsu/hello-grpc/image-uploader/rpc"
)

type ImageUploadHandler struct {
	sync.Mutex
	files map[string][]byte
	rpc.UnsafeImageUploadServiceServer
}

func NewImageUploadHandler() rpc.ImageUploadServiceServer {
	return &ImageUploadHandler{
		files: make(map[string][]byte),
	}
}

func (h *ImageUploadHandler) Upload(stream rpc.ImageUploadService_UploadServer) error {
	// 最初のリクエストを受け取る
	req, err := stream.Recv()
	if err != nil {
		return err
	}

	// 初回はメタ情報が送られる仕様
	meta := req.GetFileMeta()
	log.Printf("get meta data")
	filename := meta.Filename

	// UUIDの生成
	u, err := uuid.NewRandom()
	if err != nil {
		return err
	}
	uuid := u.String()

	buf := &bytes.Buffer{}

	// アップロードされたバイナリをループしながら受け取る(チャンクドアップロード)
	for {
		// MEMO: Recv関数は全てのリクエストを受け取るとio.EOFを返す
		r, err := stream.Recv()
		if err == io.EOF {
			break
		} else if err != nil {
			return err
		}

		chunk := r.GetData()
		log.Printf("stream byte: %d", len(chunk))
		if _, err = buf.Write(chunk); err != nil {
			return err
		}
	}

	data := buf.Bytes()
	mimeType := http.DetectContentType(data)

	h.files[filename] = data

	return stream.SendAndClose(&rpc.ImageUploadResponse{
		Uuid:        uuid,
		Size:        int32(len(data)),
		Filename:    filename,
		ContentType: mimeType,
	})
}
