package accounts

import (
	"context"
	"help-save-a-life/server/storage"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *CoreSvc) AccountsStats(ctx context.Context, filter storage.Filter) (storage.Stats, error) {
	acntStats, err := s.st.AccountsStats(ctx, filter)
	if err != nil {
		return storage.Stats{}, status.Error(codes.Internal, "processing failed")
	}

	return acntStats, nil
}
