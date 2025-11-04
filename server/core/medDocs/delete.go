package medDocs

import (
	"context"
	"help-save-a-life/server/storage"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *CoreSvc) DeleteMedDocs(ctx context.Context, md storage.MedDocs) error {
	if err := s.st.DeleteMedDocs(ctx, md); err != nil {
		return status.Error(codes.Internal, "processing failed")
	}
	return nil
}
