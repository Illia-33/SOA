package postgres

import (
	"context"
	"errors"
	"fmt"
	"log"
	"soa-socialnetwork/services/accounts/internal/models"
	"soa-socialnetwork/services/accounts/internal/storage/postgres/errs"

	"github.com/jackc/pgx/v5"
)

type accountsRepo struct {
	ctx   context.Context
	scope pgxScope
}

func (r accountsRepo) CheckPasswordByLogin(login string, password string) (models.AccountParams, error) {
	params, err := r.fetchAccountParams("login", login)
	if err != nil {
		return models.AccountParams{}, err
	}

	if params.password != password {
		return models.AccountParams{}, errs.PasswordsDoNotMatch{}
	}

	return models.AccountParams{Id: models.AccountId(params.id)}, nil
}

func (r accountsRepo) CheckPasswordByEmail(email string, password string) (models.AccountParams, error) {
	params, err := r.fetchAccountParams("email", email)
	if err != nil {
		return models.AccountParams{}, err
	}

	if params.password != password {
		return models.AccountParams{}, errs.PasswordsDoNotMatch{}
	}

	return models.AccountParams{Id: models.AccountId(params.id)}, nil
}

func (r accountsRepo) CheckPasswordByPhoneNumber(phoneNumber string, password string) (models.AccountParams, error) {
	params, err := r.fetchAccountParams("phone_number", phoneNumber)
	if err != nil {
		return models.AccountParams{}, err
	}

	if params.password != password {
		return models.AccountParams{}, errs.PasswordsDoNotMatch{}
	}

	return models.AccountParams{Id: models.AccountId(params.id)}, nil
}

func (r accountsRepo) New(data models.RegistrationData) (models.AccountId, error) {
	sql := `
	INSERT INTO accounts (login, password, email, phone_number)
	VALUES ($1, $2, $3, $4)
	RETURNING id
	`
	row := r.scope.QueryRow(r.ctx, sql, data.Login, data.Password, data.Email, data.PhoneNumber)

	var id int
	err := row.Scan(&id)
	if err != nil {
		return -1, err
	}

	return models.AccountId(id), nil
}

func (r accountsRepo) Delete(id models.AccountId) error {
	sql := `
	WITH cte AS (
		DELETE FROM accounts
		WHERE id = $1
		RETURNING 1
	)
	SELECT count(*) FROM cte;
	`

	row := r.scope.QueryRow(r.ctx, sql, id)

	var cnt int
	err := row.Scan(&cnt)
	if err != nil {
		return err
	}

	if cnt == 0 {
		return errs.AccountNotFound{}
	}

	if cnt != 1 {
		log.Printf("warning: there more than one users with account id %d", id)
	}

	return nil
}

type accountParams struct {
	id       int
	password string
}

func (r accountsRepo) fetchAccountParams(colName string, colValue string) (accountParams, error) {
	sql := fmt.Sprintf(`
	SELECT id, password
	FROM accounts
	WHERE %s = $1;
	`, colName)

	row := r.scope.QueryRow(r.ctx, sql, colValue)

	var params accountParams
	err := row.Scan(&params.id, &params.password)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return accountParams{}, errs.AccountNotFound{}
		}

		return accountParams{}, err
	}

	return params, nil
}
