package service

import "errors"

var (
	ErrNotFound      = errors.New("no rows in result set")
	ErrAlreadyExists = errors.New("unique constraint violation")
	ErrBadInput      = errors.New("foreign key violation")
	Err              = errors.New("application error")
)
