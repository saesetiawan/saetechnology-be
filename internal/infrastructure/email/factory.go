package email

import (
	"strings"

	"go-platform-core/internal/config"
	emailDomain "go-platform-core/internal/domain/email"
)

func NewEmailSender(cfg *config.Config) emailDomain.EmailSender {
	if strings.EqualFold(strings.TrimSpace(cfg.AppEnv), "production") {
		return NewBrevoClient(cfg)
	}

	return &MailpitClient{
		Addr: cfg.MailpitSMTPAddr,
	}
}
