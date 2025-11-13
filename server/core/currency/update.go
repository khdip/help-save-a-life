package currency

import (
	"context"

	"github.com/khdip/help-save-a-life/server/storage"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *CoreSvc) UpdateCurrency(ctx context.Context, curr storage.Currency) (*storage.Currency, error) {
	c, err := s.st.UpdateCurrency(ctx, curr)
	if err != nil {
		return nil, status.Error(codes.Internal, "processing failed")
	}

	return c, nil
}
