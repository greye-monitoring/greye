package application

import (
	"github.com/prometheus/client_golang/prometheus"
	"greye/pkg/metrics/domain/ports"
)

type ClusterMetrics struct {
}

func NewClusterMetric() *ClusterMetrics {
	return &ClusterMetrics{}
}

var _ ports.MetricPorts = (*ClusterMetrics)(nil)

func (m ClusterMetrics) Alarm(label string, value float64) {
	metricsClusterInAlarm.WithLabelValues(label).Set(value)
}

func (m ClusterMetrics) Monitoring(label string, value float64) {
	metricsClusterUnderMonitoring.WithLabelValues(label).Set(value)
}

func (m ClusterMetrics) DeleteMetrics(label string) {
	lbls := prometheus.Labels{"name": label}

	metricsClusterInAlarm.Delete(lbls)
	metricsClusterUnderMonitoring.Delete(lbls)
	metricsClusterUnderMonitoringCounter.Delete(lbls)
	metricsClusterUnderMonitoringLatency.Delete(lbls)
}

func (m ClusterMetrics) DeleteMonitoring(label string) {
	metricsClusterUnderMonitoring.DeleteLabelValues(label)
}

func (m ClusterMetrics) MonitoringCounter(label string, value float64) {
	metricsClusterUnderMonitoringCounter.WithLabelValues(label).Set(value)
}

func (m ClusterMetrics) MonitoringLatency(label string, value float64) {
	metricsClusterUnderMonitoringLatency.WithLabelValues(label).Set(value)
}

var (
	metricsClusterInAlarm = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "cluster_in_alarm",
			Help: "Indicates if an cluster is in alarm (1) or not (0).",
		},
		[]string{"name"},
	)
	metricsClusterUnderMonitoring = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "cluster_under_monitoring",
			Help: "Indicates if an cluster is under monitoring (1) or not (0) or not presente.",
		},
		[]string{"name"},
	)

	metricsClusterUnderMonitoringCounter = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "cluster_counter",
			Help: "Count the number of check for am cluster.",
		},
		[]string{"name"},
	)

	metricsClusterUnderMonitoringLatency = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "cluster_latency",
			Help: "The latency of latest check for an cluster.",
		},
		[]string{"name"},
	)
)

func init() {
	// Register the gauge with Prometheus
	prometheus.MustRegister(metricsClusterInAlarm)
	prometheus.MustRegister(metricsClusterUnderMonitoring)

	prometheus.MustRegister(metricsClusterUnderMonitoringCounter)
	prometheus.MustRegister(metricsClusterUnderMonitoringLatency)
}
