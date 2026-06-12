package email_register

import (
	"fmt"
	"github.com/bagusaditiasetiawan/saetechnology-be/internal/pkg/logger"
	"strings"
	"time"

	"github.com/bagusaditiasetiawan/saetechnology-be/internal/config"
	"github.com/bagusaditiasetiawan/saetechnology-be/internal/domain/email"
)

type UseCase interface {
	SendActivationEmail(email string, fullName string, activationToken string) error
	GetTemplate(name string, activationToken string) string
}

type UseCaseImpl struct {
	EmailSender  email.EmailSender
	SenderName   string
	SenderEmail  string
	RegisterLink string
	Logger       logger.Logger
}

func NewUseCase(emailSender email.EmailSender, config *config.Config, Logger logger.Logger) UseCase {
	return &UseCaseImpl{
		Logger:       Logger,
		EmailSender:  emailSender,
		SenderName:   config.RegisterSenderName,
		SenderEmail:  config.RegisterSenderEmail,
		RegisterLink: config.RegisterLinkUrl,
	}
}

func (useCase *UseCaseImpl) SendActivationEmail(
	receipt string,
	fullName string,
	activationToken string,
) error {
	useCase.Logger.Info("Send activation email", receipt, fullName)
	htmlContent := useCase.GetTemplate(fullName, activationToken)
	request := email.EmailRequest{
		Sender: email.Sender{
			Name:  useCase.SenderName,
			Email: useCase.SenderEmail,
		},
		Subject:     "Aktivasi Akun Ecommerce",
		HTMLContent: htmlContent,
		MessageVersions: []email.MessageVersion{
			{
				To: []email.Recipient{
					{
						Email: receipt,
						Name:  fullName,
					},
				},
				Subject:     "Aktivasi Akun Ecommerce",
				HTMLContent: htmlContent,
			},
		},
	}
	return useCase.EmailSender.SendEmail(request)
}

func (u *UseCaseImpl) GetTemplate(name string, activationToken string) string {
	if strings.TrimSpace(name) == "" {
		name = "Customer"
	}

	activationURL := fmt.Sprintf(
		"%s/activate-account?token=%s",
		strings.TrimRight(u.RegisterLink, "/"),
		activationToken,
	)

	htmlContent := fmt.Sprintf(`
		<div style="font-family: Arial, sans-serif; color: #333;">
		  <h2>Hi %s,</h2>
		  <p>Terima kasih sudah mendaftar di <strong>Ecommerce Panel</strong>.</p>
		  <p>Akun kamu sudah dibuat, tetapi belum aktif. Klik tombol di bawah ini untuk mengaktifkan akun.</p>
		  <p style="margin-top: 20px;">Link aktivasi ini berlaku selama 24 jam.</p>
		  <p style="margin: 30px 0;">
			<a href="%s"
			   style="background: #007bff; color: #fff; padding: 10px 20px; text-decoration: none; border-radius: 6px;">
			   Aktivasi Akun
			</a>
		  </p>
		  <p>Jika tombol tidak bisa diklik, salin link berikut ke browser:</p>
		  <p style="word-break: break-all; color: #555;">%s</p>
		  <p>Jika kamu tidak merasa mendaftar, abaikan email ini.</p>
		  <hr style="margin: 30px 0; border: none; border-top: 1px solid #eee;">
		  <p style="font-size: 13px; color: #888;">© %d Ecommerce Panel. All rights reserved.</p>
		</div>
		`,
		name,
		activationURL,
		activationURL,
		time.Now().Year(),
	)
	return htmlContent
}
