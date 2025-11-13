package accounts

import (
	"context"
	"database/sql"

	acntgrpc "github.com/khdip/help-save-a-life/proto/accounts"
	"github.com/khdip/help-save-a-life/server/storage"
)

func (s *Svc) DeleteAccounts(ctx context.Context, req *acntgrpc.DeleteAccountsRequest) (*acntgrpc.DeleteAccountsResponse, error) {
	if err := s.accst.DeleteAccounts(ctx, storage.Accounts{
		ID: req.Acnt.ID,
		CRUDTimeDate: storage.CRUDTimeDate{
			DeletedBy: sql.NullString{String: req.Acnt.DeletedBy, Valid: true},
		},
	}); err != nil {
		return nil, err
	}

	return &acntgrpc.DeleteAccountsResponse{}, nil
}
