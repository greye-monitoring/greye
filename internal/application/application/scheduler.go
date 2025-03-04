package application

import (
	"fmt"
	"greye/internal/application/domain/models"
	"time"
)

func (s *Scheduler) SendNotification(app *models.SchedulerApplication, title string, message string) {

	if app.StopMonitoringUntil.After(time.Now()) {
		s.logger.Warn("The application %s is under allarm but it has been stopped until %s", app.Host, app.StopMonitoringUntil)
		return
	}

	app.FailedRequest = app.FailedRequest + 1
	actualFailedCount := app.FailedRequest
	maxFail := app.MaxFailRequests
	host := app.Host

	s.logger.Error("Error application %s: %s", host, message)
	if actualFailedCount == maxFail {
		s.metrics.Alarm(host, 1)
		alarmSend := false
		for k, alarm := range s.ReadAlarms() {
			//_, err := alarm.Send(title, message)
			//if err != nil {
			//s.logger.Error("Failed to send alarm '%s': %v", alarm, err)
			s.logger.Error("Failed to send alarm '%s': %v", k, alarm)
			//} else {
			//	alarmSend = true
			//}
			if !alarmSend {
				app.FailedRequest = app.FailedRequest - 1
			}
		}
	}
	s.WriteToApplicationMap(host, *app)
}

func (s *Scheduler) MonitorApplication(app *models.SchedulerApplication, startupPhase bool) error {
	svcHostname := app.Host
	application, exist := s.ReadFromApplicationMap(svcHostname)

	if app.Quit == nil {
		app.Quit = make(chan struct{})
	}

	if !startupPhase {
		err := s.deleteApplication(application)
		if err != nil {
			return err
		}
	}
	if !exist {
		application = *app
		s.logger.Error("application not found")
	}

	if app.Quit == nil {
		app.Quit = make(chan struct{})
	}
	s.WriteToApplicationMap(svcHostname, *app)
	s.logger.Info(fmt.Sprintf("Application added: %v\n", svcHostname))
	requestCounter := 0
	s.metrics.MonitoringCounter(svcHostname, float64(requestCounter))
	s.metrics.Monitoring(svcHostname, 1)
	s.metrics.Alarm(svcHostname, 0)
	go func() {
		for {
			c := app.Ticker.C
			q := app.Quit
			select {
			case t := <-c:
				application, _ = s.ReadFromApplicationMap(svcHostname)
				s.logger.Info("Monitoring application %s at %v", svcHostname, t)
				s.metrics.Monitoring(svcHostname, 1)

				if application.Name == "" && application.Host == "" {
					s.logger.Warn("The application %s has been deleted, but another process added it again, deleting...", svcHostname)
					s.deleteApplication(application)
				}

				method, exists := s.ReadFromClient(application.Protocol)
				if !exists {
					title := "Protocol undefined"
					message := fmt.Sprint("Unsupported protocol %s", application.Protocol)
					s.SendNotification(&application, title, message)
					s.logger.Error(message)
					s.WriteToApplicationMap(svcHostname, application)
					break
				}

				res := method.MakeMonitoringRequest(application.MonitoringHttpRequest)
				s.metrics.MonitoringLatency(svcHostname, res.Metrics.Latency)
				requestCounter = requestCounter + 1
				s.metrics.MonitoringCounter(svcHostname, float64(requestCounter))
				s.logger.Debug("The application %s has status %v", svcHostname, res.IsOk)
				if !res.IsOk {
					message := fmt.Sprintf("Application %s is unavailable.", application.Host)
					s.SendNotification(&application, "Application unavailable", message)
					s.logger.Error(message)
					s.WriteToApplicationMap(svcHostname, application)
					_, exist := s.ReadFromApplicationMap(svcHostname)

					if !exist {
						s.metrics.DeleteMetrics(svcHostname)
					}

					break
				}
				if application.FailedRequest != 0 {
					application.FailedRequest = 0
					s.WriteToApplicationMap(svcHostname, application)

				}
				s.metrics.Alarm(svcHostname, 0)
				s.logger.Debug("SUCCESS " + svcHostname)
			case <-q:
				app.Ticker.Stop()
				return
			}
		}
	}()

	return nil
}

// todo teoricamente questo metodo non ha bisogno della lock!
func (s *Scheduler) ChooseHostname(app *models.SchedulerApplication) string {
	svcHostname := app.GetSvcHostname()
	application, _ := s.GetApplication(svcHostname)
	forceScheduledApp := app.ForcePodMonitorInstance
	if len(application) != 0 {
		// application is under monitoring!
		scheduledApplication := application[svcHostname].ScheduledApplication
		ca := application[svcHostname]
		if forceScheduledApp == scheduledApplication || forceScheduledApp == "" {
			app.ScheduledApplication = application[svcHostname].ScheduledApplication
			return application[svcHostname].ScheduledApplication
		}
		s.deleteApplication(ca)
	}

	if forceScheduledApp != "" {
		if val, exist := s.greyesStatus[forceScheduledApp]; exist {
			s.greyesStatus[forceScheduledApp] = val + 1
			return forceScheduledApp
		}
	}

	var minKey = ""
	var minValue int
	for k, v := range s.greyesStatus {
		if minKey == "" || minValue > v {
			minKey = k
			minValue = v
		}
	}
	minValue = minValue + 1
	app.ScheduledApplication = minKey
	s.greyesStatus[minKey] = minValue

	return minKey
}
