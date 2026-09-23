package repository

import (
	"context"
	"errors"

	"github.com/afrizalsebastian/football-team-management/application/helper"
	"github.com/afrizalsebastian/football-team-management/constants"
	"github.com/afrizalsebastian/football-team-management/domain/dao"
	teamsdb "github.com/afrizalsebastian/football-team-management/domain/repository/queries/db"
	"github.com/afrizalsebastian/football-team-management/module/logger"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
)

type ITeamsRepository interface {
	CreateTeam(ctx context.Context, team *dao.Teams) error
	GetListTeam(ctx context.Context) ([]dao.Teams, error)
	SoftDeleteTeam(ctx context.Context, id string) error
	GetTeamDetail(ctx context.Context, teamId string) (*dao.Teams, error)
	UpdatePartialTeam(ctx context.Context, team *dao.Teams) error
	IsTeamExists(ctx context.Context, teamId string) bool
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
		Name:        *team.Name,
		Logo:        StringPtrToPgtypeText(team.Logo),
		FoundedYear: StringPtrToPgtypeText(team.FoundedYear),
		Address:     *team.Address,
		City:        *team.City,
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

	if _, err := d.db.SoftDeleteTeams(ctx, uuidTeams); err != nil {
		l.WithContext(ctx).Error("error when soft delete team").Attr("team_id", id).Msg()
		if errors.Is(err, pgx.ErrNoRows) {
			return constants.ErrNotFoundRow
		}
		return err
	}

	l.WithContext(ctx).Debug("[SoftDeleteTeam].domain: Completed").Attr("team_id", id).Msg()
	return nil
}

func (d *teamsRepository) GetTeamDetail(ctx context.Context, teamId string) (*dao.Teams, error) {
	l := logger.LoggerNew()

	l.WithContext(ctx).Debug("[GetTeamDetail].domain: Started").Attr("team_id", teamId).Msg()

	var uuidTeams pgtype.UUID
	if err := uuidTeams.Scan(teamId); err != nil {
		return nil, constants.InvalidUUIDValue
	}

	row, err := d.db.GetTeamDetail(ctx, uuidTeams)
	if err != nil {
		l.WithContext(ctx).Error("error when get team detail").Attr("team_id", teamId).Msg()
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, constants.ErrNotFoundRow
		}
		return nil, err
	}

	result := &dao.Teams{
		Id:          row.ID.String(),
		Name:        helper.StringPtr(row.Name),
		Logo:        helper.StringPtr(row.Logo.String),
		FoundedYear: helper.StringPtr(row.FoundedYear.String),
		Address:     helper.StringPtr(row.Address),
		City:        helper.StringPtr(row.City),
		CreatedAt:   row.CreatedAt.Time,
		UpdatedAt:   row.UpdatedAt.Time,
		IsDeleted:   row.IsDeleted.Bool,
		DeletedAt:   row.CreatedAt.Time,
	}

	return result, nil
}

func (d *teamsRepository) UpdatePartialTeam(ctx context.Context, team *dao.Teams) error {
	l := logger.LoggerNew()

	l.WithContext(ctx).Debug("[GetTeamDetail].domain: Started").Attr("team_id", team.Id).Msg()
	var uuidTeams pgtype.UUID
	if err := uuidTeams.Scan(team.Id); err != nil {
		return constants.InvalidUUIDValue
	}

	params := &teamsdb.UpdateTeamsParams{
		Name:        StringPtrToPgtypeText(team.Name),
		Logo:        StringPtrToPgtypeText(team.Logo),
		FoundedYear: StringPtrToPgtypeText(team.FoundedYear),
		Address:     StringPtrToPgtypeText(team.Address),
		City:        StringPtrToPgtypeText(team.City),
		ID:          uuidTeams,
	}

	result, err := d.db.UpdateTeams(ctx, params)
	if err != nil {
		l.WithContext(ctx).Error("error when update team detail").Attr("team_id", team.Id).Msg()
		if errors.Is(err, pgx.ErrNoRows) {
			return constants.ErrNotFoundRow
		}
		return err
	}

	team.Id = result.ID.String()
	team.Name = helper.StringPtr(result.Name)
	team.Logo = helper.StringPtr(result.Logo.String)
	team.FoundedYear = helper.StringPtr(result.FoundedYear.String)
	team.Address = helper.StringPtr(result.Address)
	team.City = helper.StringPtr(result.City)
	team.CreatedAt = result.CreatedAt.Time
	team.UpdatedAt = result.UpdatedAt.Time
	team.IsDeleted = result.IsDeleted.Bool
	team.DeletedAt = result.DeletedAt.Time

	return nil
}

func (d *teamsRepository) IsTeamExists(ctx context.Context, teamId string) bool {
	l := logger.LoggerNew()

	l.WithContext(ctx).Debug("[GetTeamDetail].domain: Started").Attr("team_id", teamId).Msg()

	var uuidTeams pgtype.UUID
	if err := uuidTeams.Scan(teamId); err != nil {
		return false
	}
	_, err := d.db.CheckTeamExisits(ctx, uuidTeams)
	if err != nil {
		l.WithContext(ctx).Error("error when check is team exists").Attr("team_id", teamId).Attr("error", err).Msg()
		return false
	}

	return true
}
