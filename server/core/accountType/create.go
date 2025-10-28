package accountType

import (
	"context"
	"help-save-a-life/server/storage"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *CoreSvc) CreateAccountType(ctx context.Context, acct storage.AccountType) (string, error) {
	acctid, err := s.st.CreateAccountType(ctx, acct)
	if err != nil {
		return "", status.Error(codes.Internal, "processing failed")
	}

	return acctid, nil
}
