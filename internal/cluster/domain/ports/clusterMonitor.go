package ports

import "greye/internal/cluster/domain/models"

type ClusterMonitor interface {
	Status(app models.ClusterInfoResponse) (*models.ClusterInfoResponse, error)
	UpdateSingleNode(app models.SingleUpdateNode) (*models.ClusterInfoDetails, error)
	ReadClustersStatuses() map[string]models.ClusterInfoDetails
	Remove() bool
}
