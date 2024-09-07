package service

import (
	"app/internal/repo"
	"app/pkg/hasher"
	"context"
	"time"
)

type (
	CreateAccountInput struct {
		Username string
		Password string
		Role string
	}

	GenerateTokenInput struct {
		Username string
		Password string
	}

	ParsedToken struct {
		AccId int
		Role   string
	}
)

type Account interface {
	CreateAccount(ctx context.Context, input CreateAccountInput) (int, error)
	GenerateToken(ctx context.Context, input GenerateTokenInput) (string, error)
	ParseToken(token string) (ParsedToken, error)
}

type (
	Services struct {
		Account
	}

	ServicesDependencies struct {
		Repos    *repo.Repositories
		Hasher   hasher.Hasher
		TokenTTL time.Duration
		SignKey  string
	}
)

func NewServices(d ServicesDependencies) *Services {
	return &Services{
		Account: NewAccountService(d.Repos.Account, d.Hasher, d.TokenTTL, d.SignKey),
	}
}
