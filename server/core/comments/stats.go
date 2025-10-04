package comments

import (
	"context"
	"help-save-a-life/server/storage"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *CoreSvc) CommentStats(ctx context.Context, filter storage.Filter) (storage.Stats, error) {
	commStats, err := s.st.CommentStats(ctx, filter)
	if err != nil {
		return storage.Stats{}, status.Error(codes.Internal, "processing failed")
	}

	return commStats, nil
}
