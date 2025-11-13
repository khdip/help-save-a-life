package dailyreport

import (
	"context"

	"github.com/khdip/help-save-a-life/server/storage"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *CoreSvc) CreateDailyReport(ctx context.Context, dre storage.DailyReport) (string, error) {
	rid, err := s.st.CreateDailyReport(ctx, dre)
	if err != nil {
		return "", status.Error(codes.Internal, "processing failed")
	}

	return rid, nil
}
