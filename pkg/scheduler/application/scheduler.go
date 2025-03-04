package application

import (
	"greye/pkg/scheduler/domain/ports"
	"time"
)

type Job struct {
	Interval time.Duration `json:"-"`
	Ticker   *time.Ticker  `json:"-"`
	Quit     chan struct{} `json:"-"`
}

var _ ports.Operation = (*Job)(nil)

func NewJob() *Job {
	return &Job{
		Interval: 60 * time.Second,
		Quit:     make(chan struct{}),
	}
}

func (j Job) Add() {
	//TODO implement me
	panic("implement me")
}

func (j Job) Delete() {
	//TODO implement me
	panic("implement me")
}

func (j Job) Update() {
	//TODO implement me
	panic("implement me")
}

func (j Job) GetById() {
	//TODO implement me
	panic("implement me")
}

func (j Job) GetAll() {
	//TODO implement me
	panic("implement me")
}
