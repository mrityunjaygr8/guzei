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

func TestPostgresStoreUserInsert(t *testing.T) {
	t.Run("test UserInsert method happy path", func(t *testing.T) {
		store, teardownTest := setupTest(t)
		defer teardownTest(t)

		email := "im@parham.im"
		password := "password"
		id := "123"
		admin := true

		user, err := store.UserInsert(email, password, id, admin)

		fmt.Println(err)
		require.Nil(t, err)

		require.Equal(t, email, user.Email)
		require.Equal(t, admin, user.Admin)
		require.Equal(t, id, user.ID)
		require.NotNil(t, user.CreatedAt)
		require.Zero(t, user.UpdatedAt)
	})

	t.Run("test UserInsert for duplicates", func(t *testing.T) {
		store, teardownTest := setupTest(t)
		defer teardownTest(t)

		email := "im@parham.im"
		password := "password"
		id := "123"
		admin := true

		_, _ = store.UserInsert(email, password, id, admin)
		_, err := store.UserInsert(email, password, id, admin)

		require.Error(t, err)
		require.Equal(t, ErrUserExists, err)

	})
}

func TestPostgresStoreUserList(t *testing.T) {
	t.Run("test UserList method happy path", func(t *testing.T) {
		store, teardownTest := setupTest(t)
		defer teardownTest(t)

		email := "im@parham.im123"
		password := "password"
		admin := true
		id := "123"

		_, _ = store.UserInsert("a"+email, password, "1", admin)
		_, _ = store.UserInsert(email, password, id, admin)
		_, _ = store.UserInsert(email+"a", password, "3", admin)
		_, _ = store.UserInsert(email+"b", password, "4", admin)
		_, _ = store.UserInsert(email+"c", password, "5", admin)
		_, _ = store.UserInsert(email+"d", password, "6", admin)
		_, _ = store.UserInsert(email+"e", password, "7", admin)
		_, _ = store.UserInsert(email+"f", password, "8", admin)
		_, _ = store.UserInsert(email+"g", password, "9", admin)
		_, _ = store.UserInsert(email+"h", password, "10", admin)
		_, _ = store.UserInsert(email+"i", password, "11", admin)
		_, _ = store.UserInsert(email+"j", password, "12", admin)
		_, _ = store.UserInsert(email+"k", password, "13", admin)
		_, _ = store.UserInsert(email+"l", password, "14", admin)
		users, err := store.UserList(1, 10)
		require.Nil(t, err)
		require.NotNil(t, users)

		require.Equal(t, 14, users.totalObjects)
		require.Equal(t, 2, users.totalPages)

		require.Equal(t, email, users.data[1].Email)
		require.Equal(t, admin, users.data[1].Admin)
		require.Equal(t, id, users.data[1].ID)
		require.NotNil(t, users.data[1].CreatedAt)
		require.Zero(t, users.data[1].UpdatedAt)

	})

	t.Run("test UserList method no data", func(t *testing.T) {
		store, teardownTest := setupTest(t)
		defer teardownTest(t)

		users, err := store.UserList(1, 10)
		require.Nil(t, err)
		require.NotNil(t, users)

		require.Equal(t, 0, users.totalObjects)
		require.Equal(t, 0, users.totalPages)
	})

}

