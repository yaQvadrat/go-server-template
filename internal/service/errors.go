package service

import "errors"

var (
	ErrAccountAlreadyExists = errors.New("account already exists")
	ErrAccountNotFound      = errors.New("account not found")
	ErrCannotGetAccount     = errors.New("cannot get account")
	ErrCannotSignJWT        = errors.New("cannot sign JWT")
	ErrCannotParseJWT       = errors.New("cannot parse JWT")
)
