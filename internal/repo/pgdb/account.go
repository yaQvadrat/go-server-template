package pgdb

import (
	"app/internal/entity"
	"app/internal/repo/repoerrors"
	"app/pkg/postgres"
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

type AccountRepo struct {
	*postgres.Postgres
}

func NewAccountRepo(pg *postgres.Postgres) *AccountRepo {
	return &AccountRepo{pg}
}

func (r *AccountRepo) CreateAccount(ctx context.Context, account entity.Account) (int, error) {
	sql := `
		INSERT INTO account (username, password, role)
		VALUES ($1, $2, $3)
		RETURNING ID
	`

	var id int
	err := r.Pool.QueryRow(ctx, sql, account.Username, account.Password, account.Role).Scan(&id)
	if err != nil {
		var pgErr *pgconn.PgError
		if ok := errors.As(err, &pgErr); ok && pgErr.Code == UniqueViolationCode {
			return 0, repoerrors.ErrAlreadyExists
		}
		return 0, fmt.Errorf("pgdb - CreateAccount: %w", err)
	}

	return id, nil
}

func (r *AccountRepo) GetAccountById(ctx context.Context, id int) (entity.Account, error) {
	sql := `
		SELECT id, username, password, role, created_at
		FROM account
		WHERE id = $1
	`

	var account entity.Account
	err := r.Pool.QueryRow(ctx, sql, id).Scan(
		&account.Id,
		&account.Username,
		&account.Password,
		&account.Role,
		&account.CreatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return entity.Account{}, repoerrors.ErrNotFound
		}
		return entity.Account{}, fmt.Errorf("pgdb - GetAccountById: %w", err)
	}

	return account, nil
}

func (r *AccountRepo) GetAccountByUsername(ctx context.Context, username string) (entity.Account, error) {
	sql := `
		SELECT id, username, password, role, created_at
		FROM account
		WHERE username = $1
	`

	var account entity.Account
	err := r.Pool.QueryRow(ctx, sql, username).Scan(
		&account.Id,
		&account.Username,
		&account.Password,
		&account.Role,
		&account.CreatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return entity.Account{}, repoerrors.ErrNotFound
		}
		return entity.Account{}, fmt.Errorf("pgdb - GetAccountByUsername: %w", err)
	}

	return account, nil
}

func (r *AccountRepo) GetAccountByUsernameAndPassword(ctx context.Context, username, password string) (entity.Account, error) {
	sql := `
		SELECT id, username, password, role, created_at
		FROM account
		WHERE username = $1 AND password = $2
	`

	var account entity.Account
	err := r.Pool.QueryRow(ctx, sql, username, password).Scan(
		&account.Id,
		&account.Username,
		&account.Password,
		&account.Role,
		&account.CreatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return entity.Account{}, repoerrors.ErrNotFound
		}
		return entity.Account{}, fmt.Errorf("pgdb - GetAccountByUsername: %w", err)
	}

	return account, nil
}
