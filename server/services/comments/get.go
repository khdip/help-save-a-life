package comments

import (
	"context"

	commgrpc "help-save-a-life/proto/comments"
	"help-save-a-life/server/storage"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (s *Svc) GetComment(ctx context.Context, req *commgrpc.GetCommentRequest) (*commgrpc.GetCommentResponse, error) {
	r, err := s.cst.GetComment(ctx, storage.Comment{
		CommentID: req.Comm.CommentID,
	})
	if err != nil {
		return nil, status.Error(codes.NotFound, "comment doesn't exist")
	}
	return &commgrpc.GetCommentResponse{
		Comm: &commgrpc.Comment{
			CommentID: r.CommentID,
			Name:      r.Name,
			Email:     r.Email,
			Comment:   r.Comment,
			CreatedAt: timestamppb.New(r.CreatedAt),
		},
	}, nil
}
