package dailyreport

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	dregrpc "help-save-a-life/proto/dailyReport"
	"help-save-a-life/server/storage"
)

func (s *Svc) CreateDailyReport(ctx context.Context, req *dregrpc.CreateDailyReportRequest) (*dregrpc.CreateDailyReportResponse, error) {
	res, err := s.drst.CreateDailyReport(ctx, storage.DailyReport{
		ReportID: req.Dre.ReportID,
		Date:     req.Dre.Date,
		Amount:   req.Dre.Amount,
		Currency: req.Dre.Currency,
		CRUDTimeDate: storage.CRUDTimeDate{
			CreatedBy: req.Dre.CreatedBy,
			UpdatedBy: req.Dre.UpdatedBy,
		},
	})
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to create daily report entry")
	}

	return &dregrpc.CreateDailyReportResponse{
		ReportID: res,
	}, nil
}
