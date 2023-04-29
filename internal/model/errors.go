package model

import (
	"errors"
	"github.com/jackc/pgx/v5/pgconn"
)

var ErrDuplicateKey = errors.New("duplicate key not allowed")
var ErrRecordNotFound = errors.New("record not found")
var ErrDuplicateEmail = errors.New("duplicate email not allowed")

func IsDuplicateError(err error) bool {
	if err == nil {
		return false
	}
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		return pgErr.Code == "23505"
	}
	return false
}
