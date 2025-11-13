package collection

import (
	"context"
	"database/sql"

	collgrpc "github.com/khdip/help-save-a-life/proto/collection"
	"github.com/khdip/help-save-a-life/server/storage"
)

func (s *Svc) DeleteCollection(ctx context.Context, req *collgrpc.DeleteCollectionRequest) (*collgrpc.DeleteCollectionResponse, error) {
	if err := s.cst.DeleteCollection(ctx, storage.Collection{
		CollectionID: req.Coll.CollectionID,
		CRUDTimeDate: storage.CRUDTimeDate{
			DeletedBy: sql.NullString{String: req.Coll.DeletedBy, Valid: true},
		},
	}); err != nil {
		return nil, err
	}

	return &collgrpc.DeleteCollectionResponse{}, nil
}
