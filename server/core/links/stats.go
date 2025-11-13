package links

import (
	"context"

	"github.com/khdip/help-save-a-life/server/storage"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *CoreSvc) LinkStats(ctx context.Context, filter storage.Filter) (storage.Stats, error) {
	linkStats, err := s.st.LinkStats(ctx, filter)
	if err != nil {
		return storage.Stats{}, status.Error(codes.Internal, "processing failed")
	}

	return linkStats, nil
}
