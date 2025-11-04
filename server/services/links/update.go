package links

import (
	"context"

	linkgrpc "help-save-a-life/proto/links"
	"help-save-a-life/server/storage"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *Svc) UpdateLink(ctx context.Context, req *linkgrpc.UpdateLinkRequest) (*linkgrpc.UpdateLinkResponse, error) {
	res, err := s.ls.UpdateLink(ctx, storage.Link{
		ID:       req.Link.ID,
		LinkText: req.Link.LinkText,
		LinkURL:  req.Link.LinkURL,
		CRUDTimeDate: storage.CRUDTimeDate{
			UpdatedBy: req.Link.UpdatedBy,
		},
	})
	if err != nil {
		if status.Code(err) != codes.Unknown {
			return nil, err
		}
		return nil, status.Error(codes.Internal, "failed to update record")
	}
	return &linkgrpc.UpdateLinkResponse{
		Link: &linkgrpc.Link{
			ID:       res.ID,
			LinkText: res.LinkText,
			LinkURL:  res.LinkURL,
		},
	}, nil
}
