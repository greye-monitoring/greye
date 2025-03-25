package application

import (
	"greye/pkg/config/domain/models"
	"greye/pkg/config/domain/ports"
	validator "greye/pkg/validator/domain/ports"
)

type ConfigService struct {
	clusters      []string
	repository    ports.ConfigRepository
	configuration *models.Config
	validator     validator.ValidatorApplication
	//logger        loggerModels.LoggerApplication
}

var _ ports.ConfigApplication = (*ConfigService)(nil)

func NewConfigService(
	repository ports.ConfigRepository,
	validator validator.ValidatorApplication,
) *ConfigService {
	return &ConfigService{
		repository: repository,
		validator:  validator,
	}
}
