package collection

import (
	"context"

	collgrpc "github.com/khdip/help-save-a-life/proto/collection"
	"github.com/khdip/help-save-a-life/server/storage"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	tspb "google.golang.org/protobuf/types/known/timestamppb"
)

func (s *Svc) ListCollection(ctx context.Context, req *collgrpc.ListCollectionRequest) (*collgrpc.ListCollectionResponse, error) {
	coll, err := s.cst.ListCollection(ctx, storage.Filter{
		Offset:     req.Filter.Offset,
		Limit:      req.Filter.Limit,
		SortBy:     req.Filter.SortBy,
		Order:      req.Filter.Order,
		SearchTerm: req.Filter.SearchTerm,
	})
	if err != nil {
		return nil, status.Error(codes.NotFound, "no collection found")
	}

	list := make([]*collgrpc.Collection, len(coll))
	for i, c := range coll {
		list[i] = &collgrpc.Collection{
			CollectionID:  c.CollectionID,
			SerialNumber:  c.SerialNumber,
			AccountType:   c.AccountType,
			AccountNumber: c.AccountNumber,
			Sender:        c.Sender,
			Date:          c.Date,
			Amount:        c.Amount,
			Currency:      c.Currency,
			CreatedAt:     tspb.New(c.CreatedAt),
			CreatedBy:     c.CreatedBy,
			UpdatedAt:     tspb.New(c.UpdatedAt),
			UpdatedBy:     c.UpdatedBy,
		}
	}

	return &collgrpc.ListCollectionResponse{
		Coll: list,
	}, nil
}
