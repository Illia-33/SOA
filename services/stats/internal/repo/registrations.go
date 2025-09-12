package repo

import "soa-socialnetwork/services/stats/pkg/models"

type RegistrationsRepo interface {
	Put(...models.RegistrationEvent) error
}
