package users

import (
	"context"

	"github.com/khdip/help-save-a-life/server/storage"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *CoreSvc) UserStats(ctx context.Context, filter storage.Filter) (storage.Stats, error) {
	userStats, err := s.st.UserStats(ctx, filter)
	if err != nil {
		return storage.Stats{}, status.Error(codes.Internal, "processing failed")
	}

	return userStats, nil
}
