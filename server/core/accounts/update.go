package accounts

import (
	"context"

	"github.com/khdip/help-save-a-life/server/storage"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *CoreSvc) UpdateAccounts(ctx context.Context, acnt storage.Accounts) (*storage.Accounts, error) {
	acc, err := s.st.UpdateAccounts(ctx, acnt)
	if err != nil {
		return nil, status.Error(codes.Internal, "processing failed")
	}

	return acc, nil
}
