package models

import (
	"encoding/json"
	"errors"
	"fmt"
	annotations "greye/pkg/annotations/domain/models"
	models2 "greye/pkg/authentication/domain/models"
	modelsHttp "greye/pkg/client/domain/models"
	"greye/pkg/config/domain/models"
	"greye/pkg/scheduler/application"
	v1 "k8s.io/api/core/v1"
	"strconv"
	"strings"
	"time"
)

type SchedulerApplication struct {
	application.Job
	modelsHttp.MonitoringHttpRequest
	ScheduledApplication    string `json:"scheduledApplication"`
	ForcePodMonitorInstance string `json:"forcePodMonitorInstance"`
	MaxFailRequests         int    `json:"maxFailedRequests"`
	FailedRequest           int    `json:"failedRequest"`
}

func GetPortUsed(svc *v1.Service) int {
	port := int(svc.Spec.Ports[0].Port)
	if portAnnotation := svc.Annotations[annotations.Port]; portAnnotation != "" {
		portInt, err := strconv.Atoi(portAnnotation)
		if err == nil {
			port = int(portInt)
		}
	}
	return port
}

func NewSchedulerApplicationFromService(svc *v1.Service, defaultValue *models.Application) *SchedulerApplication {

	interval := defaultValue.IntervalSeconds * time.Second
	protocol := defaultValue.Protocol
	method := defaultValue.Method
	maxFailedRequests := defaultValue.MaxFailedRequests
	timeout := defaultValue.TimeoutSeconds * time.Second
	monitorInstance := ""

	authMethod := ""
	authUsername := ""
	authPassword := ""

	// If the service has annotations, use them to override the default values
	interval, err := time.ParseDuration(svc.Annotations[annotations.IntervalSeconds] + "s")

	if err != nil {
		interval = defaultValue.IntervalSeconds * time.Second
	}

	if protocolAnnotation := svc.Annotations[annotations.Protocol]; protocolAnnotation != "" {
		protocol = protocolAnnotation
	}

	if maxFailedRequestsAnnotation := svc.Annotations[annotations.MaxFailedRequests]; maxFailedRequestsAnnotation != "" {
		maxFailedRequests, err = strconv.Atoi(maxFailedRequestsAnnotation)
		if err != nil {
			maxFailedRequests = defaultValue.MaxFailedRequests
		}
	}

	if timeoutAnnotation := svc.Annotations[annotations.TimeoutSeconds]; timeoutAnnotation != "" {
		timeout, err = time.ParseDuration(timeoutAnnotation + "s")
		if err != nil {
			timeout = defaultValue.TimeoutSeconds * time.Second
		}
	}

	stopTimeStr := svc.Annotations[annotations.StopMonitoringUntil]
	stopTime := time.Now()
	if stopTimeStr != "" {
		stopTime, err = time.Parse("2006-01-02T15:04:05", stopTimeStr)
		if err == nil {
		}
	}

	if forcePodMonitorInstance := svc.Annotations[annotations.ForcePodMonitor]; forcePodMonitorInstance != "" {
		monitorInstance = forcePodMonitorInstance
	}

	paths := strings.Split(svc.Annotations[annotations.Paths], "\n")

	for _, path := range paths {
		if !strings.HasPrefix(path, "/") {
			path = method + path
		}
	}

	if authMethodAnnotation := svc.Annotations[annotations.AuthenticationMethod]; authMethodAnnotation != "" {
		authMethod = authMethodAnnotation
	}

	if authUsernameAnnotation := svc.Annotations[annotations.AuthenticationUsername]; authUsernameAnnotation != "" {
		authUsername = authUsernameAnnotation
	}

	if authPasswordAnnotation := svc.Annotations[annotations.AuthenticationPassword]; authPasswordAnnotation != "" {
		authPassword = authPasswordAnnotation
	}

	body := svc.Annotations[annotations.Body]
	headerString := svc.Annotations[annotations.Headers]
	var headers = make(map[string]string)
	err = json.Unmarshal([]byte(headerString), &headers)
	if err != nil {
		headers = make(map[string]string)
	}

	port := GetPortUsed(svc)
	svcHost := fmt.Sprint(svc.Name, ".", svc.Namespace, ".svc.cluster.local")
	return &SchedulerApplication{
		Job: application.Job{
			Interval: interval,
			Ticker:   time.NewTicker(interval),
			Quit:     make(chan struct{}),
		},

		MonitoringHttpRequest: modelsHttp.MonitoringHttpRequest{
			Name:                svc.Name,
			Host:                svcHost,
			Timeout:             timeout,
			Protocol:            protocol,
			Path:                paths,
			Interval:            interval,
			Port:                port,
			StopMonitoringUntil: stopTime,
			Body:                body,
			Header:              headers,
			Authentication: models2.AuthenticationData{
				Method:   authMethod,
				Username: authUsername,
				Password: authPassword,
			},
		},
		MaxFailRequests:         maxFailedRequests,
		FailedRequest:           0,
		ForcePodMonitorInstance: monitorInstance,
	}
}

func GenerateJobSchedulerApplication(ri SchedulerApplication) *SchedulerApplication {
	interval := ri.MonitoringHttpRequest.Interval

	ri.Job = application.Job{
		Interval: interval,
		Ticker:   time.NewTicker(interval),
		Quit:     make(chan struct{}),
	}

	return &ri
}

func (s SchedulerApplication) Validate() error {
	if s.MaxFailRequests == 0 {
		return errors.New("MaxFailedRequests must be greater than 0")
	}

	if s.Protocol == "" {
		return errors.New("Protocol must be provided")
	}
	if s.MonitoringHttpRequest.Interval < 30 {
		s.MonitoringHttpRequest.Interval = 30
	}
	if s.Timeout < 1 {
		return errors.New("TimeoutSeconds must be greater or equal to 1 second")
	}
	return nil
}

func (s SchedulerApplication) GetSvcHostname() string {
	return s.Host
}

func (s SchedulerApplication) AddPortToForcePodMonitorInstanceIfMissing() string {
	if s.ForcePodMonitorInstance == "" {
		return ""
	}
	if !strings.Contains(s.ForcePodMonitorInstance, ":") {
		return fmt.Sprintf("%s:%d", s.ForcePodMonitorInstance, s.Port)
	}
	return s.ForcePodMonitorInstance
}
