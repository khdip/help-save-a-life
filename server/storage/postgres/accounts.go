package postgres

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/khdip/help-save-a-life/server/storage"
)

const insertAccount = `
INSERT INTO accounts (
	account_type, 
	details,
	created_by,
	updated_by
) VALUES (
	:account_type, 
	:details,
	:created_by,
	:updated_by
) RETURNING
	id;
`

func (s *Storage) CreateAccounts(ctx context.Context, acnt storage.Accounts) (string, error) {
	stmt, err := s.db.PrepareNamed(insertAccount)
	if err != nil {
		return "", err
	}

	var id string
	if err := stmt.Get(&id, acnt); err != nil {
		return "", err
	}

	return id, nil
}

const getAccount = `
SELECT *
FROM accounts
WHERE id = $1 AND deleted_at IS NULL; 
`

func (s *Storage) GetAccounts(ctx context.Context, acnt storage.Accounts) (*storage.Accounts, error) {
	var res storage.Accounts
	if err := s.db.Get(&res, getAccount, acnt.ID); err != nil {
		return nil, fmt.Errorf("executing account details: %w", err)
	}
	return &res, nil
}

const deleteAccount = `
UPDATE
	accounts
SET
	deleted_at = now(),
	deleted_by = $1
WHERE 
	id = $2;
`

func (s *Storage) DeleteAccounts(ctx context.Context, acnt storage.Accounts) error {
	_, err := s.db.Exec(deleteAccount, acnt.DeletedBy, acnt.ID)
	if err != nil {
		return err
	}
	return nil
}

const updateAccount = `
UPDATE
	 accounts
SET
	account_type = :account_type,
	details = :details,
	updated_at = now(),
	updated_by = :updated_by
WHERE 
	id = :id
RETURNING 
	updated_at;
`

func (s *Storage) UpdateAccounts(ctx context.Context, acnt storage.Accounts) (*storage.Accounts, error) {
	stmt, err := s.db.PrepareNamedContext(ctx, updateAccount)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()
	if err := stmt.Get(&acnt, acnt); err != nil {
		if err == sql.ErrNoRows {
			return nil, err
		}
		return nil, fmt.Errorf("executing account update: %w", err)
	}
	return &acnt, nil
}

func (s *Storage) ListAccounts(ctx context.Context, f storage.Filter) ([]storage.Accounts, error) {
	var acnt []storage.Accounts
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

	listAcnt := fmt.Sprintf("SELECT * FROM accounts WHERE deleted_at IS NULL AND details ILIKE '%%' || '%s' || '%%' ORDER BY %s %s", f.SearchTerm, sortBy, order)
	fullQuery := listAcnt + limit
	if err := s.db.Select(&acnt, fullQuery); err != nil {
		if err == sql.ErrNoRows {
			return nil, err
		}
		return nil, err
	}

	return acnt, nil
}

func (s *Storage) AccountsStats(ctx context.Context, f storage.Filter) (storage.Stats, error) {
	var acntStat = fmt.Sprintf("SELECT COUNT(*) FROM accounts where deleted_at IS NULL AND details ILIKE '%%' || '%s' || '%%';", f.SearchTerm)
	var stat storage.Stats
	if err := s.db.Get(&stat, acntStat); err != nil {
		if err == sql.ErrNoRows {
			return storage.Stats{}, err
		}
		return storage.Stats{}, err
	}

	return stat, nil
}
