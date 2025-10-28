package accountType

import (
	"context"
	"help-save-a-life/server/storage"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *CoreSvc) ListAccountType(ctx context.Context, filter storage.Filter) ([]storage.AccountType, error) {
	lst, err := s.st.ListAccountType(ctx, filter)
	if err != nil {
		return nil, status.Error(codes.Internal, "processing failed")
	}

	return lst, nil
}
