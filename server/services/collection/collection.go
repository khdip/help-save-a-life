package collection

import (
	"context"

	collgrpc "help-save-a-life/proto/collection"
	"help-save-a-life/server/storage"
)

type CollectionStore interface {
	CreateCollection(ctx context.Context, cst storage.Collection) (string, error)
	GetCollection(ctx context.Context, cst storage.Collection) (*storage.Collection, error)
	UpdateCollection(ctx context.Context, cst storage.Collection) (*storage.Collection, error)
	DeleteCollection(ctx context.Context, cst storage.Collection) error
	ListCollection(ctx context.Context, flt storage.Filter) ([]storage.Collection, error)
	CollectionStats(ctx context.Context, flt storage.Filter) (storage.Stats, error)
}

type Svc struct {
	collgrpc.UnimplementedCollectionServiceServer
	cst CollectionStore
}

func New(cs CollectionStore) *Svc {
	return &Svc{
		cst: cs,
	}
}
