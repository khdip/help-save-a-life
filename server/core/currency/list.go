package currency

import (
	"context"
	"help-save-a-life/server/storage"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *CoreSvc) ListCurrency(ctx context.Context, filter storage.Filter) ([]storage.Currency, error) {
	lst, err := s.st.ListCurrency(ctx, filter)
	if err != nil {
		return nil, status.Error(codes.Internal, "processing failed")
	}

	return lst, nil
}
