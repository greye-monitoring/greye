package models

type ResultSingleRequest struct {
	Response interface{}
	Err      error
}

type Result struct {
	Response []interface{}
	Metrics  MetricsResponse
	Err      []error
	IsOk     bool
}

type MetricsResponse struct {
	Latency float64
}
