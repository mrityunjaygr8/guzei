package main

import (
	"errors"
	"time"
)

type GuzeiStore interface {
	UserInsert(email, password string, admin bool) (*User, error)
}

type User struct {
	Email     string
	ID        string
	Admin     bool
	CreatedAt time.Time
	UpdatedAt time.Time
}

var ErrUserExists = errors.New("user with specified email already exists")
