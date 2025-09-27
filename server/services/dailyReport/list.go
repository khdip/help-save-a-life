package dailyreport

import (
	"context"

	dregrpc "help-save-a-life/proto/dailyReport"
	"help-save-a-life/server/storage"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	tspb "google.golang.org/protobuf/types/known/timestamppb"
)

func (s *Svc) ListDailyReport(ctx context.Context, req *dregrpc.ListDailyReportRequest) (*dregrpc.ListDailyReportResponse, error) {
	dr, err := s.drst.ListDailyReport(ctx, storage.Filter{
		Offset:     req.Filter.Offset,
		Limit:      req.Filter.Limit,
		SortBy:     req.Filter.SortBy,
		Order:      req.Filter.Order,
		SearchTerm: req.Filter.SearchTerm,
	})
	if err != nil {
		return nil, status.Error(codes.NotFound, "no daily report found")
	}

	list := make([]*dregrpc.DailyReport, len(dr))
	for i, r := range dr {
		list[i] = &dregrpc.DailyReport{
			ReportID:     r.ReportID,
			SerialNumber: r.SerialNumber,
			Date:         r.Date,
			Amount:       r.Amount,
			Currency:     r.Currency,
			CreatedAt:    tspb.New(r.CreatedAt),
			CreatedBy:    r.CreatedBy,
			UpdatedAt:    tspb.New(r.UpdatedAt),
			UpdatedBy:    r.UpdatedBy,
		}
	}

	return &dregrpc.ListDailyReportResponse{
		Dre: list,
	}, nil
}
