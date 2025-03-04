package models

type Role string

const (
	Controller Role = "controller"
	Worker     Role = "worker"
)
