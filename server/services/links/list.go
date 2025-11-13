package links

import (
	"context"

	linkgrpc "github.com/khdip/help-save-a-life/proto/links"
	"github.com/khdip/help-save-a-life/server/storage"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	tspb "google.golang.org/protobuf/types/known/timestamppb"
)

func (s *Svc) ListLink(ctx context.Context, req *linkgrpc.ListLinkRequest) (*linkgrpc.ListLinkResponse, error) {
	lnks, err := s.ls.ListLink(ctx, storage.Filter{
		Offset:     req.Filter.Offset,
		Limit:      req.Filter.Limit,
		SortBy:     req.Filter.SortBy,
		Order:      req.Filter.Order,
		SearchTerm: req.Filter.SearchTerm,
	})
	if err != nil {
		return nil, status.Error(codes.NotFound, "no link found")
	}

	list := make([]*linkgrpc.Link, len(lnks))
	for i, a := range lnks {
		list[i] = &linkgrpc.Link{
			ID:           a.ID,
			SerialNumber: a.SerialNumber,
			LinkText:     a.LinkText,
			LinkURL:      a.LinkURL,
			CreatedAt:    tspb.New(a.CreatedAt),
			CreatedBy:    a.CreatedBy,
			UpdatedAt:    tspb.New(a.UpdatedAt),
			UpdatedBy:    a.UpdatedBy,
		}
	}

	return &linkgrpc.ListLinkResponse{
		Link: list,
	}, nil
}
