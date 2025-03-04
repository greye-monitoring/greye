package application

import (
	"github.com/gofiber/fiber/v2/log"
	"github.com/prometheus/client_golang/prometheus"
	"greye/pkg/metrics/domain/ports"
)

type ApplicationMetrics struct {
	Te string
}

func NewApplicationMetric() *ApplicationMetrics {
	return &ApplicationMetrics{}
}

var _ ports.MetricPorts = (*ApplicationMetrics)(nil)

func (m ApplicationMetrics) DeleteMetrics(label string) {
	lbls := prometheus.Labels{"name": label}

	metricsApplicationInAlarm.Delete(lbls)
	metricsApplicationUnderMonitoring.Delete(lbls)
	metricsApplicationUnderMonitoringCounter.Delete(lbls)
	metricsApplicationUnderMonitoringLatency.Delete(lbls)

}

func (m ApplicationMetrics) DeleteMonitoring(label string) {
	lbls := prometheus.Labels{"name": label}
	deleted := metricsApplicationUnderMonitoring.Delete(lbls)
	if deleted {
		log.Infof("Deleted monitoring for label: %s", label)
	} else {
		log.Warnf("No monitoring found for label: %labels", label)
	}
}

func (m ApplicationMetrics) Alarm(label string, value float64) {
	metricsApplicationInAlarm.WithLabelValues(label).Set(value)
}

func (m ApplicationMetrics) Monitoring(label string, value float64) {
	metricsApplicationUnderMonitoring.WithLabelValues(label).Set(value)
}

func (m ApplicationMetrics) MonitoringCounter(label string, value float64) {
	metricsApplicationUnderMonitoringCounter.WithLabelValues(label).Set(value)
}

func (m ApplicationMetrics) MonitoringLatency(label string, value float64) {
	metricsApplicationUnderMonitoringLatency.WithLabelValues(label).Set(value)
}

var (
	metricsApplicationInAlarm = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "application_in_alarm",
			Help: "Indicates if an application is in alarm (1) or not (0).",
		},
		[]string{"name"},
	)

	// Metric to track if an application is under monitoring
	metricsApplicationUnderMonitoring = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "application_under_monitoring",
			Help: "Indicates if an application is under monitoring (1) or not (0) or not presente.",
		},
		[]string{"name"},
	)

	metricsApplicationUnderMonitoringCounter = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "application_counter",
			Help: "Count the number of check for am application.",
		},
		[]string{"name"},
	)

	metricsApplicationUnderMonitoringLatency = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "application_latency",
			Help: "The latency of latest check for an application.",
		},
		[]string{"name"},
	)
)

func init() {
	// Register the gauge with Prometheus
	prometheus.MustRegister(metricsApplicationInAlarm)
	prometheus.MustRegister(metricsApplicationUnderMonitoring)

	prometheus.MustRegister(metricsApplicationUnderMonitoringCounter)
	prometheus.MustRegister(metricsApplicationUnderMonitoringLatency)

}
