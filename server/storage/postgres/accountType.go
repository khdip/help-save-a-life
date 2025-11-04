package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"help-save-a-life/server/storage"
)

const insertAccountType = `
INSERT INTO account_type (
	title, 
	created_by,
	updated_by
) VALUES (
	:title, 
	:created_by,
	:updated_by
) RETURNING
	id;
`

func (s *Storage) CreateAccountType(ctx context.Context, acct storage.AccountType) (string, error) {
	stmt, err := s.db.PrepareNamed(insertAccountType)
	if err != nil {
		return "", err
	}

	var id string
	if err := stmt.Get(&id, acct); err != nil {
		return "", err
	}

	return id, nil
}

const getAccountType = `
SELECT *
FROM account_type
WHERE id = $1 AND deleted_at IS NULL; 
`

func (s *Storage) GetAccountType(ctx context.Context, acct storage.AccountType) (*storage.AccountType, error) {
	var res storage.AccountType
	if err := s.db.Get(&res, getAccountType, acct.ID); err != nil {
		return nil, fmt.Errorf("executing account type details: %w", err)
	}
	return &res, nil
}

const deleteAccountType = `
UPDATE
	account_type
SET
	deleted_at = now(),
	deleted_by = $1
WHERE 
	id = $2;
`

func (s *Storage) DeleteAccountType(ctx context.Context, acct storage.AccountType) error {
	_, err := s.db.Exec(deleteAccountType, acct.DeletedBy, acct.ID)
	if err != nil {
		return err
	}
	return nil
}

const updateAccountType = `
UPDATE
	 account_type
SET
	title = :title,
	updated_at = now(),
	updated_by = :updated_by
WHERE 
	id = :id
RETURNING 
	updated_at;
`

func (s *Storage) UpdateAccountType(ctx context.Context, acct storage.AccountType) (*storage.AccountType, error) {
	stmt, err := s.db.PrepareNamedContext(ctx, updateAccountType)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()
	if err := stmt.Get(&acct, acct); err != nil {
		if err == sql.ErrNoRows {
			return nil, err
		}
		return nil, fmt.Errorf("executing account type update: %w", err)
	}
	return &acct, nil
}

func (s *Storage) ListAccountType(ctx context.Context, f storage.Filter) ([]storage.AccountType, error) {
	var acct []storage.AccountType
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

	listAcct := fmt.Sprintf("SELECT * FROM account_type WHERE deleted_at IS NULL AND title ILIKE '%%' || '%s' || '%%' ORDER BY %s %s", f.SearchTerm, sortBy, order)
	fullQuery := listAcct + limit
	if err := s.db.Select(&acct, fullQuery); err != nil {
		if err == sql.ErrNoRows {
			return nil, err
		}
		return nil, err
	}

	return acct, nil
}

func (s *Storage) AccountTypeStats(ctx context.Context, f storage.Filter) (storage.Stats, error) {
	var acctStat = fmt.Sprintf("SELECT COUNT(*) FROM account_type where deleted_at IS NULL AND title ILIKE '%%' || '%s' || '%%';", f.SearchTerm)
	var stat storage.Stats
	if err := s.db.Get(&stat, acctStat); err != nil {
		if err == sql.ErrNoRows {
			return storage.Stats{}, err
		}
		return storage.Stats{}, err
	}

	return stat, nil
}
