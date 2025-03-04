package models

import (
	"errors"
	"greye/pkg/validator/domain/ports"
	"time"
)

type Application struct {
	IntervalSeconds   time.Duration `json:"intervalSeconds"`
	Protocol          string        `json:"protocol"`
	Method            string        `json:"method"`
	MaxFailedRequests int           `json:"maxFailedRequests"`
	TimeoutSeconds    time.Duration `json:"timeoutSeconds"`
}

var _ ports.Evaluable = (*Application)(nil)

func (s *Application) Validate() error {
	if s.MaxFailedRequests == 0 {
		return errors.New("MaxFailedRequests must be greater than 0")
	}
	if s.Method == "" {
		return errors.New("Method must be provided")
	}
	if s.Protocol == "" {
		return errors.New("Protocol must be provided")
	}
	if s.IntervalSeconds < 30 {
		return errors.New("IntervalSeconds must be greater than 30")
	}
	if s.TimeoutSeconds < 1 {
		return errors.New("TimeoutSeconds must be greater or equal to 1 second")
	}
	return nil
}
