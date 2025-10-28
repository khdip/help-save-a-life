package accountType

import (
	"context"

	acctgrpc "help-save-a-life/proto/accountType"
	"help-save-a-life/server/storage"
)

type AccountTypeStore interface {
	CreateAccountType(ctx context.Context, atst storage.AccountType) (string, error)
	GetAccountType(ctx context.Context, atst storage.AccountType) (*storage.AccountType, error)
	UpdateAccountType(ctx context.Context, atst storage.AccountType) (*storage.AccountType, error)
	DeleteAccountType(ctx context.Context, atst storage.AccountType) error
	ListAccountType(ctx context.Context, flt storage.Filter) ([]storage.AccountType, error)
	AccountTypeStats(ctx context.Context, flt storage.Filter) (storage.Stats, error)
}

type Svc struct {
	acctgrpc.UnimplementedAccountTypeServiceServer
	atst AccountTypeStore
}

func New(ats AccountTypeStore) *Svc {
	return &Svc{
		atst: ats,
	}
}
