package email

import (
	"context"
	"fmt"
	"net/smtp"

	"github.com/Girmex/go-ecommerce/microservices/notification/internal/ports"
)

type SMTPEmailSender struct {
	host     string
	port     string
	username string
	password string
	from     string
}

func NewSMTPEmailSender(host, port, username, password, from string) ports.EmailSender {
	return &SMTPEmailSender{
		host:     host,
		port:     port,
		username: username,
		password: password,
		from:     from,
	}
}

func (s *SMTPEmailSender) Send(
	ctx context.Context,
	to string,
	subject string,
	body string,
) error {
	addr := fmt.Sprintf("%s:%s", s.host, s.port)
	auth := smtp.PlainAuth("", s.username, s.password, s.host)

	msg := []byte(fmt.Sprintf("To: %s\r\nSubject: %s\r\nContent-Type: text/plain; charset=UTF-8\r\n\r\n%s", to, subject, body))

	if err := smtp.SendMail(addr, auth, s.from, []string{to}, msg); err != nil {
		return fmt.Errorf("failed to send SMTP email: %w", err)
	}

	return nil
}
