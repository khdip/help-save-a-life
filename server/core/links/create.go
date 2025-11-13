package links

import (
	"context"

	"github.com/khdip/help-save-a-life/server/storage"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *CoreSvc) CreateLink(ctx context.Context, link storage.Link) (string, error) {
	id, err := s.st.CreateLink(ctx, link)
	if err != nil {
		return "", status.Error(codes.Internal, "processing failed")
	}

	return id, nil
}
