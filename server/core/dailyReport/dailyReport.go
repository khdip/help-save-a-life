package dailyreport

import "help-save-a-life/server/storage/postgres"

type CoreSvc struct {
	st *postgres.Storage
}

func New(st *postgres.Storage) *CoreSvc {
	return &CoreSvc{
		st: st,
	}
}
