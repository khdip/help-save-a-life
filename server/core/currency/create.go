package currency

import (
	"context"
	"help-save-a-life/server/storage"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *CoreSvc) CreateCurrency(ctx context.Context, curr storage.Currency) (int32, error) {
	currid, err := s.st.CreateCurrency(ctx, curr)
	if err != nil {
		return 0, status.Error(codes.Internal, "processing failed")
	}

	return currid, nil
}
