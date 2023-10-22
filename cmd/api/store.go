package main

import (
	"errors"
	"time"
)

type GuzeiStore interface {
	UserInsert(email, password, id string, admin bool) (*User, error)
	UserList(pageSize, pageNumber int) (*UsersList, error)
	UserRetrieve(email string) (*User, error)
}

type User struct {
	Email     string
	ID        string
	Admin     bool
	CreatedAt time.Time
	UpdatedAt time.Time
}

type UsersList struct {
	data         []User
	totalObjects int
	totalPages   int
}

var ErrUserExists = errors.New("user with specified email already exists")
var ErrUserNotFound = errors.New("specified user does not exists")
