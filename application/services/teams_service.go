package services

import (
	"context"
	"errors"

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
	GetListTeam(ctx context.Context) api.WebResponse[[]dto.GetListTeamItem]
	SoftDeleteTeam(ctx context.Context, id string) api.WebResponse[any]
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
		l.WithContext(ctx).Error("error when create team").Attr("error", err).Msg()
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

	l.WithContext(ctx).Debug("[CreateTeam].service: Completed").Msg()
	return api.SuccessResponse(
		ctx,
		constants.CreatedDefault.GetMessage(),
		constants.CreatedDefault.GetCode(),
		constants.CreatedDefault.GetHttpCode(),
		response,
	)
}

func (s *teamsService) GetListTeam(ctx context.Context) api.WebResponse[[]dto.GetListTeamItem] {
	l := logger.LoggerNew()

	l.WithContext(ctx).Debug("[GetListTeam].service: Started").Msg()

	teams, err := s.teamsRepository.GetListTeam(ctx)
	if err != nil {
		l.WithContext(ctx).Error("error when create team").Attr("error", err).Msg()
		return api.ErrorResponse[[]dto.GetListTeamItem](
			ctx,
			constants.InternalServerError.GetMessage(),
			constants.InternalServerError.GetCode(),
			constants.InternalServerError.GetHttpCode(),
			nil,
		)
	}

	if len(teams) == 0 {
		l.WithContext(ctx).Debug("empty teams data").Msg()
		return api.SuccessResponse[[]dto.GetListTeamItem](
			ctx,
			constants.SuccesssWithEmptyList.GetMessage(),
			constants.SuccesssWithEmptyList.GetCode(),
			constants.SuccesssWithEmptyList.GetHttpCode(),
			nil,
		)
	}

	response := make([]dto.GetListTeamItem, 0)
	for _, t := range teams {
		response = append(response, dto.GetListTeamItem{
			Id:          t.Id,
			Name:        helper.GetStringPtrValue(t.Name),
			Logo:        helper.GetStringPtrValue(t.Logo),
			FoundedYear: helper.GetStringPtrValue(t.FoundedYear),
		})
	}

	l.WithContext(ctx).Debug("[CreateTeam].service: Completed").Msg()
	return api.SuccessResponse(
		ctx,
		constants.SuccessDefault.GetMessage(),
		constants.SuccessDefault.GetCode(),
		constants.SuccessDefault.GetHttpCode(),
		response,
	)
}

func (s *teamsService) SoftDeleteTeam(ctx context.Context, id string) api.WebResponse[any] {
	l := logger.LoggerNew()

	l.WithContext(ctx).Debug("[SoftDeleteTeam].service: Started").Msg()
	if err := s.teamsRepository.SoftDeleteTeam(ctx, id); err != nil {
		l.WithContext(ctx).Debug("error when soft delete team").Attr("error", err).Msg()

		if errors.Is(err, constants.InvalidUUIDValue) {
			return api.ErrorResponse[any](
				ctx,
				constants.BadRequestDefault.GetMessage(),
				constants.BadRequestDefault.GetCode(),
				constants.BadRequestDefault.GetHttpCode(),
				nil,
			)
		}

		return api.ErrorResponse[any](
			ctx,
			constants.InternalServerError.GetMessage(),
			constants.InternalServerError.GetCode(),
			constants.InternalServerError.GetHttpCode(),
			nil,
		)
	}

	return api.SuccessResponse[any](
		ctx,
		constants.AcceptDefault.GetMessage(),
		constants.AcceptDefault.GetCode(),
		constants.AcceptDefault.GetHttpCode(),
		nil,
	)
}
