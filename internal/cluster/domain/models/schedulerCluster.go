package models

import (
	"time"
)

type SchedulerCluster struct {
	Job
	MaxFailedRequest int `json:"max_failed_request"`
}

type Job struct {
	Interval time.Duration
	Ticker   *time.Ticker  `json:"-"`
	Quit     chan struct{} `json:"-"`
}
