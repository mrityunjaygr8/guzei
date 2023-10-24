package store

import (
	"errors"
	"time"
)

type GuzeiStore interface {
	UserInsert(email, password, id string, admin bool) (*User, error)
	UserList(pageSize, pageNumber int) (*UsersList, error)
	UserRetrieve(id string) (*User, error)
	UserUpdatePassword(id, newPassword string) error
	UserUpdateAdmin(id string, newAdminValue bool) error
	UserDelete(id string) error
}

type User struct {
	Email     string    `json:"email"`
	ID        string    `json:"id"`
	Admin     bool      `json:"admin"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type UsersList struct {
	Data         []User `json:"data"`
	TotalObjects int    `json:"total"`
	TotalPages   int    `json:"pages"`
}

var ErrUserExists = errors.New("user with specified email already exists")
var ErrUserNotFound = errors.New("specified user does not exists")
var ErrStoreError = errors.New("error persisting in storage")
