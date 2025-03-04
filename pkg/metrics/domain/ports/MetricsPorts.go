package ports

type MetricPorts interface {
	Alarm(label string, value float64)
	Monitoring(label string, value float64)

	MonitoringCounter(label string, value float64)
	MonitoringLatency(label string, value float64)

	DeleteMetrics(label string)
}
