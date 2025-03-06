package application

import (
	b64 "encoding/base64"
	"fmt"
	"greye/pkg/authentication/domain/models"
	"greye/pkg/authentication/domain/ports"
)

type BasicAuth struct {
}

var _ ports.Authentication = (*BasicAuth)(nil)

func newBasicAuth() *BasicAuth {
	return &BasicAuth{}
}

func (b BasicAuth) GetAuthorization(data models.AuthenticationData) (string, error) {
	stringToEncode := fmt.Sprintf("%s:%s", data.Username, data.Password)
	token := b64.StdEncoding.EncodeToString([]byte(stringToEncode))
	result := fmt.Sprintf("Basic %s", token)
	return result, nil
}
