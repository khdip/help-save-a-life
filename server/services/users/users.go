package users

import (
	"context"

	usergrpc "github.com/khdip/help-save-a-life/proto/users"
	"github.com/khdip/help-save-a-life/server/storage"
)

type UserStore interface {
	CreateUser(ctx context.Context, ust storage.User) (string, error)
	GetUser(ctx context.Context, ust storage.User) (*storage.User, error)
	UpdateUser(ctx context.Context, ust storage.User) (*storage.User, error)
	DeleteUser(ctx context.Context, ust storage.User) error
	ListUser(ctx context.Context, flt storage.Filter) ([]storage.User, error)
	UserStats(ctx context.Context, flt storage.Filter) (storage.Stats, error)
}

type Svc struct {
	usergrpc.UnimplementedUserServiceServer
	ust UserStore
}

func New(cs UserStore) *Svc {
	return &Svc{
		ust: cs,
	}
}
