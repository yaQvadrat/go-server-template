package repo

import (
	"app/internal/entity"
	"app/internal/repo/pgdb"
	"app/pkg/postgres"
	"context"
)

type Account interface {
	CreateAccount(ctx context.Context, account entity.Account) (int, error)
	GetAccountById(ctx context.Context, id int) (entity.Account, error)
	GetAccountByUsername(ctx context.Context, username string) (entity.Account, error)
	GetAccountByUsernameAndPassword(ctx context.Context, username, password string) (entity.Account, error)
}

type Repositories struct {
	Account
}

func NewPostgresRepo(pg *postgres.Postgres) *Repositories {
	return &Repositories{
		Account: pgdb.NewAccountRepo(pg),
	}
}
