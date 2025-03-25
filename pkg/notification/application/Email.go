package application

import (
	"errors"
	"fmt"
)

// TelegramSender è una struttura che implementa l'interfaccia Sender per inviare notifiche su Telegram.

type EmailSender struct {
	Token       string
	Destination string
}

func NewEmailSender(config map[string]interface{}) (*EmailSender, error) {
	// Verifica e converte il campo "token"
	token, ok := config["token"].(string)
	if !ok || token == "" {
		return nil, errors.New("invalid or missing 'token' in TelegramSender config")
	}

	// Verifica e converte il campo "destination"
	destination, ok := config["destination"].(string)
	if !ok || destination == "" {
		return nil, errors.New("invalid or missing 'destination' in TelegramSender config")
	}

	// Crea l'istanza di TelegramSender
	return &EmailSender{
		Token:       token,
		Destination: destination,
	}, nil
}

func (t *EmailSender) Send(title string, message string) (interface{}, error) {
	cm := fmt.Sprintf("%s %s", title, message)
	fmt.Println(cm)
	return "Telegram notification sent successfully", nil
}
