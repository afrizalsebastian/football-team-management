package repository

import (
	"errors"

	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgtype"
)

func isDuplicateError(err error) bool {
	var pgErr *pgconn.PgError

	if errors.As(err, &pgErr) {
		return pgErr.Code == "23505"
	}

	return false
}

func StringPtrToPgtypeText(s *string) pgtype.Text {
	if s == nil || *s == "" {
		return pgtype.Text{Valid: false}
	}

	return pgtype.Text{String: *s, Valid: true}
}

func StringToPgtypeText(s string) pgtype.Text {
	if s == "" {
		return pgtype.Text{Valid: false}
	}

	return pgtype.Text{String: s, Valid: true}
}
