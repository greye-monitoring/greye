package ports

import (
	"github.com/go-resty/resty/v2"
	"greye/pkg/client/domain/models"
)

type HttpMethod interface {
	//CreateRequest(name string, url string, path []string, duration time.Duration, protocol string) *models.HttpRequest
	//MakeMonitoringRequest(request *models.MonitoringHttpRequest) ([]*resty.Response, error)
	MakeRequest(request *models.HttpRequest) (*resty.Response, error)
}
