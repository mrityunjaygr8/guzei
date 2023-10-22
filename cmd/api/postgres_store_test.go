package main

import (
	"fmt"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/stretchr/testify/require"
	"log"
	"os"
	"testing"
)

var TestDBString = fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s", os.Getenv("DB_USER"), os.Getenv("DB_PASS"), os.Getenv("DB_HOST"), os.Getenv("DB_PORT"), os.Getenv("DB_TEST_NAME"), os.Getenv("DB_SSLMODE"))

func setupTest(t testing.TB) (*PostgresStore, func(tb testing.TB)) {
	t.Log("setup suite")
	t.Log(TestDBString)
	m, err := migrate.New("file://./../../internal/db/migrations", TestDBString)
	t.Log(err)
	require.Nil(t, err)
	err = m.Up()
	require.Nil(t, err)
	store, err := NewPostgresStore(TestDBString)
	require.Nil(t, err)
	return store, func(tb testing.TB) {
		defer store.db.Close()
		log.Println("teardown suite")
		m, err := migrate.New("file://./../../internal/db/migrations", TestDBString)
		require.Nil(t, err)
		err = m.Down()
		require.Nil(t, err)
	}
}

func TestPostgresStore(t *testing.T) {

	t.Run("test UserInsert method happy path", func(t *testing.T) {
		store, teardownTest := setupTest(t)
		defer teardownTest(t)

		email := "im@parham.im"
		password := "password"
		admin := true

		user, err := store.UserInsert(email, password, admin)

		fmt.Println(err)
		require.Nil(t, err)

		require.Equal(t, email, user.Email)
		require.Equal(t, admin, user.Admin)
		require.NotNil(t, user.ID)
		require.NotNil(t, user.CreatedAt)
		require.Zero(t, user.UpdatedAt)
	})

	t.Run("test UserInsert for duplicates", func(t *testing.T) {
		store, teardownTest := setupTest(t)
		defer teardownTest(t)

		email := "im@parham.im"
		password := "password"
		admin := true

		_, _ = store.UserInsert(email, password, admin)
		_, err := store.UserInsert(email, password, admin)

		require.Error(t, err)
		require.Equal(t, ErrUserExists, err)

	})

	t.Run("test UserList method happy path", func(t *testing.T) {
		store, teardownTest := setupTest(t)
		defer teardownTest(t)

		email := "im@parham.im123"
		password := "password"
		admin := true

		_, _ = store.UserInsert("a"+email, password, admin)
		_, _ = store.UserInsert(email, password, admin)
		users, err := store.UserList(1, 10)
		require.Nil(t, err)
		require.NotNil(t, users)

		require.Equal(t, 2, len(users))

		require.Equal(t, email, users[1].Email)
		require.Equal(t, admin, users[1].Admin)
		require.NotNil(t, users[1].ID)
		require.NotNil(t, users[1].CreatedAt)
		require.Zero(t, users[1].UpdatedAt)

	})

}
