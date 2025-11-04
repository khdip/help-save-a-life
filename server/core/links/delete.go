package links

import (
	"context"
	"help-save-a-life/server/storage"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *CoreSvc) DeleteLink(ctx context.Context, link storage.Link) error {
	if err := s.st.DeleteLink(ctx, link); err != nil {
		return status.Error(codes.Internal, "processing failed")
	}
	return nil
}
