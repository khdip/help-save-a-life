package currency

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	currgrpc "help-save-a-life/proto/currency"
	"help-save-a-life/server/storage"
)

func (s *Svc) CreateCurrency(ctx context.Context, req *currgrpc.CreateCurrencyRequest) (*currgrpc.CreateCurrencyResponse, error) {
	res, err := s.cst.CreateCurrency(ctx, storage.Currency{
		Name:         req.Curr.Name,
		ExchangeRate: req.Curr.ExchangeRate,
		CRUDTimeDate: storage.CRUDTimeDate{
			CreatedBy: req.Curr.CreatedBy,
			UpdatedBy: req.Curr.UpdatedBy,
		},
	})
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to create currency")
	}

	return &currgrpc.CreateCurrencyResponse{
		ID: res,
	}, nil
}
