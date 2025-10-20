package settings

import (
	"context"

	settgrpc "help-save-a-life/proto/settings"
	"help-save-a-life/server/storage"
)

type SettingsStore interface {
	GetSettings(ctx context.Context, sst storage.Settings) (*storage.Settings, error)
	UpdateSettings(ctx context.Context, sst storage.Settings) (*storage.Settings, error)
}

type Svc struct {
	settgrpc.UnimplementedSettingsServiceServer
	sst SettingsStore
}

func New(ss SettingsStore) *Svc {
	return &Svc{
		sst: ss,
	}
}
