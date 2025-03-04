package models

type EasyResponse struct {
	Message string      `json:"message"`
	Error   interface{} `json:"error,omitempty"`
}
