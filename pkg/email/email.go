package email

import (
	"fmt"
	"github.com/go-mail/mail/v2"
	"os"
	"strconv"
)

const (
	Port = 25
)

type Email struct {
	From      string
	To        string
	Subject   string
	Plaintext string
	HTML      string
}

type Service struct {
	dialer *mail.Dialer
}

func NewService() *Service {
	port, err := strconv.Atoi(os.Getenv("SMTP_PORT"))
	if err != nil {
		port = Port
	}

	return &Service{
		dialer: &mail.Dialer{
			Host:     os.Getenv("SMTP_HOST"),
			Port:     port,
			Username: os.Getenv("SMTP_USERNAME"),
			Password: os.Getenv("SMTP_PASSWORD"),
		},
	}
}

func (es *Service) Send(email Email) error {
	op := "email.Send"

	msg := mail.NewMessage()

	msg.SetHeader("From", email.From)
	msg.SetHeader("To", email.To)
	msg.SetHeader("Subject", email.Subject)

	switch {
	case email.Plaintext != "" && email.HTML != "":
		msg.SetBody("text/plain", email.Plaintext)
		msg.AddAlternative("text/html", email.HTML)
	case email.Plaintext != "":
		msg.SetBody("text/plain", email.Plaintext)
	case email.HTML != "":
		msg.AddAlternative("text/html", email.HTML)
	}

	if err := es.dialer.DialAndSend(msg); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}
