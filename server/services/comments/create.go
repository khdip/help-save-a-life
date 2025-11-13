package comments

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	commgrpc "github.com/khdip/help-save-a-life/proto/comments"
	"github.com/khdip/help-save-a-life/server/storage"
)

func (s *Svc) CreateComment(ctx context.Context, req *commgrpc.CreateCommentRequest) (*commgrpc.CreateCommentResponse, error) {
	res, err := s.cst.CreateComment(ctx, storage.Comment{
		CommentID: req.Comm.CommentID,
		Name:      req.Comm.Name,
		Email:     req.Comm.Email,
		Comment:   req.Comm.Comment,
	})
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to create comment")
	}

	return &commgrpc.CreateCommentResponse{
		CommentID: res,
	}, nil
}
