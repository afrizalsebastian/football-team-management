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
	GetTeamDetail(ctx context.Context, teamId string) api.WebResponse[*dto.GetTeamDetail]
	SoftDeleteTeam(ctx context.Context, id string) api.WebResponse[any]
	CreatePlayerTeam(ctx context.Context, teamId string, request *dto.CreatePlayerTeamRequest) api.WebResponse[*dto.CreatePlayerTeamResponse]
	GetListPlayerTeam(ctx context.Context, teamId string) api.WebResponse[[]dto.GetListTeamPlayerItem]
	UpdateTeamData(ctx context.Context, teamId string, request *dto.UpdateTeamRequest) api.WebResponse[*dto.GetTeamDetail]
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
		l.WithContext(ctx).Warn("empty teams data").Msg()
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
		l.WithContext(ctx).Error("error when soft delete team").Attr("error", err).Msg()

		errMsg := constants.InternalServerError
		if errors.Is(err, constants.InvalidUUIDValue) {
			errMsg = constants.BadRequestDefault
		}
		if errors.Is(err, constants.ErrNotFoundRow) {
			errMsg = constants.NotFoundDefault
		}
		return api.ErrorResponse[any](
			ctx,
			errMsg.GetMessage(),
			errMsg.GetCode(),
			errMsg.GetHttpCode(),
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
	if found := s.teamsRepository.IsTeamExists(ctx, teamId); !found {
		l.WithContext(ctx).Warn("not found team").Msg()
		return api.ErrorResponse[*dto.CreatePlayerTeamResponse](
			ctx,
			constants.NotFoundTeam.GetMessage(),
			constants.NotFoundTeam.GetCode(),
			constants.NotFoundTeam.GetHttpCode(),
			nil,
		)
	}

	player := &dao.Players{
		TeamId:       teamId,
		Name:         helper.StringPtr(request.Name),
		Height:       helper.Float64Ptr(request.Height),
		Weight:       helper.Float64Ptr(request.Weight),
		Position:     helper.StringPtr(request.Position),
		JerseyNumber: helper.IntPtr(request.JerseyNumber),
	}

	if err := s.playerRepository.CreateTeamPlayer(ctx, player); err != nil {
		l.WithContext(ctx).Error("error when create team player").Attr("error", err).Msg()

		errMsg := constants.InternalServerError
		if errors.Is(err, constants.InvalidUUIDValue) || errors.Is(err, constants.InvalidFieldValue) {
			errMsg = constants.BadRequestDefault
		}

		if errors.Is(err, constants.DuplicateRow) {
			errMsg = constants.InvalidJerseyNumber
		}

		return api.ErrorResponse[*dto.CreatePlayerTeamResponse](
			ctx,
			errMsg.GetMessage(),
			errMsg.GetCode(),
			errMsg.GetHttpCode(),
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

func (s *teamsService) GetListPlayerTeam(ctx context.Context, teamId string) api.WebResponse[[]dto.GetListTeamPlayerItem] {
	l := logger.LoggerNew()

	l.WithContext(ctx).Debug("[GetListPlayerTeam].service: Started").Msg()
	if found := s.teamsRepository.IsTeamExists(ctx, teamId); !found {
		l.WithContext(ctx).Warn("not found team").Msg()
		return api.ErrorResponse[[]dto.GetListTeamPlayerItem](
			ctx,
			constants.NotFoundTeam.GetMessage(),
			constants.NotFoundTeam.GetCode(),
			constants.NotFoundTeam.GetHttpCode(),
			nil,
		)
	}

	result, err := s.playerRepository.GetPlayerTeamByTeamId(ctx, teamId)
	if err != nil {
		l.WithContext(ctx).Error("error when get team player").Attr("error", err).Attr("team_id", teamId).Msg()
		errMsg := constants.InternalServerError
		if errors.Is(err, constants.InvalidUUIDValue) {
			errMsg = constants.BadRequestDefault
		}
		return api.ErrorResponse[[]dto.GetListTeamPlayerItem](
			ctx,
			errMsg.GetMessage(),
			errMsg.GetCode(),
			errMsg.GetHttpCode(),
			nil,
		)
	}

	if len(result) == 0 {
		l.WithContext(ctx).Debug("empty teams data").Msg()
		return api.SuccessResponse[[]dto.GetListTeamPlayerItem](
			ctx,
			constants.SuccesssWithEmptyList.GetMessage(),
			constants.SuccesssWithEmptyList.GetCode(),
			constants.SuccesssWithEmptyList.GetHttpCode(),
			nil,
		)
	}

	response := make([]dto.GetListTeamPlayerItem, 0)
	for _, r := range result {
		var positionCode, positionTitle string
		position := constants.DictPlayerPosition.GetValue(helper.GetStringPtrValue(r.Position))
		if position != nil {
			positionCode, positionTitle = position.Code, position.Title
		}

		response = append(response, dto.GetListTeamPlayerItem{
			Id:     r.Id,
			TeamId: r.TeamId,
			Name:   helper.GetStringPtrValue(r.Name),
			Position: dto.PlayerPosition{
				Code:  positionCode,
				Title: positionTitle,
			},
			JerseyNumber: helper.GetIntPtrValue(r.JerseyNumber),
		})
	}

	l.WithContext(ctx).Debug("[GetListPlayerTeam].service: Completed").Msg()
	return api.SuccessResponse(
		ctx,
		constants.SuccessDefault.GetMessage(),
		constants.SuccessDefault.GetCode(),
		constants.SuccessDefault.GetHttpCode(),
		response,
	)
}

func (s *teamsService) GetTeamDetail(ctx context.Context, teamId string) api.WebResponse[*dto.GetTeamDetail] {
	l := logger.LoggerNew()

	l.WithContext(ctx).Debug("[GetTeamDetail].service: Started").Msg()
	result, err := s.teamsRepository.GetTeamDetail(ctx, teamId)
	if err != nil {
		l.WithContext(ctx).Error("error when get team detail").Attr("error", err).Attr("team_id", teamId).Msg()

		errMsg := constants.InternalServerError
		if errors.Is(err, constants.InvalidUUIDValue) {
			errMsg = constants.BadRequestDefault
		}
		if errors.Is(err, constants.ErrNotFoundRow) {
			errMsg = constants.NotFoundDefault
		}
		return api.ErrorResponse[*dto.GetTeamDetail](
			ctx,
			errMsg.GetMessage(),
			errMsg.GetCode(),
			errMsg.GetHttpCode(),
			nil,
		)
	}

	response := &dto.GetTeamDetail{
		Id:          result.Id,
		Name:        helper.GetStringPtrValue(result.Name),
		Logo:        helper.GetStringPtrValue(result.Logo),
		Address:     helper.GetStringPtrValue(result.Address),
		City:        helper.GetStringPtrValue(result.City),
		FoundedYear: helper.GetStringPtrValue(result.FoundedYear),
		CreatedAt:   result.CreatedAt,
		UpdatedAt:   result.UpdatedAt,
		IsDeleted:   result.IsDeleted,
		DeletedAt:   result.DeletedAt,
	}

	l.WithContext(ctx).Debug("[GetTeamDetail].service: Completed").Msg()
	return api.SuccessResponse(
		ctx,
		constants.SuccessDefault.GetMessage(),
		constants.SuccessDefault.GetCode(),
		constants.SuccessDefault.GetHttpCode(),
		response,
	)
}

func (s *teamsService) UpdateTeamData(ctx context.Context, teamId string, request *dto.UpdateTeamRequest) api.WebResponse[*dto.GetTeamDetail] {
	l := logger.LoggerNew()

	l.WithContext(ctx).Debug("[UpdateTeamData].service: Started").Msg()
	team := &dao.Teams{
		Id:          teamId,
		Name:        request.Name,
		Logo:        request.Logo,
		FoundedYear: request.FoundedYear,
		Address:     request.Address,
		City:        request.City,
	}

	if err := s.teamsRepository.UpdatePartialTeam(ctx, team); err != nil {
		l.WithContext(ctx).Error("error when update team data").
			Attr("error", err).Attr("team_id", teamId).Msg()

		errMsg := constants.InternalServerError
		if errors.Is(err, constants.InvalidUUIDValue) {
			errMsg = constants.BadRequestDefault
		}
		if errors.Is(err, constants.ErrNotFoundRow) {
			errMsg = constants.NotFoundDefault
		}
		return api.ErrorResponse[*dto.GetTeamDetail](
			ctx,
			errMsg.GetMessage(),
			errMsg.GetCode(),
			errMsg.GetHttpCode(),
			nil,
		)
	}

	response := &dto.GetTeamDetail{
		Id:          team.Id,
		Name:        helper.GetStringPtrValue(team.Name),
		Logo:        helper.GetStringPtrValue(team.Logo),
		Address:     helper.GetStringPtrValue(team.Address),
		City:        helper.GetStringPtrValue(team.City),
		FoundedYear: helper.GetStringPtrValue(team.FoundedYear),
		CreatedAt:   team.CreatedAt,
		UpdatedAt:   team.UpdatedAt,
		IsDeleted:   team.IsDeleted,
		DeletedAt:   team.DeletedAt,
	}

	l.WithContext(ctx).Debug("[UpdateTeamData].service: Completed").Msg()
	return api.SuccessResponse(
		ctx,
		constants.SuccessDefault.GetMessage(),
		constants.SuccessDefault.GetCode(),
		constants.SuccessDefault.GetHttpCode(),
		response,
	)
}
