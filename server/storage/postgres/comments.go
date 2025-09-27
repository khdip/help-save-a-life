package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"help-save-a-life/server/storage"
)

const insertComment = `
INSERT INTO comments (
	name, 
	email, 
	comment
) VALUES (
	:name, 
	:email,
	:comment
) RETURNING
	comment_id;
`

func (s *Storage) CreateComment(ctx context.Context, comm storage.Comment) (string, error) {
	stmt, err := s.db.PrepareNamed(insertComment)
	if err != nil {
		return "", err
	}

	var id string
	if err := stmt.Get(&id, comm); err != nil {
		return "", err
	}

	return id, nil
}

const getComment = `
SELECT *
FROM comments
WHERE comment_id = $1; 
`

func (s *Storage) GetComment(ctx context.Context, comm storage.Comment) (*storage.Comment, error) {
	var res storage.Comment
	if err := s.db.Get(&res, getComment, comm.CommentID); err != nil {
		return nil, fmt.Errorf("executing comment details: %w", err)
	}
	return &res, nil
}

func (s *Storage) ListComment(ctx context.Context, f storage.Filter) ([]storage.Comment, error) {
	var comm []storage.Comment
	order := "DESC"
	sortBy := "created_at"

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

	listColl := fmt.Sprintf("SELECT * FROM comments WHERE name ILIKE '%%' || '%s' || '%%' ORDER BY %s %s", f.SearchTerm, sortBy, order)
	fullQuery := listColl + limit
	if err := s.db.Select(&comm, fullQuery); err != nil {
		if err == sql.ErrNoRows {
			return nil, err
		}
		return nil, err
	}

	return comm, nil
}
