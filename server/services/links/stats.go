package links

import (
	"context"

	linkgrpc "help-save-a-life/proto/links"
	"help-save-a-life/server/storage"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *Svc) LinkStats(ctx context.Context, req *linkgrpc.LinkStatsRequest) (*linkgrpc.LinkStatsResponse, error) {
	r, err := s.ls.LinkStats(ctx, storage.Filter{
		SearchTerm: req.Filter.SearchTerm,
	})
	if err != nil {
		return nil, status.Error(codes.NotFound, "link doesn't exist")
	}
	return &linkgrpc.LinkStatsResponse{
		Stats: &linkgrpc.Stats{
			Count: r.Count,
		},
	}, nil
}
