package application

import (
	"fmt"
	"greye/pkg/client/domain/ports"
	logger "greye/pkg/logging/domain/ports"
)

func PrtocolFactory(protocol string, logger logger.LoggerApplication) (ports.MonitoringMethod, error) {
	switch protocol {
	case "http":
		httpMonitoring := NewHttpMonitoring(logger)
		return httpMonitoring, nil
	default:
		panic(fmt.Sprintf("error creating %s protocol", protocol))
	}
}
