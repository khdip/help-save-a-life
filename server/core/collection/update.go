package collection

import (
	"context"
	"help-save-a-life/server/storage"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *CoreSvc) UpdateCollection(ctx context.Context, coll storage.Collection) (*storage.Collection, error) {
	c, err := s.st.UpdateCollection(ctx, coll)
	if err != nil {
		return nil, status.Error(codes.Internal, "processing failed")
	}

	return c, nil
}
