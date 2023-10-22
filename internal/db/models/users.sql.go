// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.22.0
// source: users.sql

package models

import (
	"context"

	"github.com/jackc/pgx/v5/pgtype"
)

const userDelete = `-- name: UserDelete :exec
DELETE FROM users WHERE "ID" = $1
`

func (q *Queries) UserDelete(ctx context.Context, id string) error {
	_, err := q.db.Exec(ctx, userDelete, id)
	return err
}

const userInsert = `-- name: UserInsert :one
INSERT INTO users (email, password, "ID", admin) VALUES ($1, $2, $3, $4) RETURNING email, "createdAt", "updatedAt", "ID", admin
`

type UserInsertParams struct {
	Email    string
	Password string
	ID       string
	Admin    bool
}

type UserInsertRow struct {
	Email     string
	CreatedAt pgtype.Timestamptz
	UpdatedAt pgtype.Timestamptz
	ID        string
	Admin     bool
}

func (q *Queries) UserInsert(ctx context.Context, arg UserInsertParams) (UserInsertRow, error) {
	row := q.db.QueryRow(ctx, userInsert,
		arg.Email,
		arg.Password,
		arg.ID,
		arg.Admin,
	)
	var i UserInsertRow
	err := row.Scan(
		&i.Email,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.ID,
		&i.Admin,
	)
	return i, err
}

const userRetrieve = `-- name: UserRetrieve :one
SELECT email, "createdAt", "updatedAt", "ID", admin FROM users WHERE email = $1 LIMIT 1
`

type UserRetrieveRow struct {
	Email     string
	CreatedAt pgtype.Timestamptz
	UpdatedAt pgtype.Timestamptz
	ID        string
	Admin     bool
}

func (q *Queries) UserRetrieve(ctx context.Context, email string) (UserRetrieveRow, error) {
	row := q.db.QueryRow(ctx, userRetrieve, email)
	var i UserRetrieveRow
	err := row.Scan(
		&i.Email,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.ID,
		&i.Admin,
	)
	return i, err
}

const userUpdateAdmin = `-- name: UserUpdateAdmin :exec
UPDATE users SET admin = $2 WHERE "ID" = $1
`

type UserUpdateAdminParams struct {
	ID    string
	Admin bool
}

func (q *Queries) UserUpdateAdmin(ctx context.Context, arg UserUpdateAdminParams) error {
	_, err := q.db.Exec(ctx, userUpdateAdmin, arg.ID, arg.Admin)
	return err
}

const userUpdatePassword = `-- name: UserUpdatePassword :exec
UPDATE users SET password = $2 WHERE "ID" = $1
`

type UserUpdatePasswordParams struct {
	ID       string
	Password string
}

func (q *Queries) UserUpdatePassword(ctx context.Context, arg UserUpdatePasswordParams) error {
	_, err := q.db.Exec(ctx, userUpdatePassword, arg.ID, arg.Password)
	return err
}

const usersList = `-- name: UsersList :many
WITH row_data AS (
    SELECT email, "createdAt", "updatedAt", "ID", admin FROM users ORDER BY email LIMIT $1 OFFSET $2
) SELECT
      email, "createdAt", "updatedAt", "ID", admin,
      (SELECT COUNT(*) FROM users) AS row_data
FROM row_data
`

type UsersListParams struct {
	Limit  int32
	Offset int32
}

type UsersListRow struct {
	Email     string
	CreatedAt pgtype.Timestamptz
	UpdatedAt pgtype.Timestamptz
	ID        string
	Admin     bool
	RowData   int64
}

func (q *Queries) UsersList(ctx context.Context, arg UsersListParams) ([]UsersListRow, error) {
	rows, err := q.db.Query(ctx, usersList, arg.Limit, arg.Offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []UsersListRow
	for rows.Next() {
		var i UsersListRow
		if err := rows.Scan(
			&i.Email,
			&i.CreatedAt,
			&i.UpdatedAt,
			&i.ID,
			&i.Admin,
			&i.RowData,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}
