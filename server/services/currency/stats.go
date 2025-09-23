package currency

import (
	"context"

	currgrpc "help-save-a-life/proto/currency"
	"help-save-a-life/server/storage"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *Svc) CurrencyStats(ctx context.Context, req *currgrpc.CurrencyStatsRequest) (*currgrpc.CurrencyStatsResponse, error) {
	r, err := s.cst.CurrencyStats(ctx, storage.Filter{
		SearchTerm: req.Filter.SearchTerm,
	})
	if err != nil {
		return nil, status.Error(codes.NotFound, "currency doesn't exist")
	}
	return &currgrpc.CurrencyStatsResponse{
		Stats: &currgrpc.Stats{
			Count:       r.Count,
			TotalAmount: r.TotalAmount,
		},
	}, nil
}
