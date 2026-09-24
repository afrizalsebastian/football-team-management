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

type IPlayersService interface {
	GetListPlayer(ctx context.Context) api.WebResponse[[]dto.GetListPlayerItem]
	UpdatePlayer(ctx context.Context, playerId string, request *dto.UpdatePlayerRequest) api.WebResponse[*dto.GetPlayerDetailResponse]
	GetPlayerDetail(ctx context.Context, playerId string) api.WebResponse[*dto.GetPlayerDetailResponse]
	DeletePlayer(ctx context.Context, playerId string) api.WebResponse[any]
}

type playerService struct {
	playerRepository repository.IPlayersRepository
	teamRepository   repository.ITeamsRepository
}

func NewPlayerService(
	playerRepository repository.IPlayersRepository,
	teamRepository repository.ITeamsRepository,
) IPlayersService {
	return &playerService{
		playerRepository: playerRepository,
		teamRepository:   teamRepository,
	}
}

func (s *playerService) GetListPlayer(ctx context.Context) api.WebResponse[[]dto.GetListPlayerItem] {
	l := logger.LoggerNew()

	l.WithContext(ctx).Debug("[GetListPlayer].service: Started").Msg()

	result, err := s.playerRepository.GetListPlayer(ctx)
	if err != nil {
		l.WithContext(ctx).Error("error when get list player").Attr("error", err).Msg()
		return api.ErrorResponse[[]dto.GetListPlayerItem](
			ctx,
			constants.InternalServerError.GetMessage(),
			constants.InternalServerError.GetCode(),
			constants.InternalServerError.GetHttpCode(),
			nil,
		)
	}

	if len(result) == 0 {
		return api.SuccessResponse[[]dto.GetListPlayerItem](
			ctx,
			constants.SuccesssWithEmptyList.GetMessage(),
			constants.SuccesssWithEmptyList.GetCode(),
			constants.SuccesssWithEmptyList.GetHttpCode(),
			nil,
		)
	}

	response := make([]dto.GetListPlayerItem, 0)
	for _, r := range result {
		var positionCode, positionTitle string
		position := constants.DictPlayerPosition.GetValue(helper.GetStringPtrValue(r.Position))
		if position != nil {
			positionCode, positionTitle = position.Code, position.Title
		}
		response = append(response, dto.GetListPlayerItem{
			Id:   r.Id,
			Name: helper.GetStringPtrValue(r.Name),
			Position: &dto.PlayerPosition{
				Code:  positionCode,
				Title: positionTitle,
			},
			JerseyNumber: helper.GetIntPtrValue(r.JerseyNumber),
			Team: dto.ListPlayerItemTeam{
				Id:   r.Team.Id,
				Name: helper.GetStringPtrValue(r.Team.Name),
			},
		})
	}

	return api.SuccessResponse[[]dto.GetListPlayerItem](
		ctx,
		constants.SuccessDefault.GetMessage(),
		constants.SuccessDefault.GetCode(),
		constants.SuccessDefault.GetHttpCode(),
		response,
	)
}

func (s *playerService) UpdatePlayer(ctx context.Context, playerId string, request *dto.UpdatePlayerRequest) api.WebResponse[*dto.GetPlayerDetailResponse] {
	l := logger.LoggerNew()

	l.WithContext(ctx).Debug("[UpdatePlayer].service: Started").Msg()

	if request.TeamId != nil {
		if found := s.teamRepository.IsTeamExists(ctx, *request.TeamId); !found {
			l.WithContext(ctx).Warn("not found team").Msg()
			return api.ErrorResponse[*dto.GetPlayerDetailResponse](
				ctx,
				constants.NotFoundTeam.GetMessage(),
				constants.NotFoundTeam.GetCode(),
				constants.NotFoundTeam.GetHttpCode(),
				nil,
			)
		}
	}

	player := &dao.Players{
		Id:           playerId,
		TeamId:       helper.GetStringPtrValue(request.TeamId),
		Height:       request.Height,
		Weight:       request.Weight,
		Position:     request.Position,
		JerseyNumber: request.JerseyNumber,
	}
	if err := s.playerRepository.UpdatePartialPlayer(ctx, player); err != nil {
		l.WithContext(ctx).Error("error when update player").Attr("error", err).Msg()

		errMsg := constants.InternalServerError
		if errors.Is(err, constants.InvalidUUIDValue) {
			errMsg = constants.BadRequestDefault
		}
		if errors.Is(err, constants.ErrNotFoundRow) {
			errMsg = constants.NotFoundDefault
		}
		return api.ErrorResponse[*dto.GetPlayerDetailResponse](
			ctx,
			errMsg.GetMessage(),
			errMsg.GetCode(),
			errMsg.GetHttpCode(),
			nil,
		)
	}

	var positionCode, positionTitle string
	position := constants.DictPlayerPosition.GetValue(helper.GetStringPtrValue(player.Position))
	if position != nil {
		positionCode, positionTitle = position.Code, position.Title
	}

	response := &dto.GetPlayerDetailResponse{
		Id:     player.Id,
		TeamId: player.TeamId,
		Name:   helper.GetStringPtrValue(player.Name),
		Height: helper.GetFloat64PtrValue(player.Height),
		Weight: helper.GetFloat64PtrValue(player.Weight),
		Position: dto.PlayerPosition{
			Code:  positionCode,
			Title: positionTitle,
		},
		JerseyNumber: helper.GetIntPtrValue(player.JerseyNumber),
		CreatedAt:    player.CreatedAt,
		UpdatedAt:    player.UpdatedAt,
		IsDeleted:    player.IsDeleted,
		DeletedAt:    player.DeletedAt,
	}

	l.WithContext(ctx).Debug("[UpdatePlayer].service: Completed").Msg()
	return api.SuccessResponse(
		ctx,
		constants.SuccessDefault.GetMessage(),
		constants.SuccessDefault.GetCode(),
		constants.SuccessDefault.GetHttpCode(),
		response,
	)
}

func (s *playerService) GetPlayerDetail(ctx context.Context, playerId string) api.WebResponse[*dto.GetPlayerDetailResponse] {
	l := logger.LoggerNew()

	l.WithContext(ctx).Debug("[GetPlayerDetail].service: Started").Msg()

	result, err := s.playerRepository.GetPlayerDetail(ctx, playerId)
	if err != nil {
		l.WithContext(ctx).Error("error when get player detail").Attr("error", err).Msg()

		errMsg := constants.InternalServerError
		if errors.Is(err, constants.InvalidUUIDValue) {
			errMsg = constants.BadRequestDefault
		}
		if errors.Is(err, constants.ErrNotFoundRow) {
			errMsg = constants.NotFoundDefault
		}
		return api.ErrorResponse[*dto.GetPlayerDetailResponse](
			ctx,
			errMsg.GetMessage(),
			errMsg.GetCode(),
			errMsg.GetHttpCode(),
			nil,
		)
	}

	var positionCode, positionTitle string
	position := constants.DictPlayerPosition.GetValue(helper.GetStringPtrValue(result.Position))
	if position != nil {
		positionCode, positionTitle = position.Code, position.Title
	}

	response := &dto.GetPlayerDetailResponse{
		Id:     result.Id,
		TeamId: result.TeamId,
		Name:   helper.GetStringPtrValue(result.Name),
		Height: helper.GetFloat64PtrValue(result.Height),
		Weight: helper.GetFloat64PtrValue(result.Weight),
		Position: dto.PlayerPosition{
			Code:  positionCode,
			Title: positionTitle,
		},
		JerseyNumber: helper.GetIntPtrValue(result.JerseyNumber),
		CreatedAt:    result.CreatedAt,
		UpdatedAt:    result.UpdatedAt,
		IsDeleted:    result.IsDeleted,
		DeletedAt:    result.DeletedAt,
		Team: &dto.ListPlayerItemTeam{
			Id:   result.Team.Id,
			Name: helper.GetStringPtrValue(result.Team.Name),
		},
		GoalCount: result.GoalsCount,
	}

	l.WithContext(ctx).Debug("[GetPlayerDetail].service: Completed").Msg()
	return api.SuccessResponse(
		ctx,
		constants.SuccessDefault.GetMessage(),
		constants.SuccessDefault.GetCode(),
		constants.SuccessDefault.GetHttpCode(),
		response,
	)
}

func (s *playerService) DeletePlayer(ctx context.Context, playerId string) api.WebResponse[any] {
	l := logger.LoggerNew()

	l.WithContext(ctx).Debug("[DeletePlayer].service: Started").Msg()
	if err := s.playerRepository.DeletePlayer(ctx, playerId); err != nil {
		l.WithContext(ctx).Error("error when delete player").Attr("error", err).Msg()

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

	l.WithContext(ctx).Debug("[DeletePlayer].service: Completed").Msg()
	return api.SuccessResponse[any](
		ctx,
		constants.AcceptDefault.GetMessage(),
		constants.AcceptDefault.GetCode(),
		constants.AcceptDefault.GetHttpCode(),
		nil,
	)
}
