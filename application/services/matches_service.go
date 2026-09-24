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

type IMatchesService interface {
	CreateMatch(ctx context.Context, request *dto.CreateMatchRequest) api.WebResponse[*dto.CreateMatchResponse]
	GetListMatches(ctx context.Context, teamId *string) api.WebResponse[[]dto.GetMatchResponse]
	RescheduleMatch(ctx context.Context, matchId string, request *dto.RescheduleMatchRequest) api.WebResponse[any]
	DeleteMatch(ctx context.Context, matchId string) api.WebResponse[any]
}

type matchesService struct {
	matchRepository repository.IMatchesRepository
}

func NewMatchService(
	matchRepository repository.IMatchesRepository,
) IMatchesService {
	return &matchesService{
		matchRepository: matchRepository,
	}
}

func (s *matchesService) CreateMatch(ctx context.Context, request *dto.CreateMatchRequest) api.WebResponse[*dto.CreateMatchResponse] {
	l := logger.LoggerNew()

	l.WithContext(ctx).Debug("[CreateMatch].service: Started").Msg()
	match := &dao.Matches{
		HomeTeamId: request.HomeTeamId,
		AwayTeamId: request.AwayTeamId,
		MatchDate:  helper.StringPtr(request.Date),
		MatchTime:  helper.StringPtr(request.Time),
	}

	if err := s.matchRepository.CreateMatches(ctx, match); err != nil {
		l.WithContext(ctx).Error("error when create match").Attr("error", err).Msg()

		errMsg := constants.InternalServerError
		if errors.Is(err, constants.InvalidUUIDValue) {
			errMsg = constants.BadRequestDefault
		}

		return api.ErrorResponse[*dto.CreateMatchResponse](
			ctx,
			errMsg.GetMessage(),
			errMsg.GetCode(),
			errMsg.GetHttpCode(),
			nil,
		)
	}

	response := &dto.CreateMatchResponse{
		Id:         match.Id,
		HomeTeamId: match.HomeTeamId,
		AwayTeamId: match.AwayTeamId,
		Date:       helper.GetStringPtrValue(match.MatchDate),
		Time:       helper.GetStringPtrValue(match.MatchTime),
	}

	l.WithContext(ctx).Debug("[CreateMatch].service: Completed").Msg()
	return api.SuccessResponse(
		ctx,
		constants.CreatedDefault.GetMessage(),
		constants.CreatedDefault.GetCode(),
		constants.CreatedDefault.GetHttpCode(),
		response,
	)
}

func (s *matchesService) GetListMatches(ctx context.Context, teamId *string) api.WebResponse[[]dto.GetMatchResponse] {
	l := logger.LoggerNew()

	l.WithContext(ctx).Debug("[GetListMatches].service: Started").Msg()
	result, err := s.matchRepository.GetListMatch(ctx, teamId)
	if err != nil {
		l.WithContext(ctx).Error("error when get list match").Attr("error", err).Msg()
		return api.ErrorResponse[[]dto.GetMatchResponse](
			ctx,
			constants.InternalServerError.GetMessage(),
			constants.InternalServerError.GetCode(),
			constants.InternalServerError.GetHttpCode(),
			nil,
		)
	}

	if len(result) == 0 {
		api.SuccessResponse[[]dto.GetMatchResponse](
			ctx,
			constants.SuccesssWithEmptyList.GetMessage(),
			constants.SuccesssWithEmptyList.GetCode(),
			constants.SuccesssWithEmptyList.GetHttpCode(),
			nil,
		)
	}

	response := make([]dto.GetMatchResponse, 0)
	for _, r := range result {
		response = append(response, dto.GetMatchResponse{
			Id:         r.Id,
			HomeTeamId: r.HomeTeamId,
			AwayTeamId: r.AwayTeamId,
			Date:       helper.GetStringPtrValue(r.MatchDate),
			Time:       helper.GetStringPtrValue(r.MatchTime),
			HomeTeam: dto.MatchTeam{
				Id:   r.HomeTeam.Id,
				Name: helper.GetStringPtrValue(r.HomeTeam.Name),
				Logo: helper.GetStringPtrValue(r.HomeTeam.Logo),
			},
			AwayTeam: dto.MatchTeam{
				Id:   r.AwayTeam.Id,
				Name: helper.GetStringPtrValue(r.AwayTeam.Name),
				Logo: helper.GetStringPtrValue(r.AwayTeam.Logo),
			},
		})
	}

	l.WithContext(ctx).Debug("[GetListMatches].service: Completed").Msg()
	return api.SuccessResponse(
		ctx,
		constants.SuccessDefault.GetMessage(),
		constants.SuccessDefault.GetCode(),
		constants.SuccessDefault.GetHttpCode(),
		response,
	)
}

func (s *matchesService) RescheduleMatch(ctx context.Context, matchId string, request *dto.RescheduleMatchRequest) api.WebResponse[any] {
	l := logger.LoggerNew()

	l.WithContext(ctx).Debug("[RescheduleMatch].service: Started").Msg()
	match := &dao.Matches{
		Id:        matchId,
		MatchDate: helper.StringPtr(request.Date),
		MatchTime: helper.StringPtr(request.Time),
	}
	if err := s.matchRepository.RescheduleMatch(ctx, match); err != nil {
		l.WithContext(ctx).Error("error when reschedule match").Attr("error", err).Attr("match_id", matchId).Msg()

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

	return api.SuccessResponse[any](
		ctx,
		constants.AcceptDefault.GetMessage(),
		constants.AcceptDefault.GetCode(),
		constants.AcceptDefault.GetHttpCode(),
		nil,
	)
}

func (s *matchesService) DeleteMatch(ctx context.Context, matchId string) api.WebResponse[any] {
	l := logger.LoggerNew()

	l.WithContext(ctx).Debug("[DeleteMatch].service: Started").Msg()
	if err := s.matchRepository.DeleteMatch(ctx, matchId); err != nil {
		l.WithContext(ctx).Error("error when delete match").Attr("error", err).Attr("match_id", matchId).Msg()

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

	return api.SuccessResponse[any](
		ctx,
		constants.AcceptDefault.GetMessage(),
		constants.AcceptDefault.GetCode(),
		constants.AcceptDefault.GetHttpCode(),
		nil,
	)
}
