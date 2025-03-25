package application

import (
	"errors"
	"github.com/go-resty/resty/v2"
	"greye/pkg/client/domain/models"
	"greye/pkg/client/domain/ports"
	logger "greye/pkg/logging/domain/ports"
	"strconv"
	"strings"
	"sync"
)

type HttpMonitoring struct {
	HttpClient *resty.Client
	logger     logger.LoggerApplication
}

func NewHttpMonitoring(logger logger.LoggerApplication) *HttpMonitoring {
	restyHttpClient := resty.New()
	return &HttpMonitoring{HttpClient: restyHttpClient,
		logger: logger}
}

var _ ports.MonitoringMethod = (*HttpMonitoring)(nil)

func (h HttpMonitoring) GetMonitoring(request models.MonitoringHttpRequest, path string) (interface{}, error) {
	va, err := h.HttpClient.SetTimeout(request.Timeout).SetHeaders(request.Header).
		R().Get(request.Protocol + "://" + request.Host + ":" + strconv.Itoa(request.Port) + path)
	return va, err

}

func (h HttpMonitoring) PostMonitoring(request models.MonitoringHttpRequest, path string) (interface{}, error) {
	return h.HttpClient.SetTimeout(request.Timeout).SetHeaders(request.Header).
		R().
		SetBody(request.Body).
		Post(request.Protocol + "://" + request.Host + ":" + strconv.Itoa(request.Port) + path)
}

func (h HttpMonitoring) PutMonitoring(request models.MonitoringHttpRequest, path string) (interface{}, error) {
	return h.HttpClient.SetTimeout(request.Timeout).SetHeaders(request.Header).
		R().
		SetBody(request.Body).
		Put(request.Protocol + "://" + request.Host + ":" + strconv.Itoa(request.Port) + path)

}

func (h HttpMonitoring) DeleteMonitoring(request models.MonitoringHttpRequest, path string) (interface{}, error) {
	va, err := h.HttpClient.SetTimeout(request.Timeout).SetHeaders(request.Header).
		R().Delete(request.Protocol + "://" + request.Host + ":" + strconv.Itoa(request.Port) + path)
	return va, err

}

func (h HttpMonitoring) LogResponse(resp interface{}, err error) {
	if err != nil {
		h.logger.Error("Error: %s", err.Error())
		return
	}

}

func (h HttpMonitoring) MakeMonitoringRequest(r models.MonitoringHttpRequest) models.Result {

	results := make(chan models.ResultSingleRequest, len(r.Path))
	var wg sync.WaitGroup

	for _, paths := range r.Path {
		wg.Add(1)
		go func(path string) {
			defer wg.Done()
			// if paths start with '/' use GET method, else retrieve everything from start until first / let everything else be the method
			method := "GET"

			if !strings.HasPrefix(path, "/") && path != "" {
				parts := strings.SplitN(path, "/", 2)
				method = strings.ToUpper(parts[0])
				if len(parts) == 2 {
					path = "/" + parts[1]
				} else {
					path = ""
				}
			}
			var res interface{}
			var err error
			switch method {
			case "GET":
				res, err = h.GetMonitoring(r, path)
			case "POST":
				res, err = h.PostMonitoring(r, path)
			case "PUT":
				res, err = h.PutMonitoring(r, path)
			case "DELETE":
				res, err = h.DeleteMonitoring(r, path)
			default:
				err = errors.New("Method not implemented")
			}

			results <- models.ResultSingleRequest{Response: res, Err: err}

		}(paths)
	}

	wg.Wait()
	close(results)
	var responses models.Result
	var latency float64
	latency = 0
	for res := range results {
		responses.Response = append(responses.Response, res.Response)
		responses.Err = append(responses.Err, res.Err)
		if response, ok := res.Response.(*resty.Response); ok {
			latency += response.Time().Seconds()
		}
	}
	latency = latency / float64(len(responses.Response))
	responses.Metrics.Latency = latency
	responses.IsOk = h.checkResponse(responses)
	return responses
}

func (h HttpMonitoring) processPaths(request models.MonitoringHttpRequest, methodFunc func(string) (interface{}, error)) ([]interface{}, error) {
	var responses []interface{}

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

func (h HttpMonitoring) checkResponse(responses models.Result) bool {
	isOk := true

	if len(responses.Err) > 0 && responses.Err[0] != nil {
		isOk = false
		return isOk
	}

	for _, res := range responses.Response {
		response, ok := res.(*resty.Response)
		if !ok {
			isOk = false
			return isOk
		}
		if response.StatusCode() < 200 || response.StatusCode() > 300 {
			isOk = false
			return isOk
		}
	}

	return isOk
}
