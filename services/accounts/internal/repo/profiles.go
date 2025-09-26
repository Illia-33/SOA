package repo

import (
	"soa-socialnetwork/services/accounts/internal/models"
	opt "soa-socialnetwork/services/common/option"
	"time"
)

type ProfilesRepo interface {
	GetByAccountId(models.AccountId) (models.ProfileData, error)
	GetByProfileId(models.ProfileId) (models.ProfileData, error)

	ResolveProfileId(models.ProfileId) (models.AccountId, error)
	ResolveAccountId(models.AccountId) (models.ProfileId, error)

	New(models.ProfileId, models.AccountId, models.RegistrationData) error
	Edit(models.ProfileId, EditedProfileData) error
	Delete(models.ProfileId) error
}

type EditedProfileData struct {
	Name     opt.Option[string]
	Surname  opt.Option[string]
	Bio      opt.Option[string]
	Birthday opt.Option[time.Time]
}
