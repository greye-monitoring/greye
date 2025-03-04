package ports

import (
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/kubernetes"
)

// ImportProcessApplication defines application-level behaviors.
type ImportProcessApplication interface {
	GetKubernetesMonitoringObject(resourceVersion string) watch.Interface
	GetKubernetesServices() *v1.ServiceList
}

// ImportProcessRepository defines repository-level behaviors for Kubernetes.
type ImportProcessRepository interface {
	GetConfig() (*kubernetes.Clientset, error)
}
