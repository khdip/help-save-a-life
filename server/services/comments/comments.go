package comments

import (
	"context"

	"google.golang.org/grpc"

	commgrpc "help-save-a-life/proto/comments"
	"help-save-a-life/server/storage"
)

type CommentStore interface {
	CreateComment(ctx context.Context, cst storage.Comment) (string, error)
	GetComment(ctx context.Context, cst storage.Comment) (*storage.Comment, error)
	ListComment(ctx context.Context, flt storage.Filter) ([]storage.Comment, error)
}

type Svc struct {
	commgrpc.UnimplementedCommentServiceServer
	cst CommentStore
}

func New(cs CommentStore) *Svc {
	return &Svc{
		cst: cs,
	}
}

// RegisterService with grpc server.
func (s *Svc) RegisterSvc(srv *grpc.Server) error {
	commgrpc.RegisterCommentServiceServer(srv, s)
	return nil
}
