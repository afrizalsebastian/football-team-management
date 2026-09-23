package repository

import (
	"context"

	"github.com/afrizalsebastian/football-team-management/constants"
	"github.com/afrizalsebastian/football-team-management/domain/dao"
	matchdb "github.com/afrizalsebastian/football-team-management/domain/repository/queries/db"
	"github.com/afrizalsebastian/football-team-management/module/logger"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
)

type IMatchesRepository interface {
	CreateMatches(ctx context.Context, match *dao.Matches) error
}

type matchesRepository struct {
	db *matchdb.Queries
}

func NewMatchesRepository(
	pool *pgxpool.Pool,
) IMatchesRepository {
	return &matchesRepository{
		db: matchdb.New(pool),
	}
}

func (r *matchesRepository) CreateMatches(ctx context.Context, match *dao.Matches) error {
	l := logger.LoggerNew()

	l.WithContext(ctx).Debug("[CreateMatches].domain: Started").Msg()

	var (
		homeId pgtype.UUID
		awayId pgtype.UUID
	)

	if err := homeId.Scan(match.HomeTeamId); err != nil {
		return constants.InvalidUUIDValue
	}
	if err := awayId.Scan(match.AwayTeamId); err != nil {
		return constants.InvalidUUIDValue
	}

	params := &matchdb.CreateMatchesParams{
		MatchDate:  DateStrPtrToPgtypeDate(match.MatchDate),
		MatchTime:  TimeStrPtrToPgtypeDate(match.MatchTime),
		HomeTeamID: homeId,
		AwayTeamID: awayId,
	}

	row, err := r.db.CreateMatches(ctx, params)
	if err != nil {
		l.WithContext(ctx).
			Error("error when create match").
			Attr("home", match.HomeTeamId).
			Attr("away", match.AwayTeamId).
			Msg()

		return err
	}

	match.Id = row.String()
	l.WithContext(ctx).Debug("[CreateMatches].domain: Completed").Msg()
	return nil
}
