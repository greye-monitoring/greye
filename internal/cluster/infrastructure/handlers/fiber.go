package handlers

import (
	"greye/internal/cluster/domain/ports"
	clientHttp "greye/pkg/client/domain/ports"
	logrus "greye/pkg/logging/domain/ports"
	schedulerPort "greye/pkg/scheduler/domain/ports"
	"greye/pkg/server"
)

type ClusterHdl struct {
	cluster     ports.ClusterMonitor
	networkInfo server.NetworkInfo
	logger      logrus.LoggerApplication
	scheduler   schedulerPort.Operation
	http        clientHttp.HttpMethod
}

var _ ports.ApiExposed = (*ClusterHdl)(nil)

func NewClusterHandler(cluster ports.ClusterMonitor, info server.NetworkInfo, logger logrus.LoggerApplication, httpCLient clientHttp.HttpMethod, schedulerHandler schedulerPort.Operation) *ClusterHdl {
	return &ClusterHdl{cluster: cluster, networkInfo: info, logger: logger, http: httpCLient, scheduler: schedulerHandler}
}
