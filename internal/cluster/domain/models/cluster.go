package models

import (
	"sync"
	"time"
)

type ClusterInfo struct {
	ClusterInfo sync.Map `json:"cluster-info"`
	Ip          string   `json:"ip"`
}

type ClusterInfoResponse struct {
	ClusterInfo map[string]ClusterInfoDetails `json:"cluster-info"`
	Ip          string                        `json:"ip"`
}

type ClusterInfoDetails struct {
	Status              ClusterStatus `json:"status"`
	Error               ErrorCluster  `json:"error"`
	StopMonitoringUntil string        `json:"stopMonitoringUntil"`
	Timestamp           time.Time     `json:"timestamp"`
}

type ErrorCluster struct {
	ErrorCount int    `json:"error_count"`
	FoundBy    string `json:"found_by"`
	Count      int    `json:"count"`
}

func ConvertClusterInfoToResponse(ci *ClusterInfo) ClusterInfoResponse {
	response := ClusterInfoResponse{
		ClusterInfo: make(map[string]ClusterInfoDetails),
		Ip:          ci.Ip,
	}

	ci.ClusterInfo.Range(func(key, value interface{}) bool {
		if k, ok := key.(string); ok {
			if v, ok := value.(ClusterInfoDetails); ok {
				response.ClusterInfo[k] = v
			}
		}
		return true
	})

	return response
}

func ConvertResponseToClusterInfo(response ClusterInfoResponse) ClusterInfo {
	ci := ClusterInfo{
		ClusterInfo: sync.Map{},
		Ip:          response.Ip,
	}

	for k, v := range response.ClusterInfo {
		ci.ClusterInfo.Store(k, v)
	}

	return ci
}
