package ports

import "greye/pkg/authentication/domain/models"

type Authentication interface {
	GetAuthorization(data models.AuthenticationData) (string, error)
}
