package email

import (
	"strings"

	"github.com/bagusaditiasetiawan/saetechnology-be/internal/config"
	emailDomain "github.com/bagusaditiasetiawan/saetechnology-be/internal/domain/email"
)

func NewEmailSender(cfg *config.Config) emailDomain.EmailSender {
	if strings.EqualFold(strings.TrimSpace(cfg.AppEnv), "production") {
		return NewBrevoClient(cfg)
	}

	return &MailpitClient{
		Addr: cfg.MailpitSMTPAddr,
	}
}
