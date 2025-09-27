package currency

import (
	"context"

	"google.golang.org/grpc"

	currgrpc "help-save-a-life/proto/currency"
	"help-save-a-life/server/storage"
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

// RegisterService with grpc server.
func (s *Svc) RegisterSvc(srv *grpc.Server) error {
	currgrpc.RegisterCurrencyServiceServer(srv, s)
	return nil
}
