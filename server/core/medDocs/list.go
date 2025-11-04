package medDocs

import (
	"context"
	"help-save-a-life/server/storage"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *CoreSvc) ListMedDocs(ctx context.Context, filter storage.Filter) ([]storage.MedDocs, error) {
	lst, err := s.st.ListMedDocs(ctx, filter)
	if err != nil {
		return nil, status.Error(codes.Internal, "processing failed")
	}

	return lst, nil
}
