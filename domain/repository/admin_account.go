package repository

import (
	"context"
	"errors"

	"github.com/afrizalsebastian/football-team-management/constants"
	"github.com/afrizalsebastian/football-team-management/domain/dao"
	acctdb "github.com/afrizalsebastian/football-team-management/domain/repository/queries/db"
	"github.com/afrizalsebastian/football-team-management/module/logger"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type IAdminAccountRepository interface {
	CreateAccount(ctx context.Context, account *dao.AdminAccount) error
	GetAdminAccount(ctx context.Context, username string) (*dao.AdminAccount, error)
}

type adminAccountRepository struct {
	db *acctdb.Queries
}

func NewAdminAccountRepository(
	pool *pgxpool.Pool,
) IAdminAccountRepository {
	return &adminAccountRepository{
		db: acctdb.New(pool),
	}
}

func (r *adminAccountRepository) CreateAccount(ctx context.Context, account *dao.AdminAccount) error {
	l := logger.LoggerNew()

	l.WithContext(ctx).Debug("[CreateAccount].domain: Started").Msg()

	row, err := r.db.CreateAdminAccount(ctx, &acctdb.CreateAdminAccountParams{
		Username: account.Username,
		Password: account.Password,
	})

	if err != nil {
		l.WithContext(ctx).Error("error when create account").Attr("error", err).Msg()

		if isDuplicateError(err) {
			return constants.DuplicateRow
		}
		return err
	}

	account.Id = row.String()
	return nil
}

func (r *adminAccountRepository) GetAdminAccount(ctx context.Context, username string) (*dao.AdminAccount, error) {
	l := logger.LoggerNew()

	l.WithContext(ctx).Debug("[GetPlayerDetail].domain: Started").Msg()

	row, err := r.db.GetAdmingByUsername(ctx, username)
	if err != nil {
		l.WithContext(ctx).Error("error when get account").
			Attr("error", err).Attr("username", username).Msg()
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, constants.ErrNotFoundRow
		}
		return nil, err
	}

	l.WithContext(ctx).Debug("[GetPlayerDetail].domain: Completed").Msg()
	return &dao.AdminAccount{
		Id:       row.ID.String(),
		Username: row.Username,
		Password: row.Password,
	}, nil
}
