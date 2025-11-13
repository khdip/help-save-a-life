package postgres

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/khdip/help-save-a-life/server/storage"
)

const insertEntry = `
INSERT INTO daily_report (
	date,
	amount,
	currency,
	created_by,
	updated_by
) VALUES (
	:date,
	:amount,
	:currency,
	:created_by,
	:updated_by
) RETURNING
	report_id;
`

func (s *Storage) CreateDailyReport(ctx context.Context, dre storage.DailyReport) (string, error) {
	stmt, err := s.db.PrepareNamed(insertEntry)
	if err != nil {
		return "", err
	}

	var id string
	if err := stmt.Get(&id, dre); err != nil {
		return "", err
	}

	return id, nil
}

const getEntry = `
SELECT *
FROM daily_report
WHERE report_id = $1 AND deleted_at IS NULL; 
`

func (s *Storage) GetDailyReport(ctx context.Context, dre storage.DailyReport) (*storage.DailyReport, error) {
	var res storage.DailyReport
	if err := s.db.Get(&res, getEntry, dre.ReportID); err != nil {
		return nil, fmt.Errorf("executing daily report details: %w", err)
	}
	return &res, nil
}

const deleteEntry = `
UPDATE
	daily_report
SET
	deleted_at = now(),
	deleted_by = $1
WHERE 
	report_id = $2;

`

func (s *Storage) DeleteDailyReport(ctx context.Context, dre storage.DailyReport) error {
	_, err := s.db.Exec(deleteEntry, dre.DeletedBy, dre.ReportID)
	if err != nil {
		return err
	}
	return nil
}

const updateEntry = `
UPDATE
	daily_report
SET
	date = :date,
	amount = :amount,
	currency = :currency,
	updated_at = now(),
	updated_by = :updated_by
WHERE 
	report_id = :report_id
RETURNING 
	updated_at;
`

func (s *Storage) UpdateDailyReport(ctx context.Context, dre storage.DailyReport) (*storage.DailyReport, error) {
	stmt, err := s.db.PrepareNamedContext(ctx, updateEntry)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()
	if err := stmt.Get(&dre, dre); err != nil {
		if err == sql.ErrNoRows {
			return nil, err
		}
		return nil, fmt.Errorf("executing daily report update: %w", err)
	}
	return &dre, nil
}

func (s *Storage) ListDailyReport(ctx context.Context, f storage.Filter) ([]storage.DailyReport, error) {
	var dre []storage.DailyReport
	order := "ASC"
	sortBy := "serial_number"

	if f.Order != "" {
		order = f.Order
	}
	if f.SortBy != "" {
		sortBy = f.SortBy
	}

	limit := ""
	if f.Limit > 0 {
		limit = fmt.Sprintf(" LIMIT NULLIF(%d, 0) OFFSET %d;", f.Limit, f.Offset)
	}

	listDre := fmt.Sprintf("SELECT * FROM daily_report WHERE deleted_at IS NULL AND date ILIKE '%%' || '%s' || '%%' ORDER BY %s %s", f.SearchTerm, sortBy, order)
	fullQuery := listDre + limit
	if err := s.db.Select(&dre, fullQuery); err != nil {
		if err == sql.ErrNoRows {
			return nil, err
		}
		return nil, err
	}

	return dre, nil
}

func (s *Storage) DailyReportStats(ctx context.Context, f storage.Filter) (storage.Stats, error) {
	var dreStat = fmt.Sprintf("SELECT COUNT(*), COALESCE(SUM(amount * exchange_rate), 0) AS total_amount FROM daily_report LEFT JOIN currency ON currency = currency.id WHERE daily_report.deleted_at IS NULL AND date ILIKE '%%' || '%s' || '%%';", f.SearchTerm)
	var stat storage.Stats
	if err := s.db.Get(&stat, dreStat); err != nil {
		if err == sql.ErrNoRows {
			return storage.Stats{}, err
		}
		return storage.Stats{}, err
	}

	return stat, nil
}
