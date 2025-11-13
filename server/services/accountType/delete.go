package accountType

import (
	"context"
	"database/sql"

	acctgrpc "github.com/khdip/help-save-a-life/proto/accountType"
	"github.com/khdip/help-save-a-life/server/storage"
)

func (s *Svc) DeleteAccountType(ctx context.Context, req *acctgrpc.DeleteAccountTypeRequest) (*acctgrpc.DeleteAccountTypeResponse, error) {
	if err := s.atst.DeleteAccountType(ctx, storage.AccountType{
		ID: req.Acct.ID,
		CRUDTimeDate: storage.CRUDTimeDate{
			DeletedBy: sql.NullString{String: req.Acct.DeletedBy, Valid: true},
		},
	}); err != nil {
		return nil, err
	}

	return &acctgrpc.DeleteAccountTypeResponse{}, nil
}
