package accountType

import (
	"context"
	"help-save-a-life/server/storage"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *CoreSvc) AccountTypeStats(ctx context.Context, filter storage.Filter) (storage.Stats, error) {
	acctStats, err := s.st.AccountTypeStats(ctx, filter)
	if err != nil {
		return storage.Stats{}, status.Error(codes.Internal, "processing failed")
	}

	return acctStats, nil
}
