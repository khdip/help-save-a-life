package users

import (
	"context"

	"github.com/khdip/help-save-a-life/server/storage"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *CoreSvc) CreateUser(ctx context.Context, user storage.User) (string, error) {
	userid, err := s.st.CreateUser(ctx, user)
	if err != nil {
		return "", status.Error(codes.Internal, "processing failed")
	}

	return userid, nil
}
