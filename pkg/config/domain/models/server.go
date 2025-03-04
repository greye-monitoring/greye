package models

import (
	"errors"
	"greye/pkg/validator/domain/ports"
)

type Server struct {
	LogLevel        string `json:"logLevel"`
	Port            int    `json:"port"`
	TlsPort         int    `json:"tlsPort"`
	NumberGreye     int    `json:"numberGreye"`
	ApplicationUrl  string `json:"applicationUrl"`
	ApplicationName string `json:"applicationName"`
	ServiceHAName   string `json:"serviceHAName"`
	ServerUrl       string `json:"serverUrl"`
}

var _ ports.Evaluable = (*Server)(nil)

func (s *Server) Validate() error {
	if s.Port == 0 {
		return errors.New("the port is required")
	}
	return nil
}
