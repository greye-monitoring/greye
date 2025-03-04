package models

type ClusterStatus string

const (
	Running   ClusterStatus = "running"
	Deleted   ClusterStatus = "deleted"
	Suspended ClusterStatus = "suspended"
	Error     ClusterStatus = "error"
)
