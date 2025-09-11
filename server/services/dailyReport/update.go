package dailyreport

import (
	"context"

	dregrpc "help-save-a-life/proto/dailyReport"
	"help-save-a-life/server/storage"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *Svc) UpdateDailyReport(ctx context.Context, req *dregrpc.UpdateDailyReportRequest) (*dregrpc.UpdateDailyReportResponse, error) {
	res, err := s.drst.UpdateDailyReport(ctx, storage.DailyReport{
		ReportID: req.Dre.ReportID,
		Date:     req.Dre.Date,
		Amount:   req.Dre.Amount,
		Currency: req.Dre.Currency,
		CRUDTimeDate: storage.CRUDTimeDate{
			UpdatedBy: req.Dre.UpdatedBy,
		},
	})
	if err != nil {
		if status.Code(err) != codes.Unknown {
			return nil, err
		}
		return nil, status.Error(codes.Internal, "failed to update record")
	}
	return &dregrpc.UpdateDailyReportResponse{
		Dre: &dregrpc.DailyReport{
			ReportID: res.ReportID,
			Date:     res.Date,
			Amount:   res.Amount,
			Currency: res.Currency,
		},
	}, nil
}
