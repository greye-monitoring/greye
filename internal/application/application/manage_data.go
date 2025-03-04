package application

import (
	"fmt"
	"greye/internal/application/domain/models"
	annotations "greye/pkg/annotations/domain/models"
	modelsHttp "greye/pkg/client/domain/models"
	clientApp "greye/pkg/client/domain/ports"
	ports2 "greye/pkg/notification/domain/ports"
	v1 "k8s.io/api/core/v1"
	netUrl "net/url"
	"os"
	"regexp"
	"time"
)

func ConvertClusterInfoToResponse(ci interface{}) models.SchedulerApplication {
	response := models.SchedulerApplication{}

	if v, ok := ci.(models.SchedulerApplication); ok {
		response = v
	}

	return response
}

func (s *Scheduler) ReadApplications() map[string]models.SchedulerApplication {
	copiedMap := make(map[string]models.SchedulerApplication)
	s.applications.Range(func(key, value any) bool {
		copiedMap[key.(string)] = value.(models.SchedulerApplication)
		return true
	})

	return copiedMap
}

// todo models.schedulerapplication dovrebbe ritornare una copia!
func (s *Scheduler) ReadFromApplicationMap(key string) (models.SchedulerApplication, bool) {

	application, exist := s.applications.Load(key)
	var res models.SchedulerApplication
	if v, ok := application.(models.SchedulerApplication); ok && exist {
		res = v
	}

	return res, exist
}

func (s *Scheduler) WriteToApplicationMap(key string, applications models.SchedulerApplication) {
	s.applications.Store(key, applications)
}

func (s *Scheduler) DeleteFromApplication(url string) {
	s.applications.Delete(url)
}

func (s *Scheduler) ReadAlarms() map[string]ports2.Sender {
	s.RLock()
	defer s.RUnlock()
	return s.alarms
}

func (s *Scheduler) ReadFromClient(key string) (clientApp.MonitoringMethod, bool) {
	s.RLock()
	defer s.RUnlock()
	app, exist := s.client[key]
	return app, exist
}

func (s *Scheduler) GetApplication(url string) (map[string]models.SchedulerApplication, error) {
	if url != "" {
		app, exists := s.ReadFromApplicationMap(url)
		if !exists {
			return map[string]models.SchedulerApplication{}, nil
		}
		return map[string]models.SchedulerApplication{url: app}, nil
	}
	return s.ReadApplications(), nil
}

func (s *Scheduler) deleteApplication(app models.SchedulerApplication) error {
	url := app.Host
	if url == "" {
		return nil
	}
	hostname := s.getMyHostname()
	keyName := app.ScheduledApplication

	s.logger.Info("hostname: %s", hostname)
	s.logger.Info("keyname: %s", keyName)
	s.greyesStatus[keyName] = s.greyesStatus[keyName] - 1
	s.metrics.DeleteMetrics(url)
	if keyName == hostname {
		if app.Quit == nil || app.Ticker == nil {
			s.logger.Error("The application %s is not monitored", url)
			return nil
		}
		s.logger.Info("The application is mine")
		app.Quit <- struct{}{}
	} else {
		s.logger.Info("The application is monitored by %s", keyName)
		encodedHost := netUrl.QueryEscape(url)

		deleteRequest := &modelsHttp.HttpRequest{
			Name:     keyName,
			Host:     keyName,
			Timeout:  5 * time.Second,
			Protocol: "http",
			Path:     "/api/v1/application/monitor/" + encodedHost,
			Method:   "DELETE",
		}

		_, err := s.http.MakeRequest(deleteRequest)
		if err != nil {
			s.logger.Error("Error in the request of delete application")
			s.logger.Error(err.Error())
			return err
		}
		s.logger.Info("Request deleting send to pod %s", keyName)

	}

	s.logger.Info("The application %s is deleted", url)

	s.DeleteFromApplication(url)
	return nil
}

func (s *Scheduler) DeleteApplicationFromUrl(url string) error {
	application, exists := s.ReadFromApplicationMap(url)
	if !exists {
		return nil
	}
	err := s.deleteApplication(application)
	if err != nil {
		return err
	}
	return nil
}

func (s *Scheduler) addApplication(app *models.SchedulerApplication, startupPhase bool) error {

	hostname := s.ChooseHostname(app)
	app.ScheduledApplication = hostname
	if r, err := regexp.MatchString("-0.|0$", hostname); err == nil && r == true {
		err := s.MonitorApplication(app, startupPhase)
		if err != nil {
			return err
		}
		return nil
	}

	ha := modelsHttp.HttpRequest{
		Name:     hostname,
		Host:     hostname,
		Timeout:  5 * time.Second,
		Protocol: "http",
		Path:     "/api/v1/application/monitor",
		Body:     []models.SchedulerApplication{*app},
		Method:   "PUT",
	}

	for {
		request, err := s.http.MakeRequest(&ha)
		if err != nil && request.StatusCode() != 200 {
			continue
		} else {
			break
		}
	}

	//s.applications[app.Host] = *app
	s.WriteToApplicationMap(app.Host, *app)
	return nil
}

func (s *Scheduler) isEnabled(svc *v1.Service) bool {
	if svc.Annotations[annotations.Enabled] == "true" {
		return true
	}
	return false
}

func (s *Scheduler) getMyHostname() string {
	config, _ := s.config.GetConfig()
	hostname := os.Getenv("HOSTNAME")

	if config.Server.ApplicationName != "localhost" {
		return fmt.Sprintf("%s.%s", hostname, config.Server.ServiceHAName)
	}
	return hostname
}
