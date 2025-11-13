package accounts

import (
	"context"

	"github.com/khdip/help-save-a-life/server/storage"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *CoreSvc) CreateAccounts(ctx context.Context, acnt storage.Accounts) (string, error) {
	acntid, err := s.st.CreateAccounts(ctx, acnt)
	if err != nil {
		return "", status.Error(codes.Internal, "processing failed")
	}

	return acntid, nil
}
