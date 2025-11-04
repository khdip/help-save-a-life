package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"help-save-a-life/server/storage"
)

const insertLink = `
INSERT INTO links (
	link_text, 
	link_url,
	created_by,
	updated_by
) VALUES (
	:link_text, 
	:link_url,
	:created_by,
	:updated_by
) RETURNING
	id;
`

func (s *Storage) CreateLink(ctx context.Context, link storage.Link) (string, error) {
	stmt, err := s.db.PrepareNamed(insertLink)
	if err != nil {
		return "", err
	}

	var id string
	if err := stmt.Get(&id, link); err != nil {
		return "", err
	}

	return id, nil
}

const getLink = `
SELECT *
FROM links
WHERE id = $1 AND deleted_at IS NULL; 
`

func (s *Storage) GetLink(ctx context.Context, link storage.Link) (*storage.Link, error) {
	var res storage.Link
	if err := s.db.Get(&res, getLink, link.ID); err != nil {
		return nil, fmt.Errorf("executing link details: %w", err)
	}
	return &res, nil
}

const deleteLink = `
UPDATE
	links
SET
	deleted_at = now(),
	deleted_by = $1
WHERE 
	id = $2;
`

func (s *Storage) DeleteLink(ctx context.Context, link storage.Link) error {
	_, err := s.db.Exec(deleteLink, link.DeletedBy, link.ID)
	if err != nil {
		return err
	}
	return nil
}

const updateLink = `
UPDATE
	 links
SET
	link_text = :link_text,
	link_url = :link_url
	updated_at = now(),
	updated_by = :updated_by
WHERE 
	id = :id
RETURNING 
	updated_at;
`

func (s *Storage) UpdateLink(ctx context.Context, link storage.Link) (*storage.Link, error) {
	stmt, err := s.db.PrepareNamedContext(ctx, updateLink)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()
	if err := stmt.Get(&link, link); err != nil {
		if err == sql.ErrNoRows {
			return nil, err
		}
		return nil, fmt.Errorf("executing link update: %w", err)
	}
	return &link, nil
}

func (s *Storage) ListLink(ctx context.Context, f storage.Filter) ([]storage.Link, error) {
	var link []storage.Link
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

	listLink := fmt.Sprintf("SELECT * FROM links WHERE deleted_at IS NULL AND (link_text ILIKE '%%' || '%s' || '%%' OR link_url ILIKE '%%' || '%s' || '%%') ORDER BY %s %s", f.SearchTerm, f.SearchTerm, sortBy, order)
	fullQuery := listLink + limit
	if err := s.db.Select(&link, fullQuery); err != nil {
		if err == sql.ErrNoRows {
			return nil, err
		}
		return nil, err
	}

	return link, nil
}

func (s *Storage) LinkStats(ctx context.Context, f storage.Filter) (storage.Stats, error) {
	var linkStat = fmt.Sprintf("SELECT COUNT(*) FROM links where deleted_at IS NULL AND (link_text ILIKE '%%' || '%s' || '%%' OR link_url ILIKE '%%' || '%s' || '%%');", f.SearchTerm, f.SearchTerm)
	var stat storage.Stats
	if err := s.db.Get(&stat, linkStat); err != nil {
		if err == sql.ErrNoRows {
			return storage.Stats{}, err
		}
		return storage.Stats{}, err
	}

	return stat, nil
}
