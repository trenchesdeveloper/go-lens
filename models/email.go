package models

import (
	"fmt"

	"github.com/go-mail/mail/v2"
)

const (
	DefaultSender = "support@samuel.com"
)

type Email struct {
	From      string
	To        []string
	Subject   string
	PlainText string
	HTML      string
}

type SMTPConfig struct {
	Host     string
	Port     int
	Username string
	Password string
}

type EmailService struct {
	DefaultSender string

	dialer *mail.Dialer
}

func NewEmailService(config SMTPConfig) *EmailService {
	es := &EmailService{
		DefaultSender: DefaultSender,
		dialer:        mail.NewDialer(config.Host, config.Port, config.Username, config.Password),
	}

	return es
}

func (es *EmailService) Send(email Email) error {
	m := mail.NewMessage()
	// Set the sender address
	es.setForm(m, email)
	m.SetHeader("To", email.To...)
	m.SetHeader("Subject", email.Subject)

	switch {
		case email.PlainText != "" && email.HTML != "":
			m.AddAlternative("text/html", email.HTML)
			m.SetBody("text/plain", email.PlainText)
		case email.PlainText != "":
			m.SetBody("text/plain", email.PlainText)
		case email.HTML != "":
			m.SetBody("text/html", email.HTML)
	}

	err := es.dialer.DialAndSend(m)

	if err != nil {
		return fmt.Errorf("Send: %w", err)
	}

	return nil
}


func (es *EmailService) setForm(msg *mail.Message, email Email) {
	var from string

	switch {
		case email.From != "":
			from = email.From
		case es.DefaultSender != "":
			from = es.DefaultSender
		default:
			from = DefaultSender
		}

	msg.SetHeader("From", from)
}

func (es *EmailService) ForgotPasswordEmail(to []string, resetURL string) error {
	email := Email{
		To: to,
		Subject: "Reset your password",
		PlainText: fmt.Sprintf("Your reset token is %s", resetURL),
		HTML: fmt.Sprintf(`<p>Your reset token is <strong>%s</strong></p>`, resetURL),
	}

	err := es.Send(email)

	if err != nil {
		return fmt.Errorf("ForgotPasswordEmail: %w", err)
	}

	return nil
}
