package currency

import (
	"context"

	currgrpc "github.com/khdip/help-save-a-life/proto/currency"
	"github.com/khdip/help-save-a-life/server/storage"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (s *Svc) GetCurrency(ctx context.Context, req *currgrpc.GetCurrencyRequest) (*currgrpc.GetCurrencyResponse, error) {
	r, err := s.cst.GetCurrency(ctx, storage.Currency{
		ID: req.Curr.ID,
	})
	if err != nil {
		return nil, status.Error(codes.NotFound, "currency doesn't exist")
	}
	return &currgrpc.GetCurrencyResponse{
		Curr: &currgrpc.Currency{
			ID:           r.ID,
			SerialNumber: r.SerialNumber,
			Name:         r.Name,
			ExchangeRate: r.ExchangeRate,
			CreatedAt:    timestamppb.New(r.CreatedAt),
			CreatedBy:    r.CreatedBy,
			UpdatedAt:    timestamppb.New(r.UpdatedAt),
			UpdatedBy:    r.UpdatedBy,
		},
	}, nil
}
