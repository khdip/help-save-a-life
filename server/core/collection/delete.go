package collection

import (
	"context"

	"github.com/khdip/help-save-a-life/server/storage"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *CoreSvc) DeleteCollection(ctx context.Context, coll storage.Collection) error {

	if err := s.st.DeleteCollection(ctx, coll); err != nil {
		return status.Error(codes.Internal, "processing failed")
	}
	return nil
}
