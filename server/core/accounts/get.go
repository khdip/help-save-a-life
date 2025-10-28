package accounts

import (
	"context"
	"help-save-a-life/server/storage"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *CoreSvc) GetAccounts(ctx context.Context, acnt storage.Accounts) (*storage.Accounts, error) {
	at, err := s.st.GetAccounts(ctx, acnt)
	if err != nil {
		return nil, status.Error(codes.Internal, "processing failed")
	}
	return at, nil
}
