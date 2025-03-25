package models

type AuthenticationData struct {
	Method   string `json:"method"`
	Username string `json:"username"`
	Password string `json:"-"`
	Url      string `json:"url"`
}
