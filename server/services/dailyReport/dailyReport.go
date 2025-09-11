package dailyreport

import (
	"context"

	"google.golang.org/grpc"

	dregrpc "help-save-a-life/proto/dailyReport"
	"help-save-a-life/server/storage"
)

type DailyReportStore interface {
	CreateDailyReport(ctx context.Context, cst storage.DailyReport) (int32, error)
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

// RegisterService with grpc server.
func (s *Svc) RegisterSvc(srv *grpc.Server) error {
	dregrpc.RegisterDailyReportServiceServer(srv, s)
	return nil
}
