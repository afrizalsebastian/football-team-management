package repository

import (
	"context"
	"fmt"

	"github.com/afrizalsebastian/football-team-management/application/helper"
	"github.com/afrizalsebastian/football-team-management/constants"
	"github.com/afrizalsebastian/football-team-management/domain/dao"
	playerdb "github.com/afrizalsebastian/football-team-management/domain/repository/queries/db"
	"github.com/afrizalsebastian/football-team-management/module/logger"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
)

type IPlayersRepository interface {
	CreateTeamPlayer(ctx context.Context, player *dao.Players) error
}

type playerRepository struct {
	db *playerdb.Queries
}

func NewPlayerRepository(pool *pgxpool.Pool) IPlayersRepository {
	return &playerRepository{
		db: playerdb.New(pool),
	}
}

func (d *playerRepository) CreateTeamPlayer(ctx context.Context, player *dao.Players) error {
	l := logger.LoggerNew()

	l.WithContext(ctx).Debug("[CreateTeamPlayer].domain: Started").Msg()

	var (
		teamId   pgtype.UUID
		heightCm pgtype.Numeric
		weightKg pgtype.Numeric
	)

	if err := teamId.Scan(player.TeamId); err != nil {
		return constants.InvalidUUIDValue
	}

	if err := heightCm.Scan(fmt.Sprintf("%f", helper.GetFloat64PtrValue(player.Height))); err != nil {
		return constants.InvalidFieldValue
	}

	if err := weightKg.Scan(fmt.Sprintf("%f", helper.GetFloat64PtrValue(player.Weight))); err != nil {
		return constants.InvalidFieldValue
	}

	params := playerdb.CreatePlayerTeamParams{
		TeamID:       teamId,
		Name:         helper.GetStringPtrValue(player.Name),
		HeightCm:     heightCm,
		WeightKg:     weightKg,
		Position:     playerdb.PlayerPosition(helper.GetStringPtrValue(player.Position)),
		JerseyNumber: int16(helper.GetIntPtrValue(player.JerseyNumber)),
	}

	result, err := d.db.CreatePlayerTeam(ctx, &params)
	if err != nil {
		l.WithContext(ctx).
			Error("Error when create team player").
			Attr("error", err).
			Attr("team_id", player.TeamId).
			Attr("team_name", player.Name).
			Msg()

		if isDuplicateError(err) {
			return constants.DuplicateRow
		}
		return err
	}

	player.Id = result.String()
	l.WithContext(ctx).Debug("[CreateTeamPlayer].domain: Completed").Msg()
	return nil
}
