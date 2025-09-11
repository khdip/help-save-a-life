package users

import (
	"context"
	"help-save-a-life/server/storage"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *CoreSvc) UpdateUser(ctx context.Context, user storage.User) (*storage.User, error) {
	u, err := s.st.UpdateUser(ctx, user)
	if err != nil {
		return nil, status.Error(codes.Internal, "processing failed")
	}

	return u, nil
}
