package helper

import (
	"errors"

	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

// IsDuplicateError сверяет является ли переданная ошибка дублированием уникального ключа
func IsDuplicateError(err error) bool {
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		return pgErr.Code == pgerrcode.UniqueViolation
	}
	return false
}

func IsNoRows(err error) bool {
	return errors.As(err, &pgx.ErrNoRows)
}
