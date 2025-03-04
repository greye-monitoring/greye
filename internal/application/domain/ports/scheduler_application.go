package ports

import (
	"greye/internal/application/domain/models"
)

type SchedulerService interface {
	MonitorApplication(app *models.SchedulerApplication, startupPhase bool) error
	GetApplication(url string) (map[string]models.SchedulerApplication, error)
	DeleteApplicationFromUrl(url string) error
}
