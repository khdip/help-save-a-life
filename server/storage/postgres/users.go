package postgres

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/khdip/help-save-a-life/server/storage"
)

const insertUser = `
INSERT INTO users (
	name, 
	batch, 
	email,
	password,
	created_by,
	updated_by
) VALUES (
	:name, 
	:batch,
	:email,
	:password,
	:created_by,
	:updated_by
) RETURNING
	user_id;
`

func (s *Storage) CreateUser(ctx context.Context, user storage.User) (string, error) {
	stmt, err := s.db.PrepareNamed(insertUser)
	if err != nil {
		return "", err
	}

	var id string
	if err := stmt.Get(&id, user); err != nil {
		return "", err
	}

	return id, nil
}

const getUser = `
SELECT *
FROM users
WHERE (user_id = $1 OR email = $2) AND deleted_at IS NULL; 
`

func (s *Storage) GetUser(ctx context.Context, user storage.User) (*storage.User, error) {
	var res storage.User
	if err := s.db.Get(&res, getUser, user.UserID, user.Email); err != nil {
		return nil, fmt.Errorf("executing user details: %w", err)
	}
	return &res, nil
}

const deleteUser = `
UPDATE
	users
SET
	deleted_at = now(),
	deleted_by = $1
WHERE 
	user_id = $2;

`

func (s *Storage) DeleteUser(ctx context.Context, user storage.User) error {
	_, err := s.db.Exec(deleteUser, user.DeletedBy, user.UserID)
	if err != nil {
		return err
	}
	return nil
}

const updateUser = `
UPDATE
	 users
SET
	name = :name,
    batch = :batch,
	email = :email,
	password = :password,
	updated_at = now(),
	updated_by = :updated_by
WHERE 
	user_id = :user_id
RETURNING 
	updated_at;
`

func (s *Storage) UpdateUser(ctx context.Context, user storage.User) (*storage.User, error) {
	stmt, err := s.db.PrepareNamedContext(ctx, updateUser)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()
	if err := stmt.Get(&user, user); err != nil {
		if err == sql.ErrNoRows {
			return nil, err
		}
		return nil, fmt.Errorf("executing user update: %w", err)
	}
	return &user, nil
}

func (s *Storage) ListUser(ctx context.Context, f storage.Filter) ([]storage.User, error) {
	var users []storage.User
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

	listUser := fmt.Sprintf("SELECT * FROM users WHERE deleted_at IS NULL AND (name ILIKE '%%' || '%s' || '%%' OR email ILIKE '%%' || '%s' || '%%') ORDER BY %s %s", f.SearchTerm, f.SearchTerm, sortBy, order)
	fullQuery := listUser + limit
	if err := s.db.Select(&users, fullQuery); err != nil {
		if err == sql.ErrNoRows {
			return nil, err
		}
		return nil, err
	}

	return users, nil
}

func (s *Storage) UserStats(ctx context.Context, f storage.Filter) (storage.Stats, error) {
	var userStat = fmt.Sprintf("SELECT COUNT(*) FROM users where deleted_at IS NULL AND (name ILIKE '%%' || '%s' || '%%' OR email ILIKE '%%' || '%s' || '%%');", f.SearchTerm, f.SearchTerm)
	var stat storage.Stats
	if err := s.db.Get(&stat, userStat); err != nil {
		if err == sql.ErrNoRows {
			return storage.Stats{}, err
		}
		return storage.Stats{}, err
	}

	return stat, nil
}
