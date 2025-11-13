package links

import (
	"context"
	"database/sql"

	linkgrpc "github.com/khdip/help-save-a-life/proto/links"
	"github.com/khdip/help-save-a-life/server/storage"
)

func (s *Svc) DeleteLink(ctx context.Context, req *linkgrpc.DeleteLinkRequest) (*linkgrpc.DeleteLinkResponse, error) {
	if err := s.ls.DeleteLink(ctx, storage.Link{
		ID: req.Link.ID,
		CRUDTimeDate: storage.CRUDTimeDate{
			DeletedBy: sql.NullString{String: req.Link.DeletedBy, Valid: true},
		},
	}); err != nil {
		return nil, err
	}

	return &linkgrpc.DeleteLinkResponse{}, nil
}
