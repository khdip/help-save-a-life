package currency

import (
	"context"
	"help-save-a-life/server/storage"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *CoreSvc) CurrencyStats(ctx context.Context, filter storage.Filter) (storage.Stats, error) {
	currStats, err := s.st.CurrencyStats(ctx, filter)
	if err != nil {
		return storage.Stats{}, status.Error(codes.Internal, "processing failed")
	}

	return currStats, nil
}
