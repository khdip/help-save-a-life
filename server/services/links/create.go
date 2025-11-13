package links

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	linkgrpc "github.com/khdip/help-save-a-life/proto/links"
	"github.com/khdip/help-save-a-life/server/storage"
)

func (s *Svc) CreateLink(ctx context.Context, req *linkgrpc.CreateLinkRequest) (*linkgrpc.CreateLinkResponse, error) {
	res, err := s.ls.CreateLink(ctx, storage.Link{
		LinkText: req.Link.LinkText,
		LinkURL:  req.Link.LinkURL,
		CRUDTimeDate: storage.CRUDTimeDate{
			CreatedBy: req.Link.CreatedBy,
			UpdatedBy: req.Link.UpdatedBy,
		},
	})
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to create link entry")
	}

	return &linkgrpc.CreateLinkResponse{
		ID: res,
	}, nil
}
