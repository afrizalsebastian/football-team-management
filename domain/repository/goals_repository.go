package repository

import (
	"context"
	"errors"

	"github.com/afrizalsebastian/football-team-management/application/helper"
	"github.com/afrizalsebastian/football-team-management/constants"
	"github.com/afrizalsebastian/football-team-management/domain/dao"
	goaldb "github.com/afrizalsebastian/football-team-management/domain/repository/queries/db"
	"github.com/afrizalsebastian/football-team-management/module/logger"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
)

type IGoalsRepository interface {
	MatchGoal(ctx context.Context, goals *dao.Goals) error
	GetMatchGoalList(ctx context.Context, matchId string) ([]dao.Goals, *dao.Matches, error)
	DeleteGoals(ctx context.Context, goalId string) error
}

type goalsRepository struct {
	db *goaldb.Queries
}

func NewGoalsRepository(
	pool *pgxpool.Pool,
) IGoalsRepository {
	return &goalsRepository{
		db: goaldb.New(pool),
	}
}

func (r *goalsRepository) MatchGoal(ctx context.Context, goals *dao.Goals) error {
	l := logger.LoggerNew()

	l.WithContext(ctx).Debug("[MatchGoal].domain: Started").Msg()

	var (
		matchId  pgtype.UUID
		playerId pgtype.UUID
	)

	if err := matchId.Scan(goals.MatchId); err != nil {
		return constants.InvalidUUIDValue
	}
	if err := playerId.Scan(goals.PlayerId); err != nil {
		return constants.InvalidUUIDValue
	}

	id, err := r.db.CreateGoals(ctx, &goaldb.CreateGoalsParams{
		MatchID:    matchId,
		PlayerID:   playerId,
		GoalMinute: StringToPgtypeText(goals.GoalMinute),
	})

	if err != nil {
		l.WithContext(ctx).
			Error("error when create match").
			Attr("match_id", goals.MatchId).
			Attr("player_id", goals.PlayerId).
			Msg()

		if errors.Is(err, pgx.ErrNoRows) {
			return constants.ErrNotFoundRow
		}
		return err
	}

	goals.Id = id.String()

	l.WithContext(ctx).Debug("[MatchGoal].domain: Completed").Msg()
	return nil
}

func (r *goalsRepository) GetMatchGoalList(ctx context.Context, matchIdStr string) ([]dao.Goals, *dao.Matches, error) {
	l := logger.LoggerNew()

	l.WithContext(ctx).Debug("[GetListMatchGoal].domain: Started").Msg()

	var (
		matchId pgtype.UUID
	)

	if err := matchId.Scan(matchIdStr); err != nil {
		return nil, nil, constants.InvalidUUIDValue
	}

	rows, err := r.db.GetListMatchGoals(ctx, matchId)
	if err != nil {
		if err != nil {
			l.WithContext(ctx).
				Error("error when get list match goals").
				Attr("match_id", matchIdStr).
				Msg()

			return nil, nil, err
		}
	}

	result := make([]dao.Goals, 0)
	match := &dao.Matches{}
	for i, r := range rows {
		if i == 0 { // populate match only in first index
			match.Id = r.MatchID.String()
			match.HomeTeam = &dao.Teams{
				Id:   r.HomeTeamID.String(),
				Name: helper.StringPtr(r.HomeTeamName.String),
			}
			match.AwayTeam = &dao.Teams{
				Id:   r.AwayTeamID.String(),
				Name: helper.StringPtr(r.AwayTeamName.String),
			}
		}

		result = append(result, dao.Goals{
			Id:       r.ID.String(),
			MatchId:  r.MatchID.String(),
			PlayerId: r.PlayerID.String(),
			Player: &dao.Players{
				Id:           r.PlayerID.String(),
				Name:         helper.StringPtr(r.PlayerName),
				JerseyNumber: helper.IntPtr(int(r.JerseyNumber)),
				Team: &dao.Teams{
					Id:   r.TeamID.String(),
					Name: helper.StringPtr(r.TeamName),
				},
			},
			GoalMinute: r.GoalMinute.String,
		})
	}

	l.WithContext(ctx).Debug("[GetListMatchGoal].domain: Completed").Msg()
	return result, match, nil
}

func (r *goalsRepository) DeleteGoals(ctx context.Context, goalIdStr string) error {
	l := logger.LoggerNew()

	l.WithContext(ctx).Debug("[DeleteGoals].domain: Started").Msg()

	var (
		goalId pgtype.UUID
	)

	if err := goalId.Scan(goalIdStr); err != nil {
		return constants.InvalidUUIDValue
	}

	if _, err := r.db.DeleteGoals(ctx, goalId); err != nil {
		l.WithContext(ctx).
			Error("Error when delete goal").
			Attr("error", err).
			Msg()

		if errors.Is(err, pgx.ErrNoRows) {
			return constants.ErrNotFoundRow
		}

		return err
	}

	l.WithContext(ctx).Debug("[DeleteGoals].domain: Completed").Msg()
	return nil
}
