package services

import (
	"gopkg.in/gomail.v2"
)

type EmailService struct {
	SMTPHost string
	SMTPPort int
	Username string
	Password string
}

func NewEmailService(smtpHost string, smtpPort int, username, password string) *EmailService {
	return &EmailService{
		SMTPHost: smtpHost,
		SMTPPort: smtpPort,
		Username: username,
		Password: password,
	}
}

func (e *EmailService) SendEmail(to, subject, body string) error {
	message := gomail.NewMessage()
	message.SetHeader("From", e.Username)
	message.SetHeader("To", to)
	message.SetHeader("Subject", subject)
	message.SetBody("text/plain", body)

	dialer := gomail.NewDialer(e.SMTPHost, e.SMTPPort, e.Username, e.Password)

	if err := dialer.DialAndSend(message); err != nil {
		return err
	}
	return nil
}
