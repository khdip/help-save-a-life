package medDocs

import (
	"context"

	docsgrpc "github.com/khdip/help-save-a-life/proto/medDocs"
	"github.com/khdip/help-save-a-life/server/storage"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *Svc) MedDocsStats(ctx context.Context, req *docsgrpc.MedDocsStatsRequest) (*docsgrpc.MedDocsStatsResponse, error) {
	r, err := s.mds.MedDocsStats(ctx, storage.Filter{
		SearchTerm: req.Filter.SearchTerm,
	})
	if err != nil {
		return nil, status.Error(codes.NotFound, "med doc doesn't exist")
	}
	return &docsgrpc.MedDocsStatsResponse{
		Stats: &docsgrpc.Stats{
			Count: r.Count,
		},
	}, nil
}
