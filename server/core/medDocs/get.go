package medDocs

import (
	"context"
	"help-save-a-life/server/storage"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *CoreSvc) GetMedDocs(ctx context.Context, md storage.MedDocs) (*storage.MedDocs, error) {
	res, err := s.st.GetMedDocs(ctx, md)
	if err != nil {
		return nil, status.Error(codes.Internal, "processing failed")
	}
	return res, nil
}
