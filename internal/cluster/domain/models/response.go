package models

type Response struct {
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`  // Holds any data to be returned, omitting it if nil
	Error   string      `json:"error,omitempty"` // Holds an error message if applicable
}
