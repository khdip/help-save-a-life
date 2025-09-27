package dailyreport

import (
	"context"

	dregrpc "help-save-a-life/proto/dailyReport"
	"help-save-a-life/server/storage"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (s *Svc) GetDailyReport(ctx context.Context, req *dregrpc.GetDailyReportRequest) (*dregrpc.GetDailyReportResponse, error) {
	r, err := s.drst.GetDailyReport(ctx, storage.DailyReport{
		ReportID: req.Dre.ReportID,
	})
	if err != nil {
		return nil, status.Error(codes.NotFound, "daily report doesn't exist")
	}
	return &dregrpc.GetDailyReportResponse{
		Dre: &dregrpc.DailyReport{
			ReportID:     r.ReportID,
			SerialNumber: r.SerialNumber,
			Date:         r.Date,
			Amount:       r.Amount,
			Currency:     r.Currency,
			CreatedAt:    timestamppb.New(r.CreatedAt),
			CreatedBy:    r.CreatedBy,
			UpdatedAt:    timestamppb.New(r.UpdatedAt),
			UpdatedBy:    r.UpdatedBy,
		},
	}, nil
}
