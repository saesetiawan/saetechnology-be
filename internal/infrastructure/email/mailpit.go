package email

import (
	"bytes"
	"fmt"
	"net/smtp"
	"strings"

	emailDomain "go-platform-core/internal/domain/email"
)

type MailpitClient struct {
	Addr string
}

func (c *MailpitClient) SendEmail(payload emailDomain.EmailRequest) error {
	addr := strings.TrimSpace(c.Addr)
	if addr == "" {
		addr = "localhost:1025"
	}

	sender := payload.Sender.Email
	if strings.TrimSpace(sender) == "" {
		sender = "no-reply@localhost"
	}

	messages := payload.MessageVersions
	if len(messages) == 0 {
		messages = []emailDomain.MessageVersion{
			{
				To:          nil,
				Subject:     payload.Subject,
				HTMLContent: payload.HTMLContent,
			},
		}
	}

	for _, message := range messages {
		subject := firstNonEmpty(message.Subject, payload.Subject)
		htmlContent := firstNonEmpty(message.HTMLContent, payload.HTMLContent)
		if htmlContent == "" {
			htmlContent = "<p></p>"
		}

		for _, recipient := range message.To {
			if strings.TrimSpace(recipient.Email) == "" {
				continue
			}

			body := buildSMTPMessage(sender, recipient, subject, htmlContent)
			if err := smtp.SendMail(addr, nil, sender, []string{recipient.Email}, body); err != nil {
				return err
			}
		}
	}

	return nil
}

func buildSMTPMessage(
	sender string,
	recipient emailDomain.Recipient,
	subject string,
	htmlContent string,
) []byte {
	to := recipient.Email
	if strings.TrimSpace(recipient.Name) != "" {
		to = fmt.Sprintf("%s <%s>", recipient.Name, recipient.Email)
	}

	var buffer bytes.Buffer
	buffer.WriteString("From: " + sender + "\r\n")
	buffer.WriteString("To: " + to + "\r\n")
	buffer.WriteString("Subject: " + subject + "\r\n")
	buffer.WriteString("MIME-Version: 1.0\r\n")
	buffer.WriteString("Content-Type: text/html; charset=UTF-8\r\n")
	buffer.WriteString("\r\n")
	buffer.WriteString(htmlContent)

	return buffer.Bytes()
}

func firstNonEmpty(values ...string) string {
	for _, value := range values {
		if strings.TrimSpace(value) != "" {
			return value
		}
	}

	return ""
}
