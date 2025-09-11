package dailyreport

import (
	"context"
	"help-save-a-life/server/storage"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *CoreSvc) ListDailyReport(ctx context.Context, filter storage.Filter) ([]storage.DailyReport, error) {
	lst, err := s.st.ListDailyReport(ctx, filter)
	if err != nil {
		return nil, status.Error(codes.Internal, "processing failed")
	}

	return lst, nil
}
