package medDocs

import (
	"context"
	"database/sql"

	docsgrpc "help-save-a-life/proto/medDocs"
	"help-save-a-life/server/storage"
)

func (s *Svc) DeleteMedDocs(ctx context.Context, req *docsgrpc.DeleteMedDocsRequest) (*docsgrpc.DeleteMedDocsResponse, error) {
	if err := s.mds.DeleteMedDocs(ctx, storage.MedDocs{
		ID:        req.Docs.ID,
		DeletedBy: sql.NullString{String: req.Docs.DeletedBy, Valid: true},
	}); err != nil {
		return nil, err
	}

	return &docsgrpc.DeleteMedDocsResponse{}, nil
}
