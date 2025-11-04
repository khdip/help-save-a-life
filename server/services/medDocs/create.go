package medDocs

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	docsgrpc "help-save-a-life/proto/medDocs"
	"help-save-a-life/server/storage"
)

func (s *Svc) CreateMedDocs(ctx context.Context, req *docsgrpc.CreateMedDocsRequest) (*docsgrpc.CreateMedDocsResponse, error) {
	res, err := s.mds.CreateMedDocs(ctx, storage.MedDocs{
		Name:       req.Docs.Name,
		Type:       req.Docs.Type,
		UploadedBy: req.Docs.UploadedBy,
	})
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to create med doc")
	}

	return &docsgrpc.CreateMedDocsResponse{
		ID: res,
	}, nil
}
