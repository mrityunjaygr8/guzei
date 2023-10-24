package postgres_store

import (
	"context"
	"errors"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/mrityunjaygr8/guzei/internal/db/models"
	"github.com/mrityunjaygr8/guzei/store"
	"math"
)

type PostgresStore struct {
	db *pgxpool.Pool
}

var ErrCreatingPostgresPool = errors.New("error creating postgres pool")
var ErrConnectingToPostgres = errors.New("error connecting to postgres")

func NewPostgresStore(dbString string) (*PostgresStore, error) {
	db, err := pgxpool.New(context.Background(), dbString)
	if err != nil {
		return nil, ErrCreatingPostgresPool
	}

	err = db.Ping(context.Background())
	if err != nil {
		return nil, ErrConnectingToPostgres
	}

	return &PostgresStore{db: db}, nil
}

func (p *PostgresStore) UserInsert(email, password, id string, admin bool) (*store.User, error) {
	tx, err := p.db.BeginTx(context.Background(), pgx.TxOptions{})
	if err != nil {
		return nil, errors.New("error creating database transaction")
	}

	query := models.New(tx)
	params := models.UserInsertParams{
		Email:    email,
		Password: password,
		ID:       id,
		Admin:    admin,
	}
	dbUser, err := query.UserInsert(context.Background(), params)
	if err != nil {
		tx.Rollback(context.Background())
		var pge *pgconn.PgError
		if errors.As(err, &pge) {
			if pge.SQLState() == "23505" {
				return nil, store.ErrUserExists
			}
		}
		return nil, err
	}

	tx.Commit(context.Background())
	user := &store.User{
		Email:     dbUser.Email,
		ID:        dbUser.ID,
		Admin:     dbUser.Admin,
		CreatedAt: dbUser.CreatedAt.Time,
		UpdatedAt: dbUser.UpdatedAt.Time,
	}
	return user, nil
}

func (p *PostgresStore) UserList(pageNumber, pageSize int) (*store.UsersList, error) {
	query := models.New(p.db)
	params := models.UsersListParams{
		Limit:  int32(pageSize),
		Offset: int32((pageNumber - 1) * pageSize),
	}
	dbUsers, err := query.UsersList(context.Background(), params)
	if err != nil {
		return nil, err
	}

	users := make([]store.User, 0)
	totalObjects := 0
	totalPages := 0

	for _, user := range dbUsers {
		users = append(users, store.User{
			Email:     user.Email,
			ID:        user.ID,
			Admin:     user.Admin,
			CreatedAt: user.CreatedAt.Time,
			UpdatedAt: user.UpdatedAt.Time,
		})
	}

	if len(dbUsers) > 0 {
		totalObjects = int(dbUsers[0].RowData)
		totalPages = int(math.Ceil(float64(totalObjects) / float64(pageSize)))
	}

	return &store.UsersList{Data: users, TotalObjects: totalObjects, TotalPages: totalPages}, nil
}

func (p *PostgresStore) UserRetrieve(id string) (*store.User, error) {
	query := models.New(p.db)
	user, err := query.UserRetrieve(context.Background(), id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, store.ErrUserNotFound
		}
		return nil, err
	}

	return &store.User{
		Email:     user.Email,
		ID:        user.ID,
		Admin:     user.Admin,
		CreatedAt: user.CreatedAt.Time,
		UpdatedAt: user.UpdatedAt.Time,
	}, nil
}

func (p *PostgresStore) UserUpdatePassword(id, newPassword string) error {
	query := models.New(p.db)
	params := models.UserUpdatePasswordParams{
		ID:       id,
		Password: newPassword,
	}
	res, err := query.UserUpdatePassword(context.Background(), params)
	if err != nil {
		return err
	}
	if res.RowsAffected() == 0 {
		return store.ErrUserNotFound
	}
	return nil
}
func (p *PostgresStore) UserUpdateAdmin(id string, newAdminValue bool) error {
	query := models.New(p.db)
	params := models.UserUpdateAdminParams{
		ID:    id,
		Admin: newAdminValue,
	}
	res, err := query.UserUpdateAdmin(context.Background(), params)
	if err != nil {
		return err
	}
	if res.RowsAffected() == 0 {
		return store.ErrUserNotFound
	}
	return nil
}
