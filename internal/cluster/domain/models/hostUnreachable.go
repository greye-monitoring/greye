package models

type HostUnreachable struct {
	Host        []string `json:"host"`
	MaxAttempts int      `json:"max_attempts"`
	Attempts    int      `json:"attempts"`
}
