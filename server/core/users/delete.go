package users

import (
	"context"

	"github.com/khdip/help-save-a-life/server/storage"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *CoreSvc) DeleteUser(ctx context.Context, user storage.User) error {

	if err := s.st.DeleteUser(ctx, user); err != nil {
		return status.Error(codes.Internal, "processing failed")
	}
	return nil
}
