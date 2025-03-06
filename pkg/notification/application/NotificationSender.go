package application

import (
	"fmt"
	"greye/pkg/notification/domain/ports"
)

func NotificationSenderFactory(channel string, config interface{}) (ports.Sender, error) {
	configSender, status := config.(map[string]interface{})
	switch channel {
	case "email":
		sender, err := NewEmailSender(configSender)
		if err != nil || !status {
			panic("error creating email sender")
		}
		return sender, nil
	case "telegram":
		sender, err := NewTelegramSender(configSender)
		if err != nil || !status {
			panic("error creating telegram sender")
		}
		return sender, nil
	default:
		panic(fmt.Sprintf("error creating %s sender", channel))
	}
}
