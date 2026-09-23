package services

import (
	"context"
	"errors"

	"github.com/afrizalsebastian/football-team-management/api"
	"github.com/afrizalsebastian/football-team-management/application/dto"
	"github.com/afrizalsebastian/football-team-management/application/helper"
	"github.com/afrizalsebastian/football-team-management/constants"
	"github.com/afrizalsebastian/football-team-management/domain/dao"
	"github.com/afrizalsebastian/football-team-management/domain/repository"
	"github.com/afrizalsebastian/football-team-management/module/logger"
)

type ITeamsService interface {
	CreateTeam(ctx context.Context, request *dto.CreateTeamRequest) api.WebResponse[*dto.CreateTeamResponse]
	GetListTeam(ctx context.Context) api.WebResponse[[]dto.GetListTeamItem]
	SoftDeleteTeam(ctx context.Context, id string) api.WebResponse[any]
	CreatePlayerTeam(ctx context.Context, teamId string, request *dto.CreatePlayerTeamRequest) api.WebResponse[*dto.CreatePlayerTeamResponse]
}

type teamsService struct {
	teamsRepository  repository.ITeamsRepository
	playerRepository repository.IPlayersRepository
}

func NewTeamsService(
	teamsRepository repository.ITeamsRepository,
	playerRepository repository.IPlayersRepository,
) ITeamsService {
	return &teamsService{
		teamsRepository:  teamsRepository,
		playerRepository: playerRepository,
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

	l.WithContext(ctx).Debug("[SoftDeleteTeam].service: Completed").Msg()
	return api.SuccessResponse[any](
		ctx,
		constants.AcceptDefault.GetMessage(),
		constants.AcceptDefault.GetCode(),
		constants.AcceptDefault.GetHttpCode(),
		nil,
	)
}

func (s *teamsService) CreatePlayerTeam(ctx context.Context, teamId string, request *dto.CreatePlayerTeamRequest) api.WebResponse[*dto.CreatePlayerTeamResponse] {
	l := logger.LoggerNew()

	l.WithContext(ctx).Debug("[CreatePlayerTeam].service: Started").Msg()
	player := &dao.Players{
		TeamId:       teamId,
		Name:         helper.StringPtr(request.Name),
		Height:       helper.Float64Ptr(request.Height),
		Weight:       helper.Float64Ptr(request.Weight),
		Position:     helper.StringPtr(request.Position),
		JerseyNumber: helper.IntPtr(request.JerseyNumber),
	}

	if err := s.playerRepository.CreateTeamPlayer(ctx, player); err != nil {
		l.WithContext(ctx).Debug("error when create team player").Attr("error", err).Msg()
		if errors.Is(err, constants.InvalidUUIDValue) || errors.Is(err, constants.InvalidFieldValue) {
			return api.ErrorResponse[*dto.CreatePlayerTeamResponse](
				ctx,
				constants.BadRequestDefault.GetMessage(),
				constants.BadRequestDefault.GetCode(),
				constants.BadRequestDefault.GetHttpCode(),
				nil,
			)
		}

		if errors.Is(err, constants.DuplicateRow) {
			return api.ErrorResponse[*dto.CreatePlayerTeamResponse](
				ctx,
				constants.InvalidJerseyNumber.GetMessage(),
				constants.InvalidJerseyNumber.GetCode(),
				constants.InvalidJerseyNumber.GetHttpCode(),
				nil,
			)
		}

		return api.ErrorResponse[*dto.CreatePlayerTeamResponse](
			ctx,
			constants.InternalServerError.GetMessage(),
			constants.InternalServerError.GetCode(),
			constants.InternalServerError.GetHttpCode(),
			nil,
		)
	}

	response := &dto.CreatePlayerTeamResponse{
		Id:           player.Id,
		TeamId:       teamId,
		Height:       helper.GetFloat64PtrValue(player.Height),
		Weight:       helper.GetFloat64PtrValue(player.Weight),
		Position:     helper.GetStringPtrValue(player.Position),
		JerseyNumber: helper.GetIntPtrValue(player.JerseyNumber),
	}

	return api.SuccessResponse[*dto.CreatePlayerTeamResponse](
		ctx,
		constants.CreatedDefault.GetMessage(),
		constants.CreatedDefault.GetCode(),
		constants.CreatedDefault.GetHttpCode(),
		response,
	)
}
