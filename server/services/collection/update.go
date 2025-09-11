package collection

import (
	"context"

	collgrpc "help-save-a-life/proto/collection"
	"help-save-a-life/server/storage"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *Svc) UpdateCollection(ctx context.Context, req *collgrpc.UpdateCollectionRequest) (*collgrpc.UpdateCollectionResponse, error) {
	res, err := s.cst.UpdateCollection(ctx, storage.Collection{
		CollectionID:  req.Coll.CollectionID,
		AccountType:   req.Coll.AccountType,
		AccountNumber: req.Coll.AccountNumber,
		Sender:        req.Coll.Sender,
		Date:          req.Coll.Date,
		Amount:        req.Coll.Amount,
		Currency:      req.Coll.Currency,
		CRUDTimeDate: storage.CRUDTimeDate{
			UpdatedBy: req.Coll.UpdatedBy,
		},
	})
	if err != nil {
		if status.Code(err) != codes.Unknown {
			return nil, err
		}
		return nil, status.Error(codes.Internal, "failed to update record")
	}
	return &collgrpc.UpdateCollectionResponse{
		Coll: &collgrpc.Collection{
			CollectionID:  res.CollectionID,
			AccountType:   res.AccountType,
			AccountNumber: res.AccountNumber,
			Sender:        res.Sender,
			Date:          res.Date,
			Amount:        res.Amount,
			Currency:      res.Currency,
		},
	}, nil
}
