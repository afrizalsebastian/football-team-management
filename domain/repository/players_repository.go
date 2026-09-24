package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/afrizalsebastian/football-team-management/application/helper"
	"github.com/afrizalsebastian/football-team-management/constants"
	"github.com/afrizalsebastian/football-team-management/domain/dao"
	playerdb "github.com/afrizalsebastian/football-team-management/domain/repository/queries/db"
	"github.com/afrizalsebastian/football-team-management/module/logger"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
)

type IPlayersRepository interface {
	CreateTeamPlayer(ctx context.Context, player *dao.Players) error
	GetPlayerTeamByTeamId(ctx context.Context, teamId string) ([]dao.Players, error)
	GetListPlayer(ctx context.Context) ([]dao.Players, error)
	UpdatePartialPlayer(ctx context.Context, player *dao.Players) error
	GetPlayerDetail(ctx context.Context, playerId string) (*dao.Players, error)
	DeletePlayer(ctx context.Context, playerId string) error
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

func (d *playerRepository) GetPlayerTeamByTeamId(ctx context.Context, teamIdStr string) ([]dao.Players, error) {
	l := logger.LoggerNew()

	l.WithContext(ctx).Debug("[GetPlayerTeamByTeamId].domain: Started").Msg()

	var teamId pgtype.UUID
	if err := teamId.Scan(teamIdStr); err != nil {
		return nil, constants.InvalidUUIDValue
	}

	rows, err := d.db.GetListPlayerTeam(ctx, teamId)
	if err != nil {
		l.WithContext(ctx).
			Error("Error when get list team player").
			Attr("error", err).
			Attr("team_id", teamIdStr).
			Msg()

		return nil, err
	}

	result := make([]dao.Players, 0)
	for _, r := range rows {
		if r == nil {
			continue
		}

		result = append(result, dao.Players{
			Id:           r.ID.String(),
			TeamId:       r.TeamID.String(),
			Name:         helper.StringPtr(r.Name),
			Position:     helper.StringPtr(string(r.Position)),
			JerseyNumber: helper.IntPtr(int(r.JerseyNumber)),
		})
	}

	return result, nil
}

func (d *playerRepository) GetListPlayer(ctx context.Context) ([]dao.Players, error) {
	l := logger.LoggerNew()

	l.WithContext(ctx).Debug("[GetListPlayer].domain: Started").Msg()

	rows, err := d.db.GetListPlayer(ctx)
	if err != nil {
		l.WithContext(ctx).
			Error("Error when get list player").
			Attr("error", err).
			Msg()

		return nil, err
	}

	result := make([]dao.Players, 0)
	for _, r := range rows {
		if r == nil {
			continue
		}

		result = append(result, dao.Players{
			Id:           r.ID.String(),
			Name:         helper.StringPtr(r.Name),
			Position:     helper.StringPtr(string(r.Position)),
			JerseyNumber: helper.IntPtr(int(r.JerseyNumber)),
			TeamId:       r.TeamID.String(),
			Team: &dao.Teams{
				Id:   r.TeamID.String(),
				Name: helper.StringPtr(r.TeamName.String),
			},
		})
	}

	l.WithContext(ctx).Debug("[GetListPlayer].domain: Completed").Msg()
	return result, nil
}

func (d *playerRepository) UpdatePartialPlayer(ctx context.Context, player *dao.Players) error {
	l := logger.LoggerNew()

	l.WithContext(ctx).Debug("[UpdatePartialPlayer].domain: Started").Msg()

	params, err := d.createUpdatePlayerParams(player)
	if err != nil {
		l.WithContext(ctx).Error("error when mapping params").
			Attr("error", err).Attr("player_id", player.Id).Msg()
		return err
	}

	row, err := d.db.UpdatePlayers(ctx, params)
	if err != nil {
		l.WithContext(ctx).Error("error when update player data").
			Attr("error", err).Attr("player_id", player.Id).Msg()
		if errors.Is(err, pgx.ErrNoRows) {
			return constants.ErrNotFoundRow
		}
		return err
	}

	height, _ := row.HeightCm.Float64Value()
	weight, _ := row.WeightKg.Float64Value()

	player.Name = helper.StringPtr(row.Name)
	player.Height = helper.Float64Ptr(height.Float64)
	player.Weight = helper.Float64Ptr(weight.Float64)
	player.TeamId = row.TeamID.String()
	player.Position = helper.StringPtr(string(row.Position))
	player.JerseyNumber = helper.IntPtr(int(row.JerseyNumber))
	player.CreatedAt = row.CreatedAt.Time
	player.UpdatedAt = row.UpdatedAt.Time
	player.IsDeleted = row.IsDeleted.Bool
	player.DeletedAt = row.DeletedAt.Time

	return nil
}

func (d *playerRepository) createUpdatePlayerParams(player *dao.Players) (*playerdb.UpdatePlayersParams, error) {
	var (
		teamId   pgtype.UUID
		id       pgtype.UUID
		heightCm pgtype.Numeric
		weightKg pgtype.Numeric
	)

	if err := id.Scan(player.Id); err != nil {
		return nil, constants.InvalidUUIDValue
	}

	if player.TeamId != "" {
		if err := teamId.Scan(player.Team.Id); err != nil {
			return nil, constants.InvalidUUIDValue
		}
	}

	if player.Height != nil {
		if err := heightCm.Scan(fmt.Sprintf("%f", helper.GetFloat64PtrValue(player.Height))); err != nil {
			return nil, constants.InvalidFieldValue
		}
	}

	if player.Weight != nil {
		if err := weightKg.Scan(fmt.Sprintf("%f", helper.GetFloat64PtrValue(player.Weight))); err != nil {
			return nil, constants.InvalidFieldValue
		}
	}

	return &playerdb.UpdatePlayersParams{
		ID:           id,
		TeamID:       teamId,
		HeightCm:     heightCm,
		WeightKg:     weightKg,
		Name:         StringPtrToPgtypeText(player.Name),
		JerseyNumber: IntPtrToPgTypeInt2(player.JerseyNumber),
		Position:     PositionPtrToNullPlayerPosition(player.Position),
	}, nil
}

func (d *playerRepository) GetPlayerDetail(ctx context.Context, playerIdStr string) (*dao.Players, error) {
	l := logger.LoggerNew()

	l.WithContext(ctx).Debug("[GetPlayerDetail].domain: Started").Msg()

	var playerId pgtype.UUID
	if err := playerId.Scan(playerIdStr); err != nil {
		return nil, constants.InvalidUUIDValue
	}

	row, err := d.db.GetPlayerDetail(ctx, playerId)
	if err != nil {
		l.WithContext(ctx).Error("error when get player detail").
			Attr("error", err).Attr("player_id", playerIdStr).Msg()
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, constants.ErrNotFoundRow
		}
		return nil, err
	}

	height, _ := row.HeightCm.Float64Value()
	weight, _ := row.WeightKg.Float64Value()
	result := &dao.Players{
		Id:           row.ID.String(),
		Name:         helper.StringPtr(row.Name),
		Height:       helper.Float64Ptr(height.Float64),
		Weight:       helper.Float64Ptr(weight.Float64),
		TeamId:       row.TeamID.String(),
		Position:     helper.StringPtr(string(row.Position)),
		JerseyNumber: helper.IntPtr(int(row.JerseyNumber)),
		CreatedAt:    row.CreatedAt.Time,
		UpdatedAt:    row.UpdatedAt.Time,
		IsDeleted:    row.IsDeleted.Bool,
		DeletedAt:    row.DeletedAt.Time,
		GoalsCount:   int(row.GoalsCount),
		Team: &dao.Teams{
			Id:   row.TeamID.String(),
			Name: helper.StringPtr(row.TeamName.String),
		},
	}

	return result, nil
}

func (d *playerRepository) DeletePlayer(ctx context.Context, playerIdStr string) error {
	l := logger.LoggerNew()

	l.WithContext(ctx).Debug("[DeletePlayer].domain: Started").Msg()

	var playerId pgtype.UUID
	if err := playerId.Scan(playerIdStr); err != nil {
		return constants.InvalidUUIDValue
	}

	_, err := d.db.DeletePlayer(ctx, playerId)
	if err != nil {
		l.WithContext(ctx).Error("error when delete player").
			Attr("error", err).Attr("player_id", playerIdStr).Msg()
		if errors.Is(err, pgx.ErrNoRows) {
			return constants.ErrNotFoundRow
		}
		return err
	}

	l.WithContext(ctx).Debug("[DeletePlayer].domain: Started").Msg()
	return nil
}
