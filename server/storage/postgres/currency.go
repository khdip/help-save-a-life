package postgres

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/khdip/help-save-a-life/server/storage"
)

const insertCurrency = `
INSERT INTO currency (
	name, 
	exchange_rate, 
	created_by,
	updated_by
) VALUES (
	:name, 
	:exchange_rate,
	:created_by,
	:updated_by
) RETURNING
	id;
`

func (s *Storage) CreateCurrency(ctx context.Context, curr storage.Currency) (string, error) {
	stmt, err := s.db.PrepareNamed(insertCurrency)
	if err != nil {
		return "", err
	}

	var id string
	if err := stmt.Get(&id, curr); err != nil {
		return "", err
	}

	return id, nil
}

const getCurrency = `
SELECT *
FROM currency
WHERE id = $1 AND deleted_at IS NULL; 
`

func (s *Storage) GetCurrency(ctx context.Context, curr storage.Currency) (*storage.Currency, error) {
	var res storage.Currency
	if err := s.db.Get(&res, getCurrency, curr.ID); err != nil {
		return nil, fmt.Errorf("executing currency details: %w", err)
	}
	return &res, nil
}

const deleteCurrency = `
UPDATE
	currency
SET
	deleted_at = now(),
	deleted_by = $1
WHERE 
	id = $2;

`

func (s *Storage) DeleteCurrency(ctx context.Context, curr storage.Currency) error {
	_, err := s.db.Exec(deleteCurrency, curr.DeletedBy, curr.ID)
	if err != nil {
		return err
	}
	return nil
}

const updateCurrency = `
UPDATE
	 currency
SET
	name = :name,
    exchange_rate = :exchange_rate,
	updated_at = now(),
	updated_by = :updated_by
WHERE 
	id = :id
RETURNING 
	updated_at;
`

func (s *Storage) UpdateCurrency(ctx context.Context, curr storage.Currency) (*storage.Currency, error) {
	stmt, err := s.db.PrepareNamedContext(ctx, updateCurrency)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()
	if err := stmt.Get(&curr, curr); err != nil {
		if err == sql.ErrNoRows {
			return nil, err
		}
		return nil, fmt.Errorf("executing currency update: %w", err)
	}
	return &curr, nil
}

func (s *Storage) ListCurrency(ctx context.Context, f storage.Filter) ([]storage.Currency, error) {
	var curr []storage.Currency
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

	listCurr := fmt.Sprintf("SELECT * FROM currency WHERE deleted_at IS NULL AND name ILIKE '%%' || '%s' || '%%' ORDER BY %s %s", f.SearchTerm, sortBy, order)
	fullQuery := listCurr + limit
	if err := s.db.Select(&curr, fullQuery); err != nil {
		if err == sql.ErrNoRows {
			return nil, err
		}
		return nil, err
	}

	return curr, nil
}

func (s *Storage) CurrencyStats(ctx context.Context, f storage.Filter) (storage.Stats, error) {
	var currStat = fmt.Sprintf("SELECT COUNT(*) FROM currency where deleted_at IS NULL AND name ILIKE '%%' || '%s' || '%%';", f.SearchTerm)
	var stat storage.Stats
	if err := s.db.Get(&stat, currStat); err != nil {
		if err == sql.ErrNoRows {
			return storage.Stats{}, err
		}
		return storage.Stats{}, err
	}

	return stat, nil
}
