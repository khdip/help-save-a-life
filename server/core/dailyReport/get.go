package dailyreport

import (
	"context"
	"help-save-a-life/server/storage"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *CoreSvc) GetDailyReport(ctx context.Context, dre storage.DailyReport) (*storage.DailyReport, error) {
	dr, err := s.st.GetDailyReport(ctx, dre)
	if err != nil {
		return nil, status.Error(codes.Internal, "processing failed")
	}
	return dr, nil
}
