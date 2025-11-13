package comments

import (
	"context"

	"github.com/khdip/help-save-a-life/server/storage"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *CoreSvc) CreateComment(ctx context.Context, comm storage.Comment) (string, error) {
	commid, err := s.st.CreateComment(ctx, comm)
	if err != nil {
		return "", status.Error(codes.Internal, "processing failed")
	}

	return commid, nil
}
