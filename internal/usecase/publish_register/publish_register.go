package publish_register

import (
	"context"
	"encoding/json"

	"go-platform-core/internal/config"
	"go-platform-core/internal/domain/broker"
	"go-platform-core/internal/domain/user"
	"go-platform-core/internal/pkg/logger"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

type RegisterEmailMessage struct {
	UserID          string `json:"user_id"`
	Email           string `json:"email"`
	FullName        string `json:"full_name"`
	Role            string `json:"role"`
	ActivationToken string `json:"activation_token"`
}

type UseCase interface {
	SendToQueue(ctx context.Context, user user.User, activationToken string)
}

type UseCaseImpl struct {
	Publisher  broker.Publisher
	Logger     logger.Logger
	QueueEmail string
	Trace      trace.Tracer
}

func NewUseCase(
	cfg *config.Config,
	logger logger.Logger,
	publisher broker.Publisher,
	traceProvider trace.TracerProvider,
) UseCase {
	return &UseCaseImpl{
		Publisher:  publisher,
		Logger:     logger,
		QueueEmail: cfg.QueueRegisterEmail,
		Trace:      traceProvider.Tracer("PublishRegisterUseCase"),
	}
}

func (uc *UseCaseImpl) SendToQueue(
	ctx context.Context,
	user user.User,
	activationToken string,
) {
	ctx, span := uc.Trace.Start(ctx, "PublishRegisterUseCase.SendToQueue")
	defer span.End()

	span.SetAttributes(
		attribute.String("queue.name", uc.QueueEmail),
		attribute.String("user.id", user.ID.String()),
		attribute.String("user.email", user.Email),
		attribute.String("user.role", user.Role),
	)

	if uc.QueueEmail == "" {
		span.SetStatus(codes.Error, "register email queue is empty")
		uc.Logger.Error("register email queue is empty")
		return
	}

	message := RegisterEmailMessage{
		UserID:          user.ID.String(),
		Email:           user.Email,
		FullName:        user.FullName,
		Role:            user.Role,
		ActivationToken: activationToken,
	}

	byteMessage, err := json.Marshal(message)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, "failed marshal register email message")

		uc.Logger.Error("Error marshalling register email message", logger.Field{
			Key:   "message",
			Value: message,
		})

		return
	}

	span.AddEvent("publish_message")

	if err := uc.Publisher.Publish(
		ctx,
		uc.QueueEmail,
		broker.Message{
			Key:   "email",
			Value: byteMessage,
		},
	); err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, "failed publish register email message")
		return
	}

	span.SetAttributes(
		attribute.Int("message.size", len(byteMessage)),
	)

	span.SetStatus(codes.Ok, "message published")
}
