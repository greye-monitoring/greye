package models

import "greye/pkg/validator/domain/ports"

type Config struct {
	App          string                 `json:"app"`
	Server       Server                 `json:"server"`
	Application  Application            `json:"application"`
	Cluster      Cluster                `json:"cluster"`
	Notification map[string]interface{} `json:"notification"`
	Protocol     []string               `json:"protocol"`
	JWT          string                 `json:"jwt"`
}

var _ ports.Evaluable = (*Config)(nil)

func (c *Config) Validate() error {
	if err := c.Server.Validate(); err != nil {
		return err
	}

	return nil
}
