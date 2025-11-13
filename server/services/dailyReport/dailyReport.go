package dailyreport

import (
	"context"

	dregrpc "github.com/khdip/help-save-a-life/proto/dailyReport"
	"github.com/khdip/help-save-a-life/server/storage"
)

type DailyReportStore interface {
	CreateDailyReport(ctx context.Context, cst storage.DailyReport) (string, error)
	GetDailyReport(ctx context.Context, cst storage.DailyReport) (*storage.DailyReport, error)
	UpdateDailyReport(ctx context.Context, cst storage.DailyReport) (*storage.DailyReport, error)
	DeleteDailyReport(ctx context.Context, cst storage.DailyReport) error
	ListDailyReport(ctx context.Context, flt storage.Filter) ([]storage.DailyReport, error)
	DailyReportStats(ctx context.Context, flt storage.Filter) (storage.Stats, error)
}

type Svc struct {
	dregrpc.UnimplementedDailyReportServiceServer
	drst DailyReportStore
}

func New(cs DailyReportStore) *Svc {
	return &Svc{
		drst: cs,
	}
}
