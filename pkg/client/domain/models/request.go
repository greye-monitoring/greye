package models

import (
	"net/http"
	"time"
)

type MonitoringHttpRequest struct {
	Name     string        `json:"name"`
	Host     string        `json:"host"`
	Timeout  time.Duration `json:"timeout"`
	Port     int           `json:"port"`
	Protocol string        `json:"protocol"`
	Path     []string      `json:"path"`
	//Method              string        `json:"method"`
	Header              map[string]string `json:"header"`
	Body                interface{}       `json:"body"`
	Interval            time.Duration     `json:"interval"`
	StopMonitoringUntil time.Time         `json:"stopMonitoringUntil"`
}

type HttpRequest struct {
	Name     string        `json:"name"`
	Host     string        `json:"host"`
	Timeout  time.Duration `json:"timeout"`
	Protocol string        `json:"protocol"`
	Path     string        `json:"path"`
	Method   string        `json:"method"`
	Header   http.Header   `json:"header"`
	Body     interface{}   `json:"body"`
}
