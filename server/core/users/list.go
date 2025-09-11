package users

import (
	"context"
	"help-save-a-life/server/storage"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *CoreSvc) ListUser(ctx context.Context, filter storage.Filter) ([]storage.User, error) {
	lst, err := s.st.ListUser(ctx, filter)
	if err != nil {
		return nil, status.Error(codes.Internal, "processing failed")
	}

	return lst, nil
}
