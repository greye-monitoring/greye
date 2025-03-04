package models

type SingleUpdateNode struct {
	Ip                  string `json:"ip"`
	StopMonitoringUntil string `json:"stopMonitoringUntil"`
}
