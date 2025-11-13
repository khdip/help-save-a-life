package links

import (
	"context"

	"github.com/khdip/help-save-a-life/server/storage"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *CoreSvc) UpdateLink(ctx context.Context, link storage.Link) (*storage.Link, error) {
	res, err := s.st.UpdateLink(ctx, link)
	if err != nil {
		return nil, status.Error(codes.Internal, "processing failed")
	}

	return res, nil
}
