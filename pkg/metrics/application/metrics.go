package application

import (
	"greye/pkg/metrics/domain/ports"
	"greye/pkg/type/domain/models"
)

func MetricFactory(roleType models.RoleType) ports.MetricPorts {
	switch roleType {
	case models.Application:
		metrics := NewApplicationMetric()
		return metrics
	case models.Cluster:
		metrics := NewClusterMetric()
		return metrics
	}
	panic("error creating metric")
}
