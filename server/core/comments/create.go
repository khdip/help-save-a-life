package comments

import (
	"context"
	"help-save-a-life/server/storage"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *CoreSvc) CreateComment(ctx context.Context, comm storage.Comment) (int32, error) {
	commid, err := s.st.CreateComment(ctx, comm)
	if err != nil {
		return 0, status.Error(codes.Internal, "processing failed")
	}

	return commid, nil
}
