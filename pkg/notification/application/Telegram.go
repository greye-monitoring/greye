package application

import (
	"errors"
	"fmt"
	"github.com/go-telegram-bot-api/telegram-bot-api"
)

// TelegramSender è una struttura che implementa l'interfaccia Sender per inviare notifiche su Telegram.
type TelegramSender struct {
	Token       string
	Destination string
}

// NewTelegramSender crea una nuova istanza di TelegramSender.
func NewTelegramSender(config map[string]interface{}) (*TelegramSender, error) {
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
	return &TelegramSender{
		Token:       token,
		Destination: destination,
	}, nil
}

func (t *TelegramSender) Send(title string, message string) (interface{}, error) {
	bot, err := tgbotapi.NewBotAPI(t.Token)
	if err != nil {
		return nil, err
	}

	//for _, chatID := range destination {
	msg := tgbotapi.NewMessageToChannel(t.Destination, fmt.Sprintf("%s\n\n%s", title, message))

	_, err = bot.Send(msg)
	if err != nil {
		return nil, err
	}
	//}

	return "Telegram notification sent successfully", nil
}
