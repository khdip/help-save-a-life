package accountType

import (
	"context"

	"github.com/khdip/help-save-a-life/server/storage"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *CoreSvc) DeleteAccountType(ctx context.Context, acct storage.AccountType) error {
	if err := s.st.DeleteAccountType(ctx, acct); err != nil {
		return status.Error(codes.Internal, "processing failed")
	}
	return nil
}
