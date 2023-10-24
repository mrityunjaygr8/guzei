package postgres_store

import (
	"fmt"
	"github.com/mrityunjaygr8/guzei/store"
	"github.com/stretchr/testify/require"
	"log"
	"os"
	"testing"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

var TestDBString = fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s", os.Getenv("DB_USER"), os.Getenv("DB_PASS"), os.Getenv("DB_HOST"), os.Getenv("DB_PORT"), os.Getenv("DB_TEST_NAME"), os.Getenv("DB_SSLMODE"))

func setupTest(t testing.TB) (*PostgresStore, func(tb testing.TB)) {
	m, err := migrate.New("file://./../../internal/db/migrations", TestDBString)
	t.Log(TestDBString)
	require.Nil(t, err)
	err = m.Up()
	require.Nil(t, err)
	postgresStore, err := NewPostgresStore(TestDBString)
	require.Nil(t, err)
	return postgresStore, func(tb testing.TB) {
		defer postgresStore.db.Close()
		log.Println("teardown suite")
		m, err := migrate.New("file://./../../internal/db/migrations", TestDBString)
		require.Nil(t, err)
		err = m.Down()
		require.Nil(t, err)
	}
}

func TestPostgresStoreUserInsert(t *testing.T) {
	t.Run("test UserInsert method happy path", func(t *testing.T) {
		postgresStore, teardownTest := setupTest(t)
		defer teardownTest(t)

		email := "im@parham.im"
		password := "password"
		id := "123"
		admin := true

		user, err := postgresStore.UserInsert(email, password, id, admin)

		fmt.Println(err)
		require.Nil(t, err)

		require.Equal(t, email, user.Email)
		require.Equal(t, admin, user.Admin)
		require.Equal(t, id, user.ID)
		require.NotNil(t, user.CreatedAt)
		require.Zero(t, user.UpdatedAt)
	})

	t.Run("test UserInsert for duplicates", func(t *testing.T) {
		postgresStore, teardownTest := setupTest(t)
		defer teardownTest(t)

		email := "im@parham.im"
		password := "password"
		id := "123"
		admin := true

		_, _ = postgresStore.UserInsert(email, password, id, admin)
		_, err := postgresStore.UserInsert(email, password, id, admin)

		require.Error(t, err)
		require.Equal(t, store.ErrUserExists, err)
	})
}

func TestPostgresStoreUserList(t *testing.T) {
	t.Run("test UserList method happy path", func(t *testing.T) {
		postgresStore, teardownTest := setupTest(t)
		defer teardownTest(t)

		email := "im@parham.im123"
		password := "password"
		admin := true
		id := "123"

		_, _ = postgresStore.UserInsert("a"+email, password, "1", admin)
		_, _ = postgresStore.UserInsert(email, password, id, admin)
		_, _ = postgresStore.UserInsert(email+"a", password, "3", admin)
		_, _ = postgresStore.UserInsert(email+"b", password, "4", admin)
		_, _ = postgresStore.UserInsert(email+"c", password, "5", admin)
		_, _ = postgresStore.UserInsert(email+"d", password, "6", admin)
		_, _ = postgresStore.UserInsert(email+"e", password, "7", admin)
		_, _ = postgresStore.UserInsert(email+"f", password, "8", admin)
		_, _ = postgresStore.UserInsert(email+"g", password, "9", admin)
		_, _ = postgresStore.UserInsert(email+"h", password, "10", admin)
		_, _ = postgresStore.UserInsert(email+"i", password, "11", admin)
		_, _ = postgresStore.UserInsert(email+"j", password, "12", admin)
		_, _ = postgresStore.UserInsert(email+"k", password, "13", admin)
		_, _ = postgresStore.UserInsert(email+"l", password, "14", admin)
		users, err := postgresStore.UserList(1, 10)
		require.Nil(t, err)
		require.NotNil(t, users)

		require.Equal(t, 14, users.TotalObjects)
		require.Equal(t, 2, users.TotalPages)

		require.Equal(t, email, users.Data[1].Email)
		require.Equal(t, admin, users.Data[1].Admin)
		require.Equal(t, id, users.Data[1].ID)
		require.NotNil(t, users.Data[1].CreatedAt)
		require.Zero(t, users.Data[1].UpdatedAt)
	})

	t.Run("test UserList method no data", func(t *testing.T) {
		postgresStore, teardownTest := setupTest(t)
		defer teardownTest(t)

		users, err := postgresStore.UserList(1, 10)
		require.Nil(t, err)
		require.NotNil(t, users)

		require.Equal(t, 0, users.TotalObjects)
		require.Equal(t, 0, users.TotalPages)
	})
}

func TestPostgresStoreUserRetrieve(t *testing.T) {
	t.Run("UserRetrieve happy path", func(t *testing.T) {
		postgresStore, teardownTest := setupTest(t)
		defer teardownTest(t)

		email := "im@parham.im"
		password := "password"
		admin := true
		id := "123"

		user, err := postgresStore.UserInsert(email, password, id, admin)
		require.Nil(t, err)
		require.NotNil(t, user)

		retrieved, err := postgresStore.UserRetrieve(user.ID)
		t.Log(retrieved)
		require.Nil(t, err)
		require.NotNil(t, retrieved)

		require.Equal(t, user, retrieved)
	})

	t.Run("UserRetrieve not exists", func(t *testing.T) {
		postgresStore, teardownTest := setupTest(t)
		defer teardownTest(t)

		id := "123"

		retrieved, err := postgresStore.UserRetrieve(id)
		require.Nil(t, retrieved)
		require.NotNil(t, err)
		t.Log(err)
	})
}

func TestNewPostgresStoreUserUpdatePassword(t *testing.T) {
	t.Run("UserUpdatePassword happy path", func(t *testing.T) {
		postgresStore, teardownTest := setupTest(t)
		defer teardownTest(t)

		email := "im@parham.im123"
		password := "password"
		newPassword := "newPassword"
		admin := true
		id := "123"

		user, err := postgresStore.UserInsert(email, password, id, admin)
		require.Nil(t, err)
		require.NotNil(t, user)

		err = postgresStore.UserUpdatePassword(user.ID, newPassword)
		require.Nil(t, err)
	})
	t.Run("UserUpdatePassword user not exists", func(t *testing.T) {
		postgresStore, teardownTest := setupTest(t)
		defer teardownTest(t)

		id := "1234"
		newPassword := "newPassword"

		err := postgresStore.UserUpdatePassword(id, newPassword)
		t.Log(err)
		require.NotNil(t, err)
		require.Equal(t, store.ErrUserNotFound, err)
	})
}
func TestNewPostgresStoreUserUpdateAdmin(t *testing.T) {
	t.Run("UserUpdateAdmin happy path", func(t *testing.T) {
		postgresStore, teardownTest := setupTest(t)
		defer teardownTest(t)

		email := "im@parham.im123"
		password := "password"
		admin := true
		id := "123"
		newAdminValue := false

		user, err := postgresStore.UserInsert(email, password, id, admin)
		require.Nil(t, err)
		require.NotNil(t, user)

		err = postgresStore.UserUpdateAdmin(user.ID, newAdminValue)
		require.Nil(t, err)

		u, err := postgresStore.UserRetrieve(user.ID)
		require.Nil(t, err)
		require.NotNil(t, u)

		require.Equal(t, newAdminValue, u.Admin)

	})
	t.Run("UserUpdateAdmin user not exists", func(t *testing.T) {
		postgresStore, teardownTest := setupTest(t)
		defer teardownTest(t)

		id := "1234"
		newAdminValue := false

		err := postgresStore.UserUpdateAdmin(id, newAdminValue)
		t.Log(err)
		require.NotNil(t, err)
		require.Equal(t, store.ErrUserNotFound, err)
	})
}
func TestNewPostgresStoreUserDelete(t *testing.T) {
	t.Run("UserDelete happy path", func(t *testing.T) {
		postgresStore, teardownTest := setupTest(t)
		defer teardownTest(t)

		email := "im@parham.im123"
		password := "password"
		admin := true
		id := "123"

		user, err := postgresStore.UserInsert(email, password, id, admin)
		require.Nil(t, err)
		require.NotNil(t, user)

		err = postgresStore.UserDelete(id)
		require.Nil(t, err)
	})
	t.Run("UserUpdateAdmin user not exists", func(t *testing.T) {
		postgresStore, teardownTest := setupTest(t)
		defer teardownTest(t)

		id := "1234"

		err := postgresStore.UserDelete(id)
		require.NotNil(t, err)
		require.Equal(t, store.ErrUserNotFound, err)
	})
}
