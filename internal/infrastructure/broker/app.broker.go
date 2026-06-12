package broker

import (
	"context"
	"encoding/json"
	"saetechnology-be/internal/domain/broker"
	"saetechnology-be/internal/usecase/email_register"
	"saetechnology-be/internal/usecase/publish_register"
)

type UserConsumer struct {
	emailUseCase email_register.UseCase
}

func NewUserConsumer(u email_register.UseCase) *UserConsumer {
	return &UserConsumer{emailUseCase: u}
}

func (c *UserConsumer) Handle(ctx context.Context, msg broker.Message) error {
	var payload publish_register.RegisterEmailMessage

	if err := json.Unmarshal(msg.Value, &payload); err != nil {
		return err
	}

	return c.emailUseCase.SendActivationEmail(
		payload.Email,
		payload.FullName,
		payload.ActivationToken,
	)
}
