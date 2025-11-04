package medDocs

import (
	"context"
	"help-save-a-life/server/storage"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *CoreSvc) CreateMedDocs(ctx context.Context, md storage.MedDocs) (string, error) {
	id, err := s.st.CreateMedDocs(ctx, md)
	if err != nil {
		return "", status.Error(codes.Internal, "processing failed")
	}

	return id, nil
}
