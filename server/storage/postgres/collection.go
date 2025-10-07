package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"help-save-a-life/server/storage"
)

const insertCollection = `
INSERT INTO collection (
	account_type, 
	account_number, 
	sender,
	date,
	amount,
	currency,
	created_by,
	updated_by
) VALUES (
	:account_type, 
	:account_number,
	:sender,
	:date,
	:amount,
	:currency,
	:created_by,
	:updated_by
) RETURNING
	collection_id;
`

func (s *Storage) CreateCollection(ctx context.Context, coll storage.Collection) (string, error) {
	stmt, err := s.db.PrepareNamed(insertCollection)
	if err != nil {
		return "", err
	}

	var id string
	if err := stmt.Get(&id, coll); err != nil {
		fmt.Println(err)
		return "", err
	}

	return id, nil
}

const getCollection = `
SELECT *
FROM collection
WHERE collection_id = $1 AND deleted_at IS NULL; 
`

func (s *Storage) GetCollection(ctx context.Context, coll storage.Collection) (*storage.Collection, error) {
	var res storage.Collection
	if err := s.db.Get(&res, getCollection, coll.CollectionID); err != nil {
		return nil, fmt.Errorf("executing collection details: %w", err)
	}
	return &res, nil
}

const deleteCollection = `
UPDATE
	collection
SET
	deleted_at = now(),
	deleted_by = $1
WHERE 
	collection_id = $2;

`

func (s *Storage) DeleteCollection(ctx context.Context, coll storage.Collection) error {
	_, err := s.db.Exec(deleteCollection, coll.DeletedBy, coll.CollectionID)
	if err != nil {
		return err
	}
	return nil
}

const updateCollection = `
UPDATE
	 collection
SET
	account_type = :account_type,
    account_number = :account_number,
	sender = :sender,
	date = :date,
	amount = :amount,
	currency = :currency,
	updated_at = now(),
	updated_by = :updated_by
WHERE 
	collection_id = :collection_id
RETURNING 
	updated_at;
`

func (s *Storage) UpdateCollection(ctx context.Context, coll storage.Collection) (*storage.Collection, error) {
	stmt, err := s.db.PrepareNamedContext(ctx, updateCollection)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()
	if err := stmt.Get(&coll, coll); err != nil {
		if err == sql.ErrNoRows {
			return nil, err
		}
		return nil, fmt.Errorf("executing collection update: %w", err)
	}
	return &coll, nil
}

func (s *Storage) ListCollection(ctx context.Context, f storage.Filter) ([]storage.Collection, error) {
	var coll []storage.Collection
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

	listColl := fmt.Sprintf("SELECT * FROM collection WHERE deleted_at IS NULL AND (account_number ILIKE '%%' || '%s' || '%%' OR account_type ILIKE '%%' || '%s' || '%%' OR date ILIKE '%%' || '%s' || '%%' OR sender ILIKE '%%' || '%s' || '%%') ORDER BY %s %s", f.SearchTerm, f.SearchTerm, f.SearchTerm, f.SearchTerm, sortBy, order)
	fullQuery := listColl + limit
	if err := s.db.Select(&coll, fullQuery); err != nil {
		if err == sql.ErrNoRows {
			return nil, err
		}
		return nil, err
	}

	return coll, nil
}

func (s *Storage) CollectionStats(ctx context.Context, f storage.Filter) (storage.Stats, error) {
	var collStat = fmt.Sprintf("SELECT COUNT(*), COALESCE(SUM(amount), 0) FROM collection where deleted_at IS NULL AND (account_number ILIKE '%%' || '%s' || '%%' OR account_type ILIKE '%%' || '%s' || '%%' OR date ILIKE '%%' || '%s' || '%%' OR sender ILIKE '%%' || '%s' || '%%');", f.SearchTerm, f.SearchTerm, f.SearchTerm, f.SearchTerm)
	var stat storage.Stats
	if err := s.db.Get(&stat, collStat); err != nil {
		if err == sql.ErrNoRows {
			return storage.Stats{}, err
		}
		return storage.Stats{}, err
	}

	return stat, nil
}
