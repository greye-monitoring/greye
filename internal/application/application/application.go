package application

import (
	"encoding/json"
	"fmt"
	"greye/internal/application/domain/models"
	"greye/internal/application/domain/ports"
	portsAuth "greye/pkg/authentication/domain/ports"
	modelsHttp "greye/pkg/client/domain/models"
	clientPort "greye/pkg/client/domain/ports"
	configPort "greye/pkg/config/domain/ports"
	k8s "greye/pkg/importProcess/application"
	importProcess "greye/pkg/importProcess/domain/ports"
	logger "greye/pkg/logging/domain/ports"
	metricsPort "greye/pkg/metrics/domain/ports"
	ports2 "greye/pkg/notification/domain/ports"
	roleModel "greye/pkg/role/domain/models"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/watch"
	"log"
	"regexp"
	"sync"
	"time"
)

type Scheduler struct {
	applications sync.Map
	greyesStatus map[string]int

	k8s            importProcess.ImportProcessApplication
	http           clientPort.HttpMethod
	config         configPort.ConfigApplication
	role           roleModel.Role
	logger         logger.LoggerApplication
	alarms         map[string]ports2.Sender
	client         map[string]clientPort.MonitoringMethod
	metrics        metricsPort.MetricPorts
	authentication map[string]portsAuth.Authentication
	sync.RWMutex
}

var (
	_ ports.SchedulerService = (*Scheduler)(nil)
)

func GetApplicationInitialized(host string, http clientPort.HttpMethod, appInitialized *map[string]models.SchedulerApplication) int {
	//todo forse questo metodo va messo da un'altra parte!
	var cmstatus int
	for {
		httpRequest := modelsHttp.HttpRequest{
			Name:     host,
			Host:     host,
			Timeout:  5 * time.Second,
			Protocol: "http",
			Path:     "/api/v1/application/monitor",
			Method:   "GET",
		}

		request, err := http.MakeRequest(&httpRequest)
		if err != nil {
			log.Printf("Request failed: %v. Retrying...", err)
			time.Sleep(5 * time.Second) // Small delay before retry
			continue
		}

		if request.StatusCode() == 200 {
			cmstatus = 0
			log.Println("Received 200 response. Processing the body...")

			// Assume the body contains a map (e.g., JSON parsed as map[string]models.SchedulerApplication)
			var responseBody map[string]models.SchedulerApplication
			err = json.Unmarshal(request.Body(), &responseBody) // Ensure request.Body() is []byte
			if err != nil {
				log.Printf("Failed to parse response body: %v. Retrying...", err)
				time.Sleep(5 * time.Second)
				continue
			}

			// Add all key-value pairs to appInitialized
			for key, value := range responseBody {
				if _, exists := (*appInitialized)[key]; exists {
					panic(fmt.Sprintf("App %s already initialized", key))
				}
				cmstatus++
				(*appInitialized)[key] = value
			}

			log.Println("Received 200 response. Exiting loop.")
			break
		}

		log.Printf("Non-200 response: %d. Retrying...", request.StatusCode())
		time.Sleep(5 * time.Second)
	}
	return cmstatus
}

func (s *Scheduler) ManageStartupWorker() {

	config, err := s.config.GetConfig()
	if err != nil {
		return
	}
	appName := config.Server.ApplicationName
	var hostController string
	if appName == "localhost" {
		hostController = fmt.Sprintf("%s:8080", appName)
	} else {
		hostController = fmt.Sprintf("%s-0.%s", appName, config.Server.ServiceHAName)
	}

	s.logger.Error("hostController " + hostController)
	httpRequest := modelsHttp.HttpRequest{
		Name:     hostController,
		Host:     hostController,
		Timeout:  5 * time.Second,
		Protocol: "http",
		Path:     "/api/v1/application/monitor",
		Method:   "GET",
	}

	request, err := s.http.MakeRequest(&httpRequest)

	if request.StatusCode() == 200 {
		s.logger.Error("Received 200 response. Processing the body...")

		// Assume the body contains a map (e.g., JSON parsed as map[string]models.SchedulerApplication)
		var responseBody map[string]models.SchedulerApplication
		err = json.Unmarshal(request.Body(), &responseBody) // Ensure request.Body() is []byte
		if err != nil {
			s.logger.Error("Failed to parse response body: %v. Retrying...", err)
			time.Sleep(5 * time.Second)
			return
		}
		hostname := s.getMyHostname()
		// Add all key-value pairs to appInitialized
		for _, value := range responseBody {
			if value.ScheduledApplication == hostname {

				data := models.GenerateJobSchedulerApplication(value)

				s.MonitorApplication(data, true)
			}
		}

		s.logger.Error("Received 200 response. Exiting loop.")
	}
}

func (s *Scheduler) RemoveNoMoreUsedSvcFoundStartupPhase(svcList *v1.ServiceList, monitoredAppFromOtherPod *map[string]models.SchedulerApplication) {
	//per ogni monitoredAppFromOtherPod
	var svcMap = make(map[string]v1.Service)

	for _, svc := range svcList.Items {
		svcHost := fmt.Sprintf("%s.%s.svc.cluster.local", svc.ObjectMeta.Name, svc.ObjectMeta.Namespace)
		svcMap[svcHost] = svc
	}

	for _, app := range *monitoredAppFromOtherPod {
		if _, exists := svcMap[app.Host]; !exists { // If the service no longer exists
			s.logger.Error("Service %s no longer exists, deleting monitoring", app.Host)
			s.deleteApplication(app) // Remove monitoring for that app
		}
	}
}

func (s *Scheduler) ManageStartupController(monitoredAppFromOtherPod *map[string]models.SchedulerApplication) {
	svcList := s.k8s.GetKubernetesServices()
	s.RemoveNoMoreUsedSvcFoundStartupPhase(svcList, monitoredAppFromOtherPod)

	//var svcToMonitor = make(map[string]v1.Service)
	var bulkMonitor = make(map[string][]*models.SchedulerApplication)
	config, _ := s.config.GetConfig()

	nServicesAtStartTime := len(svcList.Items)
	servicesElaborated := 0
	go func() {
		resourceVersion := ""
		for {
			svcWatch := s.k8s.GetKubernetesMonitoringObject(resourceVersion)
			ch := svcWatch.ResultChan()

			applicationController := make(map[string]models.SchedulerApplication)

			for event := range ch {

				servicesElaborated++

				svc, ok := event.Object.(*v1.Service)
				if !ok {
					if status, isStatus := event.Object.(*metav1.Status); isStatus && status.Reason == metav1.StatusReasonExpired {
						// Resource version expired -> Get new latest version
						s.logger.Error("Resource version expired, restarting watch with latest version...")

						svcList := s.k8s.GetKubernetesServices()
						if svcList == nil || len(svcList.Items) == 0 {
							s.logger.Error("Failed to retrieve services, retrying...")
							time.Sleep(5 * time.Second) // Avoid tight loop
							break
						}

						// Restart with the latest resourceVersion
						resourceVersion = svcList.ResourceVersion
						break
					}

					s.logger.Error("Received an object that is not a Service: %T", event.Object)
					continue
				}
				resourceVersion = svc.ResourceVersion

				metadata := svc.ObjectMeta
				s.logger.Error("Service %s/%s received", svc.ObjectMeta.Namespace, svc.ObjectMeta.Name)
				isEnabled := s.isEnabled(svc)
				if isEnabled {
					s.logger.Error("Adding service %s ...", svc.ObjectMeta.Name)
				}

				host := fmt.Sprintf("%s.%s.svc.cluster.local", metadata.Name, metadata.Namespace)
				appExist := false

				//usedPort := models.GetPortUsed(svc)

				func() {

					if _, exists := (applicationController)[host]; exists {
						appExist = true
					}
					s.logger.Error("the event typer is %s", event.Type)
					if !isEnabled && !appExist && event.Type != watch.Deleted {
						s.logger.Error("The application %s is not enabled and is not under monitoring. Skipping...", host)
						return
					}

					if event.Type == watch.Deleted || (!isEnabled && appExist) {
						s.logger.Error("The application %s is not enabled and is under monitoring. Deleting...", host)
						ReadFromApplicationMap, _ := s.ReadFromApplicationMap(host)
						err := s.deleteApplication(ReadFromApplicationMap)
						if err != nil {
							s.logger.Error("Error during deleting the application %s.", host)
							s.logger.Error(err.Error())
							return
						}
						delete(applicationController, host)
						//time.Sleep(1 * time.Millisecond)
						return
					}

					defaultValue := config.Application
					appModels := models.NewSchedulerApplicationFromService(svc, &defaultValue)

					if isEnabled {
						s.logger.Error("The application %s requested to be monitored. Adding or updating...", host)

						err := appModels.Validate()
						if err != nil {
							s.logger.Error("Error validating application: %v", err)
						}
						applicationController[host] = *appModels

						if servicesElaborated > nServicesAtStartTime {
							s.addApplication(appModels, false)
						} else {

							hostname := ""
							alreadySchedueldApplication := (*monitoredAppFromOtherPod)[appModels.Host].ScheduledApplication

							if alreadySchedueldApplication != "" {
								hostname = alreadySchedueldApplication
							} else {
								hostname = s.ChooseHostname(appModels)
							}
							appModels.ScheduledApplication = hostname
							bulkMonitor[hostname] = append(bulkMonitor[hostname], appModels)

							s.WriteToApplicationMap(appModels.Host, *appModels)
						}
					}
				}()

				appName := config.Server.ApplicationName
				if servicesElaborated == nServicesAtStartTime {
					s.logger.Error("All services at startup have been elaborated, executing bulk requests")
					for bulkHostname, bulkApps := range bulkMonitor {

						s.logger.Error("Starting bulk monitoring for %s", bulkHostname)
						regexPattern := fmt.Sprintf("^%s-0.|localhost:[0-9]*0$", appName)
						if r, err := regexp.MatchString(regexPattern, bulkHostname); err == nil && r == true {

							for _, app := range bulkApps {
								err := s.MonitorApplication(app, true)
								if err != nil {
									s.logger.Error("Error monitoring %s: %s", bulkHostname, err.Error())
								}
							}
						} else {
							ha := modelsHttp.HttpRequest{
								Name:     bulkHostname,
								Host:     bulkHostname,
								Timeout:  60 * time.Second,
								Protocol: "http",
								Path:     "/api/v1/application/monitor",
								Body:     bulkApps,
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
						}
						s.logger.Error("Bulk monitoring for %s completed", bulkHostname)
					}
					s.logger.Error("All bulk monitoring requests completed")
				}
			}
			s.logger.Error("ERROR WITH CHANNEL")
		}
	}()
}

func (s *Scheduler) ManageStartup(monitoredAppFromOtherPod *map[string]models.SchedulerApplication) error {
	if s.role == roleModel.Worker {
		s.ManageStartupWorker()
	} else {
		s.ManageStartupController(monitoredAppFromOtherPod)
	}
	return nil
}

func NewScheduler(http clientPort.HttpMethod, config configPort.ConfigApplication, roleType roleModel.Role, logger logger.LoggerApplication, notification map[string]ports2.Sender, client map[string]clientPort.MonitoringMethod, metrics *metricsPort.MetricPorts, importData *k8s.ImportProcessApplication, auth map[string]portsAuth.Authentication) *Scheduler {
	c, err := config.GetConfig()
	if err != nil {
		return nil
	}
	cmstatus := make(map[string]int)
	nClusterMonitor := c.Server.NumberGreye
	appName := c.Server.ApplicationName
	svcHAName := c.Server.ServiceHAName
	port := c.Server.Port
	var monitoredAppFromOtherPod = &map[string]models.SchedulerApplication{}
	if roleType == roleModel.Controller {
		for i := 0; i < int(nClusterMonitor); i++ {
			var k string
			if appName == "localhost" {
				k = fmt.Sprintf("%s:808%d", appName, i)
			} else {
				k = fmt.Sprintf("%s-%d.%s:%d", appName, i, svcHAName, port)
			}
			cmstatus[k] = 0
			if i != 0 {
				cmstatus[k] = GetApplicationInitialized(k, http, monitoredAppFromOtherPod)
			}
		}
	}
	//
	s := &Scheduler{applications: sync.Map{},
		http:           http,
		config:         config,
		role:           roleType,
		greyesStatus:   cmstatus,
		k8s:            importData,
		logger:         logger,
		alarms:         notification,
		client:         client,
		metrics:        *metrics,
		authentication: auth,
	}

	s.ManageStartup(monitoredAppFromOtherPod)

	return s
}
