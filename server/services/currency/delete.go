package currency

import (
	"context"
	"database/sql"

	currgrpc "help-save-a-life/proto/currency"
	"help-save-a-life/server/storage"
)

func (s *Svc) DeleteCurrency(ctx context.Context, req *currgrpc.DeleteCurrencyRequest) (*currgrpc.DeleteCurrencyResponse, error) {
	if err := s.cst.DeleteCurrency(ctx, storage.Currency{
		ID: req.Curr.ID,
		CRUDTimeDate: storage.CRUDTimeDate{
			DeletedBy: sql.NullString{String: req.Curr.DeletedBy, Valid: true},
		},
	}); err != nil {
		return nil, err
	}

	return &currgrpc.DeleteCurrencyResponse{}, nil
}
