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
	GetListMatch(ctx context.Context, teamId *string, matchType *string) ([]dao.Matches, error)
	RescheduleMatch(ctx context.Context, match *dao.Matches) error
	DeleteMatch(ctx context.Context, matchId string) error
	MatchFullTime(ctx context.Context, matchId string) error
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
	admin := helper.GetContextValueClaims(ctx)
	if admin == nil {
		return constants.AdminContextNil
	}

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
		CreatedBy:  StringToPgtypeText(admin.Username),
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

func (r *matchesRepository) GetListMatch(ctx context.Context, teamId *string, matchType *string) ([]dao.Matches, error) {
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
		FilterType: StringPtrToPgtypeText(matchType),
	})
	if err != nil {
		l.WithContext(ctx).
			Error("Error when get list match").
			Attr("error", err).
			Msg()

		return nil, err
	}

	mapMatch := make(map[string]bool)
	match := make([]dao.Matches, 0)
	mapGoals := make(map[string][]dao.Goals)
	mapTeamGoals := make(map[string]map[string]int)

	for _, r := range rows {
		if r == nil {
			continue
		}

		rowMatchId := r.ID.String()
		if _, ok := mapMatch[rowMatchId]; !ok {
			match = append(match, dao.Matches{
				Id:         rowMatchId,
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
				Status:    int(r.Status.Int16),
				HomeScore: int(r.HomeScore.Int16),
				AwayScore: int(r.AwayScore.Int16),
			})

			mapMatch[rowMatchId] = true
		}

		if r.GoalID.Valid {
			playerTeamId := r.PlayerTeamID.String()
			mapGoals[rowMatchId] = append(mapGoals[rowMatchId], dao.Goals{
				Id:         r.GoalID.String(),
				GoalMinute: r.GoalMinute.String,
				PlayerId:   r.PlayerID.String(),
				Player: &dao.Players{
					Id:           r.PlayerID.String(),
					Name:         helper.StringPtr(r.PlayerName.String),
					JerseyNumber: helper.IntPtr(int(r.PlayerNumber.Int16)),
					TeamId:       playerTeamId,
				},
			})

			// live score
			_, ok := mapTeamGoals[rowMatchId]
			if !ok {
				mapTeamGoals[rowMatchId] = map[string]int{
					playerTeamId: 1,
				}
				continue
			}

			_, ok = mapTeamGoals[rowMatchId][playerTeamId]
			if !ok {
				mapTeamGoals[rowMatchId][playerTeamId] = 1
				continue
			}

			mapTeamGoals[rowMatchId][playerTeamId] += 1
		}
	}

	for i := range match {
		val, _ := mapGoals[match[i].Id]
		match[i].Goals = val

		// RESULT OF SCORE BASED ON MATCH RESULT
		// THIS IS FOR LIVESCORE WHEN MATCH RESULT DOESN'T MADE YET
		if match[i].AwayScore == 0 {
			val, ok := mapTeamGoals[match[i].Id][match[i].AwayTeamId]
			if ok {
				match[i].AwayScore = val
			}
		}

		if match[i].HomeScore == 0 {
			val, ok := mapTeamGoals[match[i].Id][match[i].HomeTeamId]
			if ok {
				match[i].HomeScore = val
			}
		}

		switch {
		case match[i].AwayScore == match[i].HomeScore:
			match[i].Status = 0

		case match[i].AwayScore > match[i].HomeScore:
			match[i].Status = -1

		case match[i].AwayScore < match[i].HomeScore:
			match[i].Status = 1
		}
		// ====================================================
	}
	l.WithContext(ctx).Debug("[GetListMatch].domain: Completed").Msg()
	return match, nil
}

func (r *matchesRepository) RescheduleMatch(ctx context.Context, match *dao.Matches) error {
	l := logger.LoggerNew()

	l.WithContext(ctx).Debug("[RescheduleMatch].domain: Started").Msg()
	admin := helper.GetContextValueClaims(ctx)
	if admin == nil {
		return constants.AdminContextNil
	}

	var matchId pgtype.UUID
	if err := matchId.Scan(match.Id); err != nil {
		return constants.InvalidUUIDValue
	}

	if _, err := r.db.RescheduleMatch(ctx, &matchdb.RescheduleMatchParams{
		MatchDate: DateStrPtrToPgtypeDate(match.MatchDate),
		MatchTime: TimeStrPtrToPgtypeTime(match.MatchTime),
		ID:        matchId,
		UpdatedBy: StringToPgtypeText(admin.Username),
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
	admin := helper.GetContextValueClaims(ctx)
	if admin == nil {
		return constants.AdminContextNil
	}

	var matchId pgtype.UUID
	if err := matchId.Scan(matchIdStr); err != nil {
		return constants.InvalidUUIDValue
	}

	if _, err := r.db.DeleteMatch(ctx, &matchdb.DeleteMatchParams{
		ID:        matchId,
		DeletedBy: StringToPgtypeText(admin.Username),
	}); err != nil {
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

func (r *matchesRepository) MatchFullTime(ctx context.Context, matchIdStr string) error {
	l := logger.LoggerNew()

	l.WithContext(ctx).Debug("[DeleteMatch].domain: Started").Msg()
	admin := helper.GetContextValueClaims(ctx)
	if admin == nil {
		return constants.AdminContextNil
	}

	var matchId pgtype.UUID
	if err := matchId.Scan(matchIdStr); err != nil {
		return constants.InvalidUUIDValue
	}

	if _, err := r.db.MatchFullTime(ctx, &matchdb.MatchFullTimeParams{
		MatchID:   matchId,
		CreatedBy: StringToPgtypeText(admin.Username),
	}); err != nil {
		l.WithContext(ctx).
			Error("Error when set match result").
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
