package services

import (
	"context"

	"github.com/afrizalsebastian/football-team-management/api"
	"github.com/afrizalsebastian/football-team-management/application/dto"
	"github.com/afrizalsebastian/football-team-management/application/helper"
	"github.com/afrizalsebastian/football-team-management/constants"
	"github.com/afrizalsebastian/football-team-management/domain/dao"
	teams_repository "github.com/afrizalsebastian/football-team-management/domain/repository/teams-repository"
	"github.com/afrizalsebastian/football-team-management/module/logger"
)

type ITeamsService interface {
	CreateTeam(ctx context.Context, request *dto.CreateTeamRequest) api.WebResponse[*dto.CreateTeamResponse]
}

type teamsService struct {
	teamsRepository teams_repository.ITeamsRepository
}

func NewTeamsService(teamsRepository teams_repository.ITeamsRepository) ITeamsService {
	return &teamsService{
		teamsRepository: teamsRepository,
	}
}

func (s *teamsService) CreateTeam(ctx context.Context, request *dto.CreateTeamRequest) api.WebResponse[*dto.CreateTeamResponse] {
	l := logger.LoggerNew()

	l.WithContext(ctx).Debug("[CreateTeam].service: Started").Msg()
	team := &dao.Teams{
		Name:        helper.StringPtr(request.Name),
		Logo:        helper.StringPtr(request.Logo),
		FoundedYear: helper.StringPtr(request.FoundedYear),
		Address:     helper.StringPtr(request.Address),
		City:        helper.StringPtr(request.City),
	}

	if err := s.teamsRepository.CreateTeam(ctx, team); err != nil {
		l.WithContext(ctx).Error("[CreateTeam].service: Started").Attr("error", err).Msg()
		return api.ErrorResponse[*dto.CreateTeamResponse](
			ctx,
			constants.InternalServerError.GetMessage(),
			constants.InternalServerError.GetCode(),
			constants.InternalServerError.GetHttpCode(),
			nil,
		)
	}

	response := &dto.CreateTeamResponse{
		Id:          team.Id,
		Name:        helper.GetStringPtrValue(team.Name),
		Logo:        helper.GetStringPtrValue(team.Logo),
		FoundedYear: helper.GetStringPtrValue(team.FoundedYear),
		Address:     helper.GetStringPtrValue(team.Address),
		City:        helper.GetStringPtrValue(team.City),
	}

	return api.SuccessResponse(
		ctx,
		constants.CreatedDefault.GetMessage(),
		constants.CreatedDefault.GetCode(),
		constants.CreatedDefault.GetHttpCode(),
		response,
	)
}
