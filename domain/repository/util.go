package repository

import (
	"errors"

	"github.com/afrizalsebastian/football-team-management/domain/repository/queries/db"
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
	if s == nil {
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

func IntToPgTypeInt2(v int) pgtype.Int2 {
	if v == 0 {
		return pgtype.Int2{Valid: false}
	}

	return pgtype.Int2{Int16: int16(v), Valid: true}
}

func IntPtrToPgTypeInt2(v *int) pgtype.Int2 {
	if v == nil {
		return pgtype.Int2{Valid: false}
	}

	return pgtype.Int2{Int16: int16(*v), Valid: true}
}

func PositionPtrToNullPlayerPosition(p *string) db.NullPlayerPosition {
	if p == nil || *p == "" {
		return db.NullPlayerPosition{Valid: false}
	}

	return db.NullPlayerPosition{PlayerPosition: db.PlayerPosition(*p), Valid: true}
}

func PositionToNullPlayerPosition(p string) db.NullPlayerPosition {
	if p == "" {
		return db.NullPlayerPosition{Valid: false}
	}

	return db.NullPlayerPosition{PlayerPosition: db.PlayerPosition(p), Valid: true}
}
