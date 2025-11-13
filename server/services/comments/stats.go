package comments

import (
	"context"

	commgrpc "github.com/khdip/help-save-a-life/proto/comments"
	"github.com/khdip/help-save-a-life/server/storage"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *Svc) CommentStats(ctx context.Context, req *commgrpc.CommentStatsRequest) (*commgrpc.CommentStatsResponse, error) {
	r, err := s.cst.CommentStats(ctx, storage.Filter{
		SearchTerm: req.Filter.SearchTerm,
	})
	if err != nil {
		return nil, status.Error(codes.NotFound, "comment doesn't exist")
	}
	return &commgrpc.CommentStatsResponse{
		Stats: &commgrpc.Stats{
			Count: r.Count,
		},
	}, nil
}
