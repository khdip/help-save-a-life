package medDocs

import (
	"context"

	docsgrpc "github.com/khdip/help-save-a-life/proto/medDocs"
	"github.com/khdip/help-save-a-life/server/storage"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (s *Svc) GetMedDocs(ctx context.Context, req *docsgrpc.GetMedDocsRequest) (*docsgrpc.GetMedDocsResponse, error) {
	r, err := s.mds.GetMedDocs(ctx, storage.MedDocs{
		ID: req.Docs.ID,
	})
	if err != nil {
		return nil, status.Error(codes.NotFound, "med doc doesn't exist")
	}
	return &docsgrpc.GetMedDocsResponse{
		Docs: &docsgrpc.MedDocs{
			ID:           r.ID,
			SerialNumber: r.SerialNumber,
			Name:         r.Name,
			Type:         r.Type,
			UploadedAt:   timestamppb.New(r.UploadedAt),
			UploadedBy:   r.UploadedBy,
		},
	}, nil
}
