package application

import (
	"fmt"
	"greye/pkg/authentication/domain/ports"
)

func AuthFactory(authenticationMethod string) ports.Authentication {
	switch authenticationMethod {
	case "basicAuth":
		metrics := newBasicAuth()
		return metrics
	}
	panic(fmt.Sprintf("The authentication method %s is not supported", authenticationMethod))
}
