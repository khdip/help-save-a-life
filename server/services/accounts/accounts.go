package accounts

import (
	"context"

	acntgrpc "github.com/khdip/help-save-a-life/proto/accounts"
	"github.com/khdip/help-save-a-life/server/storage"
)

type AccountsStore interface {
	CreateAccounts(ctx context.Context, ast storage.Accounts) (string, error)
	GetAccounts(ctx context.Context, ast storage.Accounts) (*storage.Accounts, error)
	UpdateAccounts(ctx context.Context, ast storage.Accounts) (*storage.Accounts, error)
	DeleteAccounts(ctx context.Context, ast storage.Accounts) error
	ListAccounts(ctx context.Context, flt storage.Filter) ([]storage.Accounts, error)
	AccountsStats(ctx context.Context, flt storage.Filter) (storage.Stats, error)
}

type Svc struct {
	acntgrpc.UnimplementedAccountsServiceServer
	accst AccountsStore
}

func New(as AccountsStore) *Svc {
	return &Svc{
		accst: as,
	}
}
