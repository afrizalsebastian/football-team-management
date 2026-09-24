package services

import (
	"context"
	"errors"

	"github.com/afrizalsebastian/football-team-management/api"
	"github.com/afrizalsebastian/football-team-management/constants"
	"github.com/afrizalsebastian/football-team-management/domain/repository"
	"github.com/afrizalsebastian/football-team-management/module/logger"
)

type IGoalService interface {
	DeleteGoals(ctx context.Context, goalId string) api.WebResponse[any]
}

type goalService struct {
	goalRepository repository.IGoalsRepository
}

func NewGoalService(
	goalRepository repository.IGoalsRepository,
) IGoalService {
	return &goalService{
		goalRepository: goalRepository,
	}
}

func (s *goalService) DeleteGoals(ctx context.Context, goalId string) api.WebResponse[any] {
	l := logger.LoggerNew()

	l.WithContext(ctx).Debug("[DeleteGoals].service: Started").Msg()
	if err := s.goalRepository.DeleteGoals(ctx, goalId); err != nil {
		l.WithContext(ctx).Error("error when soft delete goal").Attr("error", err).Msg()

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

	l.WithContext(ctx).Debug("[DeleteGoals].service: Completed").Msg()
	return api.SuccessResponse[any](
		ctx,
		constants.AcceptDefault.GetMessage(),
		constants.AcceptDefault.GetCode(),
		constants.AcceptDefault.GetHttpCode(),
		nil,
	)
}
