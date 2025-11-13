package currency

import (
	"context"

	currgrpc "github.com/khdip/help-save-a-life/proto/currency"
	"github.com/khdip/help-save-a-life/server/storage"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	tspb "google.golang.org/protobuf/types/known/timestamppb"
)

func (s *Svc) ListCurrency(ctx context.Context, req *currgrpc.ListCurrencyRequest) (*currgrpc.ListCurrencyResponse, error) {
	curr, err := s.cst.ListCurrency(ctx, storage.Filter{
		Offset:     req.Filter.Offset,
		Limit:      req.Filter.Limit,
		SortBy:     req.Filter.SortBy,
		Order:      req.Filter.Order,
		SearchTerm: req.Filter.SearchTerm,
	})
	if err != nil {
		return nil, status.Error(codes.NotFound, "no currency found")
	}

	list := make([]*currgrpc.Currency, len(curr))
	for i, c := range curr {
		list[i] = &currgrpc.Currency{
			ID:           c.ID,
			SerialNumber: c.SerialNumber,
			Name:         c.Name,
			ExchangeRate: c.ExchangeRate,
			CreatedAt:    tspb.New(c.CreatedAt),
			CreatedBy:    c.CreatedBy,
			UpdatedAt:    tspb.New(c.UpdatedAt),
			UpdatedBy:    c.UpdatedBy,
		}
	}

	return &currgrpc.ListCurrencyResponse{
		Curr: list,
	}, nil
}
