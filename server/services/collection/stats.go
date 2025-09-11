package collection

import (
	"context"

	collgrpc "help-save-a-life/proto/collection"
	"help-save-a-life/server/storage"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *Svc) CollectionStats(ctx context.Context, req *collgrpc.CollectionStatsRequest) (*collgrpc.CollectionStatsResponse, error) {
	r, err := s.cst.CollectionStats(ctx, storage.Filter{
		SearchTerm: req.Filter.SearchTerm,
	})
	if err != nil {
		return nil, status.Error(codes.NotFound, "collection doesn't exist")
	}
	return &collgrpc.CollectionStatsResponse{
		Stats: &collgrpc.Stats{
			Count:       r.Count,
			TotalAmount: r.TotalAmount,
		},
	}, nil
}
