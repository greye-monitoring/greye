package application

import "greye/pkg/config/domain/models"

func (c *ConfigService) GetConfig() (*models.Config, error) {
	if c.configuration == nil {
		if err := c.Config(); err != nil {
			return nil, err
		}
	}
	if err := c.validator.Struct(c.configuration); err != nil {
		return nil, err
	}
	return c.configuration, nil
}
