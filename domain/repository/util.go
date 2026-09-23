package repository

import (
	"errors"
	"time"

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

func DateStrPtrToPgtypeDate(s *string) pgtype.Date {
	if s == nil {
		return pgtype.Date{Valid: false}
	}

	t, err := time.Parse("02-01-2006", *s)
	if err != nil {
		return pgtype.Date{Valid: false}
	}

	return pgtype.Date{
		Time:  t,
		Valid: true,
	}
}

func DateStrToPgtypeDate(s string) pgtype.Date {
	if s == "" {
		return pgtype.Date{Valid: false}
	}

	t, err := time.Parse("02-01-2006", s)
	if err != nil {
		return pgtype.Date{Valid: false}
	}

	return pgtype.Date{
		Time:  t,
		Valid: true,
	}
}

func TimeStrPtrToPgtypeDate(s *string) pgtype.Time {
	if s == nil {
		return pgtype.Time{Valid: false}
	}

	t, err := time.Parse("15:04", *s)
	if err != nil {
		return pgtype.Time{Valid: false}
	}

	duration := time.Duration(t.Hour())*time.Hour + time.Duration(t.Minute())*time.Minute
	microseconds := int64(duration / time.Microsecond)
	return pgtype.Time{
		Microseconds: microseconds,
		Valid:        true,
	}
}

func TimeStrToPgtypeDate(s string) pgtype.Time {
	if s == "" {
		return pgtype.Time{Valid: false}
	}

	t, err := time.Parse("15:04", s)
	if err != nil {
		return pgtype.Time{Valid: false}
	}

	duration := time.Duration(t.Hour())*time.Hour + time.Duration(t.Minute())*time.Minute
	microseconds := int64(duration / time.Microsecond)
	return pgtype.Time{
		Microseconds: microseconds,
		Valid:        true,
	}
}
