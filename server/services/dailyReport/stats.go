package dailyreport

import (
	"context"

	dregrpc "github.com/khdip/help-save-a-life/proto/dailyReport"
	"github.com/khdip/help-save-a-life/server/storage"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *Svc) DailyReportStats(ctx context.Context, req *dregrpc.DailyReportStatsRequest) (*dregrpc.DailyReportStatsResponse, error) {
	r, err := s.drst.DailyReportStats(ctx, storage.Filter{
		SearchTerm: req.Filter.SearchTerm,
	})
	if err != nil {
		return nil, status.Error(codes.NotFound, "daily report entry doesn't exist")
	}
	return &dregrpc.DailyReportStatsResponse{
		Stats: &dregrpc.Stats{
			Count:       r.Count,
			TotalAmount: r.TotalAmount,
		},
	}, nil
}
