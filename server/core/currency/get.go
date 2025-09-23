package currency

import (
	"context"
	"help-save-a-life/server/storage"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *CoreSvc) GetCurrency(ctx context.Context, curr storage.Currency) (*storage.Currency, error) {
	c, err := s.st.GetCurrency(ctx, curr)
	if err != nil {
		return nil, status.Error(codes.Internal, "processing failed")
	}
	return c, nil
}
