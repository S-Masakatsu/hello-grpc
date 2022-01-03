package handler

import (
	"context"
	"math/rand"
	"sync"
	"time"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/S-Masakatsu/hello-grpc/bake-pancake/rpc"
)

type BakerHandler struct {
	report *report
	rpc.UnsafePancakeBakerServiceServer
}

type report struct {
	sync.Mutex // 複数人が同時に焼いても大丈夫なようにする
	data       map[rpc.Pancake_Menu]int
}

func NewBakerHandler() rpc.PancakeBakerServiceServer {
	return &BakerHandler{
		report: &report{
			data: make(map[rpc.Pancake_Menu]int),
		},
	}
}

func (h *BakerHandler) Bake(ctx context.Context, req *rpc.BakeRequest) (*rpc.BakeResponse, error) {
	// バリデーション
	if req.Menu == rpc.Pancake_UNKNOWN || req.Menu > rpc.Pancake_SPICY_CURRY {
		return nil, status.Errorf(codes.InvalidArgument, "invalid pancake")
	}

	// パンケーキを焼いて数を登録
	now := time.Now()
	h.report.Lock()
	h.report.data[req.Menu] = h.report.data[req.Menu] + 1
	h.report.Unlock()

	return &rpc.BakeResponse{
		Pancake: &rpc.Pancake{
			Menu:           req.Menu,
			ChefName:       "gami", // ワンオペ
			TechnicalScore: rand.Float32(),
			CreatedAt: &timestamppb.Timestamp{
				Seconds: now.Unix(),
				Nanos:   int32(now.Nanosecond()),
			},
		},
	}, nil
}

func (h *BakerHandler) Report(ctx context.Context, req *rpc.ReportRequest) (*rpc.ReportResponse, error) {
	counts := make([]*rpc.Report_BakeCount, 0)

	// レポート作成
	h.report.Lock()
	for k, v := range h.report.data {
		counts = append(counts, &rpc.Report_BakeCount{
			Menu:  k,
			Count: int32(v),
		})
	}
	h.report.Unlock()

	return &rpc.ReportResponse{
		Report: &rpc.Report{
			BakeCounts: counts,
		},
	}, nil
}
