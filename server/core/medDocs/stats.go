package medDocs

import (
	"context"

	"github.com/khdip/help-save-a-life/server/storage"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *CoreSvc) MedDocsStats(ctx context.Context, filter storage.Filter) (storage.Stats, error) {
	Stats, err := s.st.MedDocsStats(ctx, filter)
	if err != nil {
		return storage.Stats{}, status.Error(codes.Internal, "processing failed")
	}

	return Stats, nil
}
