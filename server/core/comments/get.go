package comments

import (
	"context"
	"help-save-a-life/server/storage"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *CoreSvc) GetComment(ctx context.Context, comm storage.Comment) (*storage.Comment, error) {
	c, err := s.st.GetComment(ctx, comm)
	if err != nil {
		return nil, status.Error(codes.Internal, "processing failed")
	}
	return c, nil
}
