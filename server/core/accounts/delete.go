package accounts

import (
	"context"
	"help-save-a-life/server/storage"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *CoreSvc) DeleteAccounts(ctx context.Context, acnt storage.Accounts) error {

	if err := s.st.DeleteAccounts(ctx, acnt); err != nil {
		return status.Error(codes.Internal, "processing failed")
	}
	return nil
}
