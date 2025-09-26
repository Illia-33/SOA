package repo

import (
	"soa-socialnetwork/services/accounts/internal/models"
)

type AccountsRepo interface {
	CheckPasswordByLogin(login string, password string) (models.AccountParams, error)
	CheckPasswordByEmail(email string, password string) (models.AccountParams, error)
	CheckPasswordByPhoneNumber(phoneNumber string, password string) (models.AccountParams, error)

	New(models.RegistrationData) (models.AccountId, error)
	Delete(models.AccountId) error
}
