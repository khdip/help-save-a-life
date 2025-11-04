package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"help-save-a-life/server/storage"
)

const insertMedDocs = `
INSERT INTO med_docs (
	name, 
	type,
	uploaded_by
) VALUES (
	:name, 
	:type,
	:uploaded_by
) RETURNING
	id;
`

func (s *Storage) CreateMedDocs(ctx context.Context, md storage.MedDocs) (string, error) {
	stmt, err := s.db.PrepareNamed(insertMedDocs)
	if err != nil {
		return "", err
	}

	var id string
	if err := stmt.Get(&id, md); err != nil {
		return "", err
	}

	return id, nil
}

const getMedDocs = `
SELECT *
FROM med_docs
WHERE id = $1 AND deleted_at IS NULL; 
`

func (s *Storage) GetMedDocs(ctx context.Context, md storage.MedDocs) (*storage.MedDocs, error) {
	var res storage.MedDocs
	if err := s.db.Get(&res, getMedDocs, md.ID); err != nil {
		return nil, fmt.Errorf("executing med doc details: %w", err)
	}
	return &res, nil
}

const deleteMedDocs = `
UPDATE
	med_docs
SET
	deleted_at = now(),
	deleted_by = $1
WHERE 
	id = $2;
`

func (s *Storage) DeleteMedDocs(ctx context.Context, md storage.MedDocs) error {
	_, err := s.db.Exec(deleteMedDocs, md.DeletedBy, md.ID)
	if err != nil {
		return err
	}
	return nil
}

func (s *Storage) ListMedDocs(ctx context.Context, f storage.Filter) ([]storage.MedDocs, error) {
	var md []storage.MedDocs
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

	listMD := fmt.Sprintf("SELECT * FROM med_docs WHERE deleted_at IS NULL AND (name ILIKE '%%' || '%s' || '%%' OR type ILIKE '%%' || '%s' || '%%') ORDER BY %s %s", f.SearchTerm, f.SearchTerm, sortBy, order)
	fullQuery := listMD + limit
	if err := s.db.Select(&md, fullQuery); err != nil {
		if err == sql.ErrNoRows {
			return nil, err
		}
		return nil, err
	}

	return md, nil
}

func (s *Storage) MedDocsStats(ctx context.Context, f storage.Filter) (storage.Stats, error) {
	var mdStat = fmt.Sprintf("SELECT COUNT(*) FROM med_docs where deleted_at IS NULL AND (name ILIKE '%%' || '%s' || '%%' OR type ILIKE '%%' || '%s' || '%%');", f.SearchTerm, f.SearchTerm)
	var stat storage.Stats
	if err := s.db.Get(&stat, mdStat); err != nil {
		if err == sql.ErrNoRows {
			return storage.Stats{}, err
		}
		return storage.Stats{}, err
	}

	return stat, nil
}
