package application

import (
	"encoding/json"
	"errors"
	"greye/internal/cluster/domain/models"
	"greye/internal/cluster/domain/ports"
	_ "greye/pkg/annotations/domain/models"
	models2 "greye/pkg/client/domain/models"
	clientApp "greye/pkg/client/domain/ports"
	configPort "greye/pkg/config/domain/ports"
	logger "greye/pkg/logging/domain/ports"
	metricsPort "greye/pkg/metrics/domain/ports"
	ports2 "greye/pkg/notification/domain/ports"
	"greye/pkg/scheduler/application"
	"math"
	"math/rand/v2"
	"reflect"
	"slices"
	"sync"
	"time"
)

type Cluster struct {
	cluster models.ClusterInfo

	hostUnreachable models.HostUnreachable

	RequestCounter sync.Map
	application.Job
	http    clientApp.HttpMethod
	config  configPort.ConfigApplication
	logger  logger.LoggerApplication
	alarms  map[string]ports2.Sender
	metrics metricsPort.MetricPorts
}

var _ ports.ClusterMonitor = (*Cluster)(nil)

func (c *Cluster) ReadClustersStatuses() map[string]models.ClusterInfoDetails {

	copiedMap := make(map[string]models.ClusterInfoDetails)
	c.cluster.ClusterInfo.Range(func(key, value interface{}) bool {
		if k, ok := key.(string); ok {
			if v, ok := value.(models.ClusterInfoDetails); ok {
				copiedMap[k] = v
			}
		}
		return true
	})

	return copiedMap
}

func (c *Cluster) ReadAlarms() map[string]ports2.Sender {
	return c.alarms
}

func NewCluster(http clientApp.HttpMethod, config configPort.ConfigApplication,
	logger logger.LoggerApplication, notification map[string]ports2.Sender,
	metrics *metricsPort.MetricPorts) *Cluster {
	getConfig, _ := config.GetConfig()
	interval := time.Duration(getConfig.Cluster.IntervalSeconds) * time.Second
	m := models.ClusterInfoDetails{
		Status:    models.Running,
		Timestamp: time.Now().Local(),
	}
	c := &Cluster{
		cluster: models.ClusterInfo{
			ClusterInfo: sync.Map{},
			Ip:          getConfig.Cluster.MyIp,
		},
		Job: application.Job{
			Interval: interval,
			Ticker:   time.NewTicker(interval),
			Quit:     make(chan struct{}),
		},

		hostUnreachable: models.HostUnreachable{
			Host:        getConfig.Cluster.ClusterIp,
			MaxAttempts: 1,
			Attempts:    0,
		},

		http:    http,
		config:  config,
		logger:  logger,
		alarms:  notification,
		metrics: *metrics,
	}

	c.cluster.ClusterInfo.Store(getConfig.Cluster.MyIp, m)

	// Initialize the cluster information
	for hostUnreached := range getConfig.Cluster.ClusterIp {
		c.metrics.Monitoring(getConfig.Cluster.ClusterIp[hostUnreached], 0)
	}
	c.monitorCluster()
	return c

}

func (c *Cluster) multipleRequest(ip []string) ([]string, []string) {

	var hostReached []string
	var hostUnreached []string

	for _, cluster := range ip {
		data, err := c.execRequest(cluster)
		if err != nil {
			c.logger.Error("Error getting cluster information from ", cluster, err)
			hostUnreached = append(hostUnreached, cluster)
			continue
		}
		hostReached = append(hostReached, cluster)
		_, _ = c.verifyAndUpdate(*data)
	}
	return hostReached, hostUnreached
}

func (c *Cluster) ReadApplications(filter models.ClusterStatus, host string, excludeMe bool) map[string]models.ClusterInfoDetails {

	copiedMap := make(map[string]models.ClusterInfoDetails)

	c.cluster.ClusterInfo.Range(func(key, value interface{}) bool {
		if k, ok := key.(string); ok {
			if v, ok := value.(models.ClusterInfoDetails); ok {
				if host != "" && host != k {
					return true
				}
				if excludeMe && k == c.cluster.Ip {
					return true
				}
				if filter == "" {
					copiedMap[k] = v
				}
				if filter != "" && v.Status == filter {
					copiedMap[k] = v
				}
			}
		}
		return true
	})

	return copiedMap
}

func (c *Cluster) calculateNumberOfCalls(nClusters int, pow float64) int {

	res := math.Pow(float64(nClusters), pow)
	return int(math.Ceil(res))
}

func (c *Cluster) countNumberOfElements() int {
	count := 0
	c.cluster.ClusterInfo.Range(func(key, value interface{}) bool {
		count++
		return true
	})
	return count
}

func (c *Cluster) monitorCluster() {
	go func() {
		for {
			ct := c.Ticker.C
			select {
			case <-ct:
				c.checkUnreachableHosts()
				c.logger.Info("Checking cluster information")
				c.rescheduleCluster()
				hostToCall := c.getHostToUpdate()
				c.multipleRequest(hostToCall)
				c.logger.Info("Number of elements in the map: %d", c.countNumberOfElements())
			}
		}
	}()
}
func (c *Cluster) rescheduleCluster() {
	applications := c.ReadApplications(models.Suspended, "", false)
	for k, v := range applications {
		stopMonitoringClusterUntil, err := time.Parse("2006-01-02T15:04:05", v.StopMonitoringUntil)
		if err != nil {
			c.logger.Error("Error parsing time: ", err)
			continue
		}
		if stopMonitoringClusterUntil.Before(time.Now()) {
			v.Status = models.Running
			//v.StopMonitoringUntil = ""
			c.cluster.ClusterInfo.Store(k, v)
		}
	}

}

func (c *Cluster) getHostToUpdate() []string {
	config, _ := c.config.GetConfig()

	numberOfExecutions := config.Cluster.MaxFailedRequests

	pow := float64(1 / float64(numberOfExecutions))

	nClusters := c.countNumberOfElements()
	nCall := c.calculateNumberOfCalls(nClusters, pow)

	ErrorClustermap := c.ReadApplications(models.Error, "", true)
	clusterMap := c.ReadApplications(models.Running, "", true)

	var hostAvailable []string
	for k := range clusterMap {
		hostAvailable = append(hostAvailable, k)
	}

	var hostToCall []string

	for k, v := range ErrorClustermap {
		if v.Error.FoundBy == c.cluster.Ip {
			hostToCall = append(hostToCall, k)
		}
	}

	for i := 0; i < nCall; i++ {
		if len(hostAvailable) > 0 {
			randomNumber := rand.IntN(len(hostAvailable))
			hostToCall = append(hostToCall, hostAvailable[randomNumber])
			hostAvailable = append(hostAvailable[:randomNumber], hostAvailable[randomNumber+1:]...)
		}
	}
	return hostToCall
}

func (c *Cluster) sendNotification(host string) {
	errorCluster := c.ReadApplications("", host, false)
	ec := errorCluster[host]
	if ec.Status != models.Error {
		ec.Status = models.Error
		ec.Error = models.ErrorCluster{
			FoundBy: c.cluster.Ip,
			Count:   1,
		}
		ec.Timestamp = time.Now().Local()
		c.cluster.ClusterInfo.Store(host, ec)
	} else {
		if ec.Error.FoundBy == c.cluster.Ip {
			ec.Error.Count++
			c.cluster.ClusterInfo.Store(host, ec)
			config, _ := c.config.GetConfig()
			if ec.Error.Count == config.Cluster.MaxFailedRequests {
				c.logger.Error("Sending notification to ", host)
				c.metrics.Alarm(host, 1)
				alarmSend := false

				for k, alarm := range c.ReadAlarms() {
					//_, err := alarm.Send(title, message)
					//if err != nil {
					//s.logger.Error("Failed to send alarm '%s': %v", alarm, err)
					c.logger.Error("Failed to send alarm '%s': %v", k, alarm)
					//} else {
					//	alarmSend = true
					//}
				}

				if alarmSend {
					ec.Error.Count--
					c.cluster.ClusterInfo.Store(host, ec)
				}
			}
		}
	}

}

// function exec request to other cluster, passing the host and all clusterInfo Map
func (c *Cluster) execRequest(host string) (*models.ClusterInfoResponse, error) {
	c.logger.Info("Executing request to %s", host)

	toResponse := models.ConvertClusterInfoToResponse(&c.cluster)

	request := &models2.HttpRequest{
		Host:     host,
		Protocol: "http",
		Path:     "/api/v1/cluster/status",
		Method:   "PUT",
		Body:     toResponse,
	}

	clusterInfo, err := c.http.MakeRequest(request)
	value, ok := c.RequestCounter.Load(host)
	if !ok {
		c.metrics.MonitoringCounter(host, 1)
		c.RequestCounter.Store(host, float64(1))
	} else {
		counter := value.(float64) + 1
		c.metrics.MonitoringCounter(host, counter)
		c.RequestCounter.Store(host, counter)
	}
	if err != nil {
		c.logger.Error("Error getting cluster information from ", host, err)
		c.sendNotification(host)
		return nil, err
	}

	c.metrics.MonitoringLatency(host, clusterInfo.Time().Seconds())
	c.metrics.Alarm(host, 0)

	var response models.ClusterInfoResponse
	err = json.Unmarshal(clusterInfo.Body(), &response) // Convert response body to struct
	if err != nil {
		c.logger.Error("Error unmarshaling response body: ", err)
		return nil, err
	}

	return &response, nil
}

func (c *Cluster) updateClusterInfo(m models.ClusterInfoDetails, y models.ClusterInfoDetails, host string, myIp string, yourIp string) {
	if m.Status == models.Error && y.Status == models.Error {
		c.metrics.Alarm(host, 1)
		switch {
		case m.Timestamp.After(y.Timestamp):
			c.logger.Info("Updating key: ", host)
			c.cluster.ClusterInfo.Store(host, y)
			break
		case m.Timestamp.Equal(y.Timestamp):
			if m.Error.FoundBy == myIp {
				break
			}

			if m.Error.FoundBy == yourIp {
				c.cluster.ClusterInfo.Store(host, y)
				break
			}

			maxCount := math.Max(float64(m.Error.Count), float64(y.Error.Count))

			m.Error.Count = int(maxCount)
			c.cluster.ClusterInfo.Store(host, m)

			break
		case m.Timestamp.Before(y.Timestamp):
			if y.Error.FoundBy != yourIp {
				c.execRequest(y.Error.FoundBy)
			}
			break
		}
	}

	if m.Status == models.Running && y.Status == models.Suspended && m.StopMonitoringUntil == "" {

		c.logger.Info("Updating key: ", host)
		c.metrics.Monitoring(host, 0)
		c.metrics.Alarm(host, 0)
		c.cluster.ClusterInfo.Store(host, y)
		return
	}

	if m.Status == models.Error && (y.Status == models.Deleted || y.Status == models.Suspended) {
		c.metrics.Monitoring(host, 0)
		c.metrics.Alarm(host, 0)
		c.cluster.ClusterInfo.Store(host, y)
		return
	}

	// Default case: update the key if the timestamp is greater than the current one
	if m.Timestamp.Before(y.Timestamp) {
		c.logger.Info("Updating key: ", host)

		if m.Status == models.Deleted && y.Status == models.Error {
			c.metrics.Monitoring(host, 0)
			c.metrics.Alarm(host, 0)
			return
		}

		if y.Status == models.Running && m.Status == models.Suspended && y.StopMonitoringUntil == "" {
			return
		}

		c.cluster.ClusterInfo.Store(host, y)

		if y.Status == models.Deleted {
			c.metrics.Monitoring(host, 0)
		}

		if y.Status == models.Error {
			c.metrics.Alarm(host, 1)
		} else {
			c.metrics.Alarm(host, 0)
		}
		c.metrics.Monitoring(host, 1)
	}
}

func (c *Cluster) verifyAndUpdate(ci models.ClusterInfoResponse) (*models.ClusterInfoResponse, error) {

	response := models.ConvertClusterInfoToResponse(&c.cluster)
	areEquals := reflect.DeepEqual(response.ClusterInfo, ci.ClusterInfo)
	if areEquals {
		c.logger.Info("Cluster information are equals")
		return &response, nil
	}

	c.logger.Info("Cluster information aren't the same, checking for differences")

	c.cluster.ClusterInfo.Range(func(key, value interface{}) bool {
		if k, ok := key.(string); ok {
			if v, ok := value.(models.ClusterInfoDetails); ok {

				ciLoadedMap, ok := ci.ClusterInfo[k]
				if !ok {
					c.logger.Info("Key not found in this cluster information: ", k)
					return true
				}

				areEquals = reflect.DeepEqual(v, ciLoadedMap)
				if !areEquals {
					c.logger.Info("Got differences in values for key: ", k)
					c.updateClusterInfo(v, ci.ClusterInfo[k], k, c.cluster.Ip, ci.Ip)
				}
				delete(ci.ClusterInfo, k)
			}
		}
		return true
	})

	for k, v := range ci.ClusterInfo {
		c.logger.Info("Adding new key to the cluster information: ", k)
		c.cluster.ClusterInfo.Store(k, v)
		c.metrics.Monitoring(k, 1)
	}

	response = models.ConvertClusterInfoToResponse(&c.cluster)

	return &response, nil
}

func (c *Cluster) Status(ci models.ClusterInfoResponse) (*models.ClusterInfoResponse, error) {

	updatedClusterInfo, err := c.verifyAndUpdate(ci)
	if err != nil {
		return nil, err
	}
	return updatedClusterInfo, nil
}

func (c *Cluster) UpdateSingleNode(updateValue models.SingleUpdateNode) (*models.ClusterInfoDetails, error) {

	applications := c.ReadApplications("", updateValue.Ip, false)
	if _, ok := applications[updateValue.Ip]; !ok {
		return nil, errors.New("Cluster not found")
	}
	details := applications[updateValue.Ip]

	stopTime := updateValue.StopMonitoringUntil
	if updateValue.StopMonitoringUntil != "" {
		_, err := time.Parse("2006-01-02T15:04:05", stopTime)
		if err != nil {
			return nil, err
		}
		//if stopMonitoringClusterUntil.After(time.Now()) {
		details.StopMonitoringUntil = stopTime
		details.Status = models.Suspended
		details.Timestamp = time.Now()
		c.cluster.ClusterInfo.Store(updateValue.Ip, details)
		//}
	}

	return &details, nil
}

func (c *Cluster) checkUnreachableHosts() {
	if len(c.hostUnreachable.Host) == 0 {
		return
	}

	c.hostUnreachable.Attempts = c.hostUnreachable.Attempts + 1
	if c.hostUnreachable.Attempts == c.hostUnreachable.MaxAttempts {
		c.hostUnreachable.Attempts = 0
		hostReached, _ := c.multipleRequest(c.hostUnreachable.Host)
		for host := range hostReached {
			c.logger.Info("Host %s is reachable", hostReached[host])
			c.metrics.Monitoring(hostReached[host], 1)
			c.hostUnreachable.Host = slices.DeleteFunc(c.hostUnreachable.Host, func(s string) bool {
				return s == hostReached[host]
			})
			c.logger.Info("Host %s is reachable", hostReached[host])
		}
	}

}

func (c *Cluster) Remove() bool {
	for k, v := range c.ReadApplications(models.Error, "", true) {
		if v.Status == models.Error {
			c.logger.Error("Cluster %s is in error status", k)
			return false
		}
	}

	c.Ticker.Stop()
	myIp := c.cluster.Ip
	details, ok := c.cluster.ClusterInfo.Load(myIp)
	if ok {
		detail := details.(models.ClusterInfoDetails)
		detail.Status = models.Deleted
		detail.Timestamp = time.Now().Local()
		c.cluster.ClusterInfo.Store(myIp, detail)
	}

	hostToCall := c.getHostToUpdate()
	hostReceivedRequest, _ := c.multipleRequest(hostToCall)

	if len(hostReceivedRequest) > 0 {
		return true
	} else {
		return false
	}
}
