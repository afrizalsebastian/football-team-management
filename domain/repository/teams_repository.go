package repository

import (
	"context"

	"github.com/afrizalsebastian/football-team-management/application/helper"
	"github.com/afrizalsebastian/football-team-management/constants"
	"github.com/afrizalsebastian/football-team-management/domain/dao"
	teamsdb "github.com/afrizalsebastian/football-team-management/domain/repository/queries/db"
	"github.com/afrizalsebastian/football-team-management/module/logger"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
)

type ITeamsRepository interface {
	CreateTeam(ctx context.Context, team *dao.Teams) error
	GetListTeam(ctx context.Context) ([]dao.Teams, error)
	SoftDeleteTeam(ctx context.Context, id string) error
}

type teamsRepository struct {
	db *teamsdb.Queries
}

func NewTeamsRepository(dbPool *pgxpool.Pool) ITeamsRepository {
	return &teamsRepository{
		db: teamsdb.New(dbPool),
	}
}

func (d *teamsRepository) CreateTeam(ctx context.Context, team *dao.Teams) error {
	l := logger.LoggerNew()

	l.WithContext(ctx).Debug("[CreateTeam].domain: Started").Msg()
	params := teamsdb.CreateTeamParams{
		Name: *team.Name,
		Logo: pgtype.Text{
			String: *team.Logo,
			Valid:  true,
		},
		FoundedYear: pgtype.Text{
			String: *team.FoundedYear,
			Valid:  true,
		},
		Address: *team.Address,
		City:    *team.City,
	}

	result, err := d.db.CreateTeam(ctx, &params)
	if err != nil {
		l.WithContext(ctx).
			Error("Error when create team").
			Attr("error", err).
			Attr("team_name", team.Name).
			Msg()

		return err
	}

	team.Id = result.String()
	l.WithContext(ctx).Debug("[CreateTeam].domain: Completed").Msg()
	return nil
}

func (d *teamsRepository) GetListTeam(ctx context.Context) ([]dao.Teams, error) {
	l := logger.LoggerNew()

	l.WithContext(ctx).Debug("[GetListTeam].domain: Started").Msg()
	rows, err := d.db.GetListTeams(ctx)
	if err != nil {
		l.WithContext(ctx).Error("Error when get list teams").Attr("error", err).Msg()
		return nil, err
	}

	result := make([]dao.Teams, 0)
	for _, r := range rows {
		if r == nil {
			continue
		}

		result = append(result, dao.Teams{
			Id:          r.ID.String(),
			Name:        helper.StringPtr(r.Name),
			Logo:        helper.StringPtr(r.Logo.String),
			FoundedYear: helper.StringPtr(r.FoundedYear.String),
		})
	}

	l.WithContext(ctx).Debug("[GetListTeam].domain: Completed").Msg()
	return result, nil
}

func (d *teamsRepository) SoftDeleteTeam(ctx context.Context, id string) error {
	l := logger.LoggerNew()

	l.WithContext(ctx).Debug("[SoftDeleteTeam].domain: Started").Attr("team_id", id).Msg()

	var uuidTeams pgtype.UUID
	if err := uuidTeams.Scan(id); err != nil {
		return constants.InvalidUUIDValue
	}

	if err := d.db.SoftDeleteTeams(ctx, uuidTeams); err != nil {
		l.WithContext(ctx).Debug("error when soft delete team").Attr("team_id", id).Msg()
		return err
	}

	l.WithContext(ctx).Debug("[SoftDeleteTeam].domain: Completed").Attr("team_id", id).Msg()
	return nil
}
