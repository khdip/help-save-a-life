package accountType

import (
	"context"
	"help-save-a-life/server/storage"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *CoreSvc) UpdateAccountType(ctx context.Context, acct storage.AccountType) (*storage.AccountType, error) {
	at, err := s.st.UpdateAccountType(ctx, acct)
	if err != nil {
		return nil, status.Error(codes.Internal, "processing failed")
	}

	return at, nil
}
