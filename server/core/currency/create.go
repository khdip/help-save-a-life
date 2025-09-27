package currency

import (
	"context"
	"help-save-a-life/server/storage"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *CoreSvc) CreateCurrency(ctx context.Context, curr storage.Currency) (string, error) {
	currid, err := s.st.CreateCurrency(ctx, curr)
	if err != nil {
		return "", status.Error(codes.Internal, "processing failed")
	}

	return currid, nil
}
