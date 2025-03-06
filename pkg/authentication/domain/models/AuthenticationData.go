package models

type AuthenticationData struct {
	Method   string `json:"method"`
	Username string `json:"username"`
	Password string
	Url      string `json:"url"`
}
