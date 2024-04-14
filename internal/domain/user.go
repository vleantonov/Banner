package domain

import (
	"errors"
)

var (
	ErrUserNotFound       = errors.New("user not found")
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrUserLoginPassword  = errors.New("login and password must be provided")
	ErrUserExists         = errors.New("user exists")
	ErrInvalidToken       = errors.New("invalid token")
)

type User struct {
	Login    string `db:"login"`
	PassHash []byte `db:"pass_hash"`
	IsAdmin  bool   `db:"is_admin"`
}
