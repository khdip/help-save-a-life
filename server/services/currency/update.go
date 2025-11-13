package currency

import (
	"context"

	currgrpc "github.com/khdip/help-save-a-life/proto/currency"
	"github.com/khdip/help-save-a-life/server/storage"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *Svc) UpdateCurrency(ctx context.Context, req *currgrpc.UpdateCurrencyRequest) (*currgrpc.UpdateCurrencyResponse, error) {
	res, err := s.cst.UpdateCurrency(ctx, storage.Currency{
		ID:           req.Curr.ID,
		Name:         req.Curr.Name,
		ExchangeRate: req.Curr.ExchangeRate,
		CRUDTimeDate: storage.CRUDTimeDate{
			UpdatedBy: req.Curr.UpdatedBy,
		},
	})
	if err != nil {
		if status.Code(err) != codes.Unknown {
			return nil, err
		}
		return nil, status.Error(codes.Internal, "failed to update record")
	}
	return &currgrpc.UpdateCurrencyResponse{
		Curr: &currgrpc.Currency{
			ID:           res.ID,
			Name:         res.Name,
			ExchangeRate: res.ExchangeRate,
		},
	}, nil
}
