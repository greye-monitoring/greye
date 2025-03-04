package application

import (
	"errors"
	"github.com/go-resty/resty/v2"
	"greye/pkg/client/domain/models"
	"greye/pkg/client/domain/ports"
	logger "greye/pkg/logging/domain/ports"
	"strings"
	"time"
)

type HttpApplication struct {
	HttpClient *resty.Client
	logger     logger.LoggerApplication
}

var _ ports.HttpMethod = (*HttpApplication)(nil)

func NewHttpApplication(logger logger.LoggerApplication) *HttpApplication {
	restyHttpClient := resty.New().SetRetryCount(3).SetRetryWaitTime(5 * time.Second).AddRetryCondition(func(response *resty.Response, err error) bool {
		return response.StatusCode() >= 500
	})
	return &HttpApplication{HttpClient: restyHttpClient,
		logger: logger}
}

func (h HttpApplication) processPaths(request *models.MonitoringHttpRequest, methodFunc func(string) (*resty.Response, error)) ([]*resty.Response, error) {
	var responses []*resty.Response

	for _, path := range request.Path {
		resp, err := methodFunc(path)
		h.LogResponse(resp, err)
		responses = append(responses, resp)
		if err != nil {
			return responses, err
		}
	}

	return responses, nil
}

// todo ho un problema con il timeout
func (h HttpApplication) GetMonitoring(request *models.MonitoringHttpRequest) ([]*resty.Response, error) {
	return h.processPaths(request, func(path string) (*resty.Response, error) {
		va, err := h.HttpClient.SetTimeout(request.Timeout*time.Second).SetHeader("Content-Type", "application/json").
			R().Get(request.Protocol + "://" + request.Host + path)
		return va, err
	})
}

func (h HttpApplication) Get(request *models.HttpRequest) (*resty.Response, error) {
	return h.HttpClient.SetTimeout(request.Timeout).
		R().Get(request.Protocol + "://" + request.Host + request.Path)
}

func (h HttpApplication) PostMonitoring(request *models.MonitoringHttpRequest) ([]*resty.Response, error) {
	return h.processPaths(request, func(path string) (*resty.Response, error) {
		return h.HttpClient.SetTimeout(request.Timeout).
			R().
			SetBody(request.Body).
			Post(request.Protocol + "://" + request.Host + path)
	})
}

func (h HttpApplication) Post(request *models.HttpRequest) (*resty.Response, error) {
	return h.HttpClient.SetTimeout(request.Timeout).
		R().
		SetBody(request.Body).
		Post(request.Protocol + "://" + request.Host + request.Path)

}

func (h HttpApplication) PutMonitoring(request *models.MonitoringHttpRequest) ([]*resty.Response, error) {
	return h.processPaths(request, func(path string) (*resty.Response, error) {
		return h.HttpClient.SetTimeout(request.Timeout).
			R().
			SetBody(request.Body).
			Put(request.Protocol + "://" + request.Host + path)
	})
}

func (h HttpApplication) DeleteMonitoring(request *models.MonitoringHttpRequest) ([]*resty.Response, error) {
	return h.processPaths(request, func(path string) (*resty.Response, error) {
		va, err := h.HttpClient.SetTimeout(request.Timeout*time.Second).SetHeader("Content-Type", "application/json").
			R().Delete(request.Protocol + "://" + request.Host + path)
		return va, err
	})
}

func (h HttpApplication) Put(request *models.HttpRequest) (*resty.Response, error) {
	return h.HttpClient.SetTimeout(request.Timeout).
		R().
		SetBody(request.Body).
		Put(request.Protocol + "://" + request.Host + request.Path)

}

func (h HttpApplication) Delete(request *models.HttpRequest) (*resty.Response, error) {
	return h.HttpClient.SetTimeout(request.Timeout).
		R().Delete(request.Protocol + "://" + request.Host + request.Path)
}

func (h HttpApplication) LogResponse(resp *resty.Response, err error) {
	if err != nil {
		h.logger.Error("Error:", err)
		return
	}
	h.logger.Info("Response Info:")
	h.logger.Error("  Error      :", err)
	h.logger.Info("  Status Code:", resp.StatusCode())
	h.logger.Info("  Status     :", resp.Status())
	h.logger.Info("  Proto      :", resp.Proto())
	h.logger.Info("  Time       :", resp.Time())
	h.logger.Info("  Received At:", resp.ReceivedAt())
	h.logger.Info("  Body       :\n", resp)

}

func (h HttpApplication) MakeRequest(r *models.HttpRequest) (*resty.Response, error) {
	switch strings.ToUpper(r.Method) {
	case "GET":
		return h.Get(r)
	case "POST":
		return h.Post(r)
	case "PUT":
		return h.Put(r)
	case "DELETE":
		return h.Delete(r)
	}
	return nil, errors.New("Method not implemented")
}
