package main

import (
	"fmt"
	"github.com/stretchr/testify/require"
	"os"
	"testing"
)

var TestDBString = fmt.Sprintf("postgres://%s:%s@%s:%s/%s", os.Getenv("DB_USER"), os.Getenv("DB_PASS"), os.Getenv("DB_HOST"), os.Getenv("DB_PORT"), os.Getenv("DB_NAME"))

func TestPostgresStore(t *testing.T) {
	store, err := NewPostgresStore(TestDBString)
	require.Nil(t, err)
	t.Run("test UserInsert method happy path", func(t *testing.T) {
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
		email := "im@parham.im"
		password := "password"
		admin := true

		_, _ = store.UserInsert(email, password, admin)
		_, err := store.UserInsert(email, password, admin)

		require.Error(t, err)
		require.Equal(t, ErrUserExists, err)

	})

	t.Run("test UserList method happy path", func(t *testing.T) {
		email := "im@parham.im123"
		password := "password"
		admin := true

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
