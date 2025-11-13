package users

import (
	"context"

	"github.com/khdip/help-save-a-life/server/storage"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *CoreSvc) GetUser(ctx context.Context, user storage.User) (*storage.User, error) {
	u, err := s.st.GetUser(ctx, user)
	if err != nil {
		return nil, status.Error(codes.Internal, "processing failed")
	}
	return u, nil
}
