package currency

import (
	"context"
	"help-save-a-life/server/storage"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *CoreSvc) DeleteCurrency(ctx context.Context, curr storage.Currency) error {

	if err := s.st.DeleteCurrency(ctx, curr); err != nil {
		return status.Error(codes.Internal, "processing failed")
	}
	return nil
}
