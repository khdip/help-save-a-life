package currency

import (
	"context"

	currgrpc "github.com/khdip/help-save-a-life/proto/currency"
	"github.com/khdip/help-save-a-life/server/storage"
)

type CurrencyStore interface {
	CreateCurrency(ctx context.Context, cst storage.Currency) (string, error)
	GetCurrency(ctx context.Context, cst storage.Currency) (*storage.Currency, error)
	UpdateCurrency(ctx context.Context, cst storage.Currency) (*storage.Currency, error)
	DeleteCurrency(ctx context.Context, cst storage.Currency) error
	ListCurrency(ctx context.Context, flt storage.Filter) ([]storage.Currency, error)
	CurrencyStats(ctx context.Context, flt storage.Filter) (storage.Stats, error)
}

type Svc struct {
	currgrpc.UnimplementedCurrencyServiceServer
	cst CurrencyStore
}

func New(cs CurrencyStore) *Svc {
	return &Svc{
		cst: cs,
	}
}
