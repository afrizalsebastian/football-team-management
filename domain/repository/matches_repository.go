package repository

import (
	"context"
	"errors"

	"github.com/afrizalsebastian/football-team-management/application/helper"
	"github.com/afrizalsebastian/football-team-management/constants"
	"github.com/afrizalsebastian/football-team-management/domain/dao"
	matchdb "github.com/afrizalsebastian/football-team-management/domain/repository/queries/db"
	"github.com/afrizalsebastian/football-team-management/module/logger"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
)

type IMatchesRepository interface {
	CreateMatches(ctx context.Context, match *dao.Matches) error
	GetListMatch(ctx context.Context, teamId *string) ([]dao.Matches, error)
	RescheduleMatch(ctx context.Context, match *dao.Matches) error
	DeleteMatch(ctx context.Context, matchId string) error
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
		MatchTime:  TimeStrPtrToPgtypeTime(match.MatchTime),
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

func (r *matchesRepository) GetListMatch(ctx context.Context, teamId *string) ([]dao.Matches, error) {
	l := logger.LoggerNew()

	l.WithContext(ctx).Debug("[GetListMatch].domain: Started").Msg()
	var (
		homeId pgtype.UUID
		awayId pgtype.UUID
	)

	if teamId != nil {
		if err := homeId.Scan(*teamId); err != nil {
			return nil, constants.InvalidUUIDValue
		}
		if err := awayId.Scan(*teamId); err != nil {
			return nil, constants.InvalidUUIDValue
		}
	}

	rows, err := r.db.GetListMatch(ctx, &matchdb.GetListMatchParams{
		HomeTeamID: homeId,
		AwayTeamID: awayId,
	})
	if err != nil {
		l.WithContext(ctx).
			Error("Error when get list match").
			Attr("error", err).
			Msg()

		return nil, err
	}

	result := make([]dao.Matches, 0)
	for _, r := range rows {
		if r == nil {
			continue
		}

		result = append(result, dao.Matches{
			Id:         r.ID.String(),
			HomeTeamId: r.HomeTeamID.String(),
			AwayTeamId: r.AwayTeamID.String(),
			MatchDate:  helper.StringPtr(PgDateToDateStr(r.MatchDate)),
			MatchTime:  helper.StringPtr(PgTimeToTimeStr(r.MatchTime)),
			HomeTeam: &dao.Teams{
				Id:   r.HomeTeamID.String(),
				Name: helper.StringPtr(r.HomeTeamName.String),
				Logo: helper.StringPtr(r.HomeLogo.String),
			},
			AwayTeam: &dao.Teams{
				Id:   r.AwayTeamID.String(),
				Name: helper.StringPtr(r.AwayTeamName.String),
				Logo: helper.StringPtr(r.AwayLogo.String),
			},
		})
	}

	l.WithContext(ctx).Debug("[GetListMatch].domain: Completed").Msg()
	return result, nil
}

func (r *matchesRepository) RescheduleMatch(ctx context.Context, match *dao.Matches) error {
	l := logger.LoggerNew()

	l.WithContext(ctx).Debug("[RescheduleMatch].domain: Started").Msg()

	var matchId pgtype.UUID
	if err := matchId.Scan(match.Id); err != nil {
		return constants.InvalidUUIDValue
	}

	if _, err := r.db.RescheduleMatch(ctx, &matchdb.RescheduleMatchParams{
		MatchDate: DateStrPtrToPgtypeDate(match.MatchDate),
		MatchTime: TimeStrPtrToPgtypeTime(match.MatchTime),
		ID:        matchId,
	}); err != nil {
		l.WithContext(ctx).
			Error("Error when get list match").
			Attr("error", err).
			Msg()

		if errors.Is(err, pgx.ErrNoRows) {
			return constants.ErrNotFoundRow
		}

		return err
	}

	l.WithContext(ctx).Debug("[RescheduleMatch].domain: Completed").Msg()
	return nil
}

func (r *matchesRepository) DeleteMatch(ctx context.Context, matchIdStr string) error {
	l := logger.LoggerNew()

	l.WithContext(ctx).Debug("[DeleteMatch].domain: Started").Msg()

	var matchId pgtype.UUID
	if err := matchId.Scan(matchIdStr); err != nil {
		return constants.InvalidUUIDValue
	}

	if _, err := r.db.DeleteMatch(ctx, matchId); err != nil {
		l.WithContext(ctx).
			Error("Error when delete match").
			Attr("error", err).
			Msg()

		if errors.Is(err, pgx.ErrNoRows) {
			return constants.ErrNotFoundRow
		}

		return err
	}

	l.WithContext(ctx).Debug("[DeleteMatch].domain: Completed").Msg()
	return nil
}
