package postgres

import (
	"context"
	"errors"
	"soa-socialnetwork/services/accounts/internal/models"
	"soa-socialnetwork/services/accounts/internal/repo"
	"time"

	"github.com/jackc/pgx/v5"
)

type apiTokensRepo struct {
	ctx   context.Context
	scope pgxScope
}

func (r apiTokensRepo) Put(token models.ApiToken, params repo.ApiTokenParams) (time.Time, error) {
	sql := `
	INSERT INTO api_tokens(account_id, token, valid_until, read_access, write_access)
	VALUES ($1, $2, NOW() + $3, $4, $5)
	RETURNING valid_until;
	`

	row := r.scope.QueryRow(r.ctx, sql, params.AccountId, token, params.Ttl, params.ReadAccess, params.WriteAccess)

	var validUntil time.Time
	err := row.Scan(&validUntil)
	if err != nil {
		return time.Time{}, err
	}

	return validUntil, nil
}

func (r apiTokensRepo) Get(token models.ApiToken) (models.ApiTokenData, error) {
	sql := `
	SELECT account_id, read_access, write_access, created_at, valid_until
	FROM api_tokens
	WHERE token = $1;
	`

	row := r.scope.QueryRow(r.ctx, sql, token)

	var data models.ApiTokenData
	err := row.Scan(&data.AccountId, &data.ReadAccess, &data.WriteAccess, &data.CreatedAt, &data.ValidUntil)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return models.ApiTokenData{}, ErrorTokenNotFound{}
		}

		return models.ApiTokenData{}, err
	}
	data.Token = token

	return data, nil
}

type ErrorTokenNotFound struct{}

func (ErrorTokenNotFound) Error() string {
	return "token not found"
}
