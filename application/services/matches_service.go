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
