package models

import (
	"errors"
	"greye/pkg/validator/domain/ports"
	"time"
)

type Cluster struct {
	IntervalSeconds   time.Duration `json:"intervalSeconds"`
	TimeoutSeconds    int           `json:"timeoutSeconds"`
	MaxFailedRequests int           `json:"maxFailedRequests"`
	MyIp              string        `json:"myIp"`
	ClusterIp         []string      `json:"ip"`
}

var _ ports.Evaluable = (*Application)(nil)

func (s *Cluster) Cluster() error {

	if s.IntervalSeconds < 30 {
		return errors.New("IntervalSeconds must be greater than 30")
	}

	if s.TimeoutSeconds < 5 {
		return errors.New("TimeoutSeconds must be greater than 5")
	}

	if s.MaxFailedRequests < 4 {
		return errors.New("MaxFailedRequests must be greater than 0")
	}

	if s.MyIp == "" {
		return errors.New("MyIp must be provided")
	}
	return nil
}
