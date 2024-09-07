package service

import (
	"app/internal/entity"
	"app/internal/repo"
	"app/internal/repo/repoerrors"
	"app/pkg/hasher"
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt"
	log "github.com/sirupsen/logrus"
)

type AccountService struct {
	accountRepo    repo.Account
	passwordHasher hasher.Hasher
	tokenTTL       time.Duration
	signKey        string
}

type TokenClaims struct {
	jwt.StandardClaims
	AccId int
	Role  string
}

func NewAccountService(repo repo.Account, hasher hasher.Hasher, tokenTTL time.Duration, signKey string) *AccountService {
	return &AccountService{
		accountRepo:    repo,
		passwordHasher: hasher,
		tokenTTL:       tokenTTL,
		signKey:        signKey,
	}
}

func (s *AccountService) CreateAccount(ctx context.Context, input CreateAccountInput) (int, error) {
	account := entity.Account{
		Username: input.Username,
		Password: s.passwordHasher.Hash(input.Password),
		Role: input.Role,
	}

	id, err := s.accountRepo.CreateAccount(ctx, account)
	if err != nil {
		if errors.Is(err, repoerrors.ErrAlreadyExists) {
			return 0, ErrAccountAlreadyExists
		}
		log.Error(fmt.Errorf("service - AccountService.CreateAccount: %w", err))
	}

	return id, nil
}

func (s *AccountService) GenerateToken(ctx context.Context, input GenerateTokenInput) (string, error) {
	acc, err := s.accountRepo.GetAccountByUsernameAndPassword(ctx, input.Username, input.Password)
	if err != nil {
		if errors.Is(err, repoerrors.ErrNotFound) {
			return "", ErrAccountNotFound
		}
		log.Error(fmt.Errorf("service - AccountService.GenerateToken: %w", err))
		return "", ErrCannotGetAccount
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &TokenClaims{
		StandardClaims: jwt.StandardClaims{
			IssuedAt:  time.Now().Unix(),
			ExpiresAt: time.Now().Add(12 * time.Hour).Unix(),
		},
		AccId: acc.Id,
		Role:  acc.Role,
	})

	tokenString, err := token.SignedString([]byte(s.signKey))
	if err != nil {
		log.Error(fmt.Errorf("service - AccountService.GenerateToken: %w", err))
		return "", ErrCannotSignJWT
	}

	return tokenString, nil
}

func (s *AccountService) ParseToken(tokenString string) (ParsedToken, error) {
	token, err := jwt.ParseWithClaims(tokenString, &TokenClaims{}, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return ParsedToken{}, fmt.Errorf("incorrect signing method")
		}

		return []byte(s.signKey), nil
	})

	if err != nil {
		return ParsedToken{}, ErrCannotParseJWT
	}

	claims, ok := token.Claims.(*TokenClaims)
	if !ok {
		return ParsedToken{}, ErrCannotParseJWT
	}

	return ParsedToken{
		AccId: claims.AccId,
		Role:  claims.Role,
	}, nil
}
