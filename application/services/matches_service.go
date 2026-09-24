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
	GetListMatches(ctx context.Context, teamId *string, matchType *string) api.WebResponse[[]dto.GetMatchResponse]
	RescheduleMatch(ctx context.Context, matchId string, request *dto.RescheduleMatchRequest) api.WebResponse[any]
	DeleteMatch(ctx context.Context, matchId string) api.WebResponse[any]
	CreateMatchGoal(ctx context.Context, matchId string, request *dto.MatchGoal) api.WebResponse[any]
	MatchGoalList(ctx context.Context, matchId string) api.WebResponse[*dto.MatchGoalsResponse]
	MatchFullTime(ctx context.Context, matchId string) api.WebResponse[any]
}

type matchesService struct {
	matchRepository repository.IMatchesRepository
	goalRepository  repository.IGoalsRepository
}

func NewMatchService(
	matchRepository repository.IMatchesRepository,
	goalRepository repository.IGoalsRepository,
) IMatchesService {
	return &matchesService{
		matchRepository: matchRepository,
		goalRepository:  goalRepository,
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

func (s *matchesService) GetListMatches(ctx context.Context, teamId *string, matchType *string) api.WebResponse[[]dto.GetMatchResponse] {
	l := logger.LoggerNew()

	l.WithContext(ctx).Debug("[GetListMatches].service: Started").Msg()
	result, err := s.matchRepository.GetListMatch(ctx, teamId, matchType)
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
		return api.SuccessResponse[[]dto.GetMatchResponse](
			ctx,
			constants.SuccesssWithEmptyList.GetMessage(),
			constants.SuccesssWithEmptyList.GetCode(),
			constants.SuccesssWithEmptyList.GetHttpCode(),
			nil,
		)
	}

	response := make([]dto.GetMatchResponse, 0)
	for _, r := range result {
		matchGoals := make([]dto.MatchGoalListItem, 0)
		for _, g := range r.Goals {
			matchGoals = append(matchGoals, dto.MatchGoalListItem{
				Id:         g.Id,
				GoalMinute: g.GoalMinute,
				Player: dto.GetListPlayerItem{
					Id:           g.Player.Id,
					Name:         helper.GetStringPtrValue(g.Player.Name),
					JerseyNumber: helper.GetIntPtrValue(g.Player.JerseyNumber),
					Team: dto.ListPlayerItemTeam{
						Id: g.Player.TeamId,
					},
				},
			})
		}

		response = append(response, dto.GetMatchResponse{
			Id:         r.Id,
			HomeTeamId: r.HomeTeamId,
			AwayTeamId: r.AwayTeamId,
			Date:       helper.GetStringPtrValue(r.MatchDate),
			Time:       helper.GetStringPtrValue(r.MatchTime),
			HomeScore:  r.HomeScore,
			AwayScore:  r.AwayScore,
			Status:     constants.MapMatchStatus[r.Status],
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
			Goals: matchGoals,
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

func (s *matchesService) CreateMatchGoal(ctx context.Context, matchId string, request *dto.MatchGoal) api.WebResponse[any] {
	l := logger.LoggerNew()

	l.WithContext(ctx).Debug("[CreateMatchGoal].service: Started").Msg()
	goals := &dao.Goals{
		MatchId:    matchId,
		PlayerId:   request.PlayerId,
		GoalMinute: request.GoalMinute,
	}
	if err := s.goalRepository.MatchGoal(ctx, goals); err != nil {
		l.WithContext(ctx).Error("error when create goals").Attr("error", err).Attr("match_id", matchId).Msg()

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

	l.WithContext(ctx).Debug("[CreateMatchGoal].service: Completed").Msg()
	return api.SuccessResponse[any](
		ctx,
		constants.CreatedDefault.GetMessage(),
		constants.CreatedDefault.GetCode(),
		constants.CreatedDefault.GetHttpCode(),
		nil,
	)
}

func (s *matchesService) MatchGoalList(ctx context.Context, matchId string) api.WebResponse[*dto.MatchGoalsResponse] {
	l := logger.LoggerNew()

	l.WithContext(ctx).Debug("[MatchGoalList].service: Started").Msg()

	result, match, err := s.goalRepository.GetMatchGoalList(ctx, matchId)
	if err != nil {
		l.WithContext(ctx).Error("error when get match goal list").Attr("error", err).Attr("match_id", matchId).Msg()

		errMsg := constants.InternalServerError
		if errors.Is(err, constants.InvalidUUIDValue) {
			errMsg = constants.BadRequestDefault
		}

		return api.ErrorResponse[*dto.MatchGoalsResponse](
			ctx,
			errMsg.GetMessage(),
			errMsg.GetCode(),
			errMsg.GetHttpCode(),
			nil,
		)
	}

	if len(result) == 0 {
		return api.SuccessResponse[*dto.MatchGoalsResponse](
			ctx,
			constants.SuccesssWithEmptyList.GetMessage(),
			constants.SuccesssWithEmptyList.GetCode(),
			constants.SuccesssWithEmptyList.GetHttpCode(),
			nil,
		)
	}

	goals := make([]dto.MatchGoalListItem, 0)
	for _, r := range result {
		goals = append(goals, dto.MatchGoalListItem{
			Id:         r.Id,
			GoalMinute: r.GoalMinute,
			Player: dto.GetListPlayerItem{
				Id:           r.Player.Id,
				Name:         helper.GetStringPtrValue(r.Player.Name),
				JerseyNumber: helper.GetIntPtrValue(r.Player.JerseyNumber),
				Team: dto.ListPlayerItemTeam{
					Id:   r.Player.Team.Id,
					Name: helper.GetStringPtrValue(r.Player.Team.Name),
				},
			},
		})
	}

	response := &dto.MatchGoalsResponse{
		Goals:   goals,
		MatchId: match.Id,
		HomeTeam: dto.ListPlayerItemTeam{
			Id:   match.HomeTeam.Id,
			Name: helper.GetStringPtrValue(match.HomeTeam.Name),
		},
		AwayTeam: dto.ListPlayerItemTeam{
			Id:   match.AwayTeam.Id,
			Name: helper.GetStringPtrValue(match.AwayTeam.Name),
		},
	}

	l.WithContext(ctx).Debug("[CreateMatchGoal].service: MatchGoalList").Msg()
	return api.SuccessResponse(
		ctx,
		constants.SuccessDefault.GetMessage(),
		constants.SuccessDefault.GetCode(),
		constants.SuccessDefault.GetHttpCode(),
		response,
	)
}

func (s *matchesService) MatchFullTime(ctx context.Context, matchId string) api.WebResponse[any] {
	l := logger.LoggerNew()

	l.WithContext(ctx).Debug("[MatchFullTime].service: Started").Msg()
	if err := s.matchRepository.MatchFullTime(ctx, matchId); err != nil {
		l.WithContext(ctx).Error("error when set match results").Attr("error", err).Attr("match_id", matchId).Msg()

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

	l.WithContext(ctx).Debug("[MatchFullTime].service: Completed").Msg()
	return api.SuccessResponse[any](
		ctx,
		constants.AcceptDefault.GetMessage(),
		constants.AcceptDefault.GetCode(),
		constants.AcceptDefault.GetHttpCode(),
		nil,
	)
}
