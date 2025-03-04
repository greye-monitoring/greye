package models

type EasyResponses struct {
	Message string      `json:"message"`
	Error   interface{} `json:"error,omitempty"`
}
