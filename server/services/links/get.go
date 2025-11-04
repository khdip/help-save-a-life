package links

import (
	"context"

	linkgrpc "help-save-a-life/proto/links"
	"help-save-a-life/server/storage"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (s *Svc) GetLink(ctx context.Context, req *linkgrpc.GetLinkRequest) (*linkgrpc.GetLinkResponse, error) {
	r, err := s.ls.GetLink(ctx, storage.Link{
		ID: req.Link.ID,
	})
	if err != nil {
		return nil, status.Error(codes.NotFound, "link entry doesn't exist")
	}
	return &linkgrpc.GetLinkResponse{
		Link: &linkgrpc.Link{
			ID:           r.ID,
			SerialNumber: r.SerialNumber,
			LinkText:     r.LinkText,
			LinkURL:      r.LinkURL,
			CreatedAt:    timestamppb.New(r.CreatedAt),
			CreatedBy:    r.CreatedBy,
			UpdatedAt:    timestamppb.New(r.UpdatedAt),
			UpdatedBy:    r.UpdatedBy,
		},
	}, nil
}
