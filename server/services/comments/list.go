package comments

import (
	"context"

	commgrpc "github.com/khdip/help-save-a-life/proto/comments"
	"github.com/khdip/help-save-a-life/server/storage"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (s *Svc) ListComment(ctx context.Context, req *commgrpc.ListCommentRequest) (*commgrpc.ListCommentResponse, error) {
	comm, err := s.cst.ListComment(ctx, storage.Filter{
		Offset:     req.Filter.Offset,
		Limit:      req.Filter.Limit,
		SortBy:     req.Filter.SortBy,
		Order:      req.Filter.Order,
		SearchTerm: req.Filter.SearchTerm,
	})
	if err != nil {
		return nil, status.Error(codes.NotFound, "no comments found")
	}

	list := make([]*commgrpc.Comment, len(comm))
	for i, c := range comm {
		list[i] = &commgrpc.Comment{
			CommentID:    c.CommentID,
			SerialNumber: c.SerialNumber,
			Name:         c.Name,
			Email:        c.Email,
			Comment:      c.Comment,
			CreatedAt:    timestamppb.New(c.CreatedAt),
		}
	}

	return &commgrpc.ListCommentResponse{
		Comm: list,
	}, nil
}
