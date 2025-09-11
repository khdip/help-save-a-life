package collection

import (
	"context"
	"help-save-a-life/server/storage"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *CoreSvc) CreateCollection(ctx context.Context, coll storage.Collection) (int32, error) {
	collid, err := s.st.CreateCollection(ctx, coll)
	if err != nil {
		return 0, status.Error(codes.Internal, "processing failed")
	}

	return collid, nil
}
