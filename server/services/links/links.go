package links

import (
	"context"

	linkgrpc "help-save-a-life/proto/links"
	"help-save-a-life/server/storage"
)

type LinkStore interface {
	CreateLink(ctx context.Context, link storage.Link) (string, error)
	GetLink(ctx context.Context, link storage.Link) (*storage.Link, error)
	UpdateLink(ctx context.Context, link storage.Link) (*storage.Link, error)
	DeleteLink(ctx context.Context, link storage.Link) error
	ListLink(ctx context.Context, flt storage.Filter) ([]storage.Link, error)
	LinkStats(ctx context.Context, flt storage.Filter) (storage.Stats, error)
}

type Svc struct {
	linkgrpc.UnimplementedLinkServiceServer
	ls LinkStore
}

func New(ls LinkStore) *Svc {
	return &Svc{
		ls: ls,
	}
}
