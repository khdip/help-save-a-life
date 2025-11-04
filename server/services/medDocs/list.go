package medDocs

import (
	"context"

	docsgrpc "help-save-a-life/proto/medDocs"
	"help-save-a-life/server/storage"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	tspb "google.golang.org/protobuf/types/known/timestamppb"
)

func (s *Svc) ListMedDocs(ctx context.Context, req *docsgrpc.ListMedDocsRequest) (*docsgrpc.ListMedDocsResponse, error) {
	mdss, err := s.mds.ListMedDocs(ctx, storage.Filter{
		Offset:     req.Filter.Offset,
		Limit:      req.Filter.Limit,
		SortBy:     req.Filter.SortBy,
		Order:      req.Filter.Order,
		SearchTerm: req.Filter.SearchTerm,
	})
	if err != nil {
		return nil, status.Error(codes.NotFound, "no med docs found")
	}

	list := make([]*docsgrpc.MedDocs, len(mdss))
	for i, a := range mdss {
		list[i] = &docsgrpc.MedDocs{
			ID:           a.ID,
			SerialNumber: a.SerialNumber,
			Name:         a.Name,
			Type:         a.Type,
			UploadedAt:   tspb.New(a.UploadedAt),
			UploadedBy:   a.UploadedBy,
		}
	}

	return &docsgrpc.ListMedDocsResponse{
		Docs: list,
	}, nil
}
