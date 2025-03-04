package models

// internal/application/domain.models/myRequest.go

type RequestInfo struct {
	Name         string   `json:"name"`
	Namespace    string   `json:"namespace"`
	Host         string   `json:"host"`
	Port         string   `json:"port"`
	Protocol     string   `json:"protocol"`
	Architecture string   `json:"architecture"`
	Path         []string `json:"path"`
}
