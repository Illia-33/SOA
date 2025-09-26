package postgres

import (
	"context"
	"errors"
	"log"
	"soa-socialnetwork/services/accounts/internal/models"
	"soa-socialnetwork/services/accounts/internal/repo"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
)

type profilesRepo struct {
	ctx   context.Context
	scope pgxScope
}

func (r profilesRepo) GetByAccountId(id models.AccountId) (models.ProfileData, error) {
	sql := `
	SELECT profile_id, name, surname, birthday, bio
	FROM profiles
	WHERE account_id = $1;
	`

	row := r.scope.QueryRow(r.ctx, sql, id)

	var (
		profileId  string
		name       string
		surname    string
		pgBirthday pgtype.Date
		pgBio      pgtype.Text
	)
	err := row.Scan(&profileId, &name, &surname, &pgBirthday, &pgBio)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return models.ProfileData{}, ErrorProfileNotFound{}
		}

		return models.ProfileData{}, err
	}

	var birthday time.Time
	if pgBirthday.Valid {
		birthday = pgBirthday.Time
	}

	var bio string
	if pgBio.Valid {
		bio = pgBio.String
	}

	return models.ProfileData{
		AccountId: id,
		ProfileId: models.ProfileId(profileId),
		Name:      name,
		Surname:   surname,
		Birthday:  birthday,
		Bio:       bio,
	}, nil
}

func (r profilesRepo) GetByProfileId(id models.ProfileId) (models.ProfileData, error) {
	sql := `
	SELECT account_id, name, surname, birthday, bio
	FROM profiles
	WHERE profile_id = $1;
	`

	row := r.scope.QueryRow(r.ctx, sql, id)

	var (
		accountId  int
		name       string
		surname    string
		pgBirthday pgtype.Date
		pgBio      pgtype.Text
	)
	err := row.Scan(&accountId, &name, &surname, &pgBirthday, &pgBio)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return models.ProfileData{}, ErrorProfileNotFound{}
		}

		return models.ProfileData{}, err
	}

	var birthday time.Time
	if pgBirthday.Valid {
		birthday = pgBirthday.Time
	}

	var bio string
	if pgBio.Valid {
		bio = pgBio.String
	}

	return models.ProfileData{
		AccountId: models.AccountId(accountId),
		ProfileId: id,
		Name:      name,
		Surname:   surname,
		Birthday:  birthday,
		Bio:       bio,
	}, nil
}

func (r profilesRepo) ResolveProfileId(id models.ProfileId) (models.AccountId, error) {
	sql := `
	SELECT account_id
	FROM profiles
	WHERE profile_id = $1;
	`

	row := r.scope.QueryRow(r.ctx, sql, id)

	var accountId int
	err := row.Scan(&accountId)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return -1, ErrorProfileNotFound{}
		}
		return -1, err
	}

	return models.AccountId(accountId), nil
}

func (r profilesRepo) ResolveAccountId(id models.AccountId) (models.ProfileId, error) {
	sql := `
	SELECT profile_id
	FROM profiles
	WHERE account_id = $1;
	`

	row := r.scope.QueryRow(r.ctx, sql, id)

	var profileId string
	err := row.Scan(&profileId)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return "", ErrorProfileNotFound{}
		}
		return "", err
	}

	return models.ProfileId(profileId), nil
}

func (r profilesRepo) New(profileId models.ProfileId, accountId models.AccountId, data models.RegistrationData) error {
	sql := `
	INSERT INTO profiles(account_id, profile_id, name, surname)
	VALUES ($1, $2::uuid, $3, $4);
	`

	_, err := r.scope.Exec(r.ctx, sql, accountId, profileId, data.Name, data.Surname)
	return err
}

func (r profilesRepo) Edit(id models.ProfileId, data repo.EditedProfileData) error {
	sql := `
	WITH cte AS (
		UPDATE profiles
		SET
			name = COALESCE($1, name),
			surname = COALESCE($2, surname),
			birthday = COALESCE($3, birthday),
			bio = COALESCE($4, bio)
		WHERE profile_id = $5
		RETURNING 1
	)
	SELECT count(*) FROM cte;
	`

	pgName := pgtype.Text{
		String: data.Name.Value,
		Valid:  data.Surname.HasValue,
	}
	pgSurname := pgtype.Text{
		String: data.Name.Value,
		Valid:  data.Surname.HasValue,
	}
	pgBirthday := pgtype.Date{
		Time:             data.Birthday.Value,
		InfinityModifier: pgtype.Finite,
		Valid:            data.Birthday.HasValue,
	}
	pgBio := pgtype.Text{
		String: data.Bio.Value,
		Valid:  data.Bio.HasValue,
	}
	row := r.scope.QueryRow(r.ctx, sql, pgName, pgSurname, pgBirthday, pgBio, id)

	var cnt int
	err := row.Scan(&cnt)
	if err != nil {
		return err
	}

	if cnt == 0 {
		return ErrorProfileNotFound{}
	}

	if cnt != 1 {
		log.Printf("warning: there %d users with profile id %s", cnt, id)
	}

	return nil
}

func (r profilesRepo) Delete(id models.ProfileId) error {
	sql := `
	WITH cte AS (
		DELETE FROM profiles
		WHERE profile_id = $1
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
		return ErrorProfileNotFound{}
	}

	if cnt != 1 {
		log.Printf("warning: there %d users with profile id %s", cnt, id)
	}

	return nil
}

type ErrorProfileNotFound struct{}

func (ErrorProfileNotFound) Error() string {
	return "profile not found"
}
