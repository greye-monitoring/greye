package handlers

import (
	"greye/internal/application/domain/ports"
	clientHttp "greye/pkg/client/domain/ports"
	configPort "greye/pkg/config/domain/ports"
	logrus "greye/pkg/logging/domain/ports"
	schedulerPort "greye/pkg/scheduler/domain/ports"
	valPort "greye/pkg/validator/domain/ports"
)

type ApplicationHdl struct {
	config        configPort.ConfigApplication
	validator     valPort.ValidatorApplication
	logger        logrus.LoggerApplication
	scheduler     schedulerPort.Operation
	http          clientHttp.HttpMethod
	schedulerData ports.SchedulerService
}

var _ ports.ApiExposed = (*ApplicationHdl)(nil)

func NewApiExposedHdl(validator valPort.ValidatorApplication, logger logrus.LoggerApplication, httpCLient clientHttp.HttpMethod,
	schedulerHandler schedulerPort.Operation, schedulerApp ports.SchedulerService, config configPort.ConfigApplication) *ApplicationHdl {

	return &ApplicationHdl{
		validator:     validator,
		logger:        logger,
		http:          httpCLient,
		scheduler:     schedulerHandler,
		schedulerData: schedulerApp,
		config:        config,
	}
}
