package models

type Annotations string

const (
	Enabled             string = "ge-enabled"
	Paths               string = "ge-paths"
	Body                string = "ge-body"
	Headers             string = "ge-headers"
	Port                string = "ge-port"
	IntervalSeconds     string = "ge-intervalSeconds"
	Protocol            string = "ge-protocol"
	MaxFailedRequests   string = "ge-maxFailedRequests"
	TimeoutSeconds      string = "ge-timeoutSeconds"
	StopMonitoringUntil string = "ge-stopMonitoringUntil"
	ForcePodMonitor     string = "ge-forcePodMonitor"
)
