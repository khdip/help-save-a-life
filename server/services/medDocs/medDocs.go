package medDocs

import (
	"context"

	docsgrpc "help-save-a-life/proto/medDocs"
	"help-save-a-life/server/storage"
)

type MedDocsStore interface {
	CreateMedDocs(ctx context.Context, md storage.MedDocs) (string, error)
	GetMedDocs(ctx context.Context, md storage.MedDocs) (*storage.MedDocs, error)
	DeleteMedDocs(ctx context.Context, md storage.MedDocs) error
	ListMedDocs(ctx context.Context, flt storage.Filter) ([]storage.MedDocs, error)
	MedDocsStats(ctx context.Context, flt storage.Filter) (storage.Stats, error)
}

type Svc struct {
	docsgrpc.UnimplementedMedDocsServiceServer
	mds MedDocsStore
}

func New(mds MedDocsStore) *Svc {
	return &Svc{
		mds: mds,
	}
}
