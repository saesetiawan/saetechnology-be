package broker

import (
	"context"
	"log"
	"strconv"

	rabbitmq "github.com/wagslane/go-rabbitmq"

	"github.com/bagusaditiasetiawan/saetechnology-be/internal/domain/broker"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/trace"
)

type rabbitConsumer struct {
	conn   *rabbitmq.Conn
	tracer trace.Tracer
}

func NewRabbitConsumer(
	conn *rabbitmq.Conn,
	tracerProvider trace.TracerProvider,
) broker.Consumer {
	return &rabbitConsumer{
		conn:   conn,
		tracer: tracerProvider.Tracer("RabbitConsumer"),
	}
}

func (r *rabbitConsumer) Start(
	ctx context.Context,
	queue string,
	handler broker.Handler,
) error {
	ctx, span := r.tracer.Start(ctx, "RabbitMQ.Consumer.Start")
	defer span.End()

	span.SetAttributes(
		attribute.String("messaging.system", "rabbitmq"),
		attribute.String("messaging.destination.name", queue),
		attribute.String("messaging.operation", "start"),
		attribute.Int("messaging.consumer.concurrency", 5),
	)

	consumer, err := rabbitmq.NewConsumer(
		r.conn,
		queue,
		rabbitmq.WithConsumerOptionsConcurrency(5),

		// PENTING:
		// Jangan set RoutingKey dan ExchangeName kosong.
		// Kalau queue sudah dideclare, consumer cukup consume queue-nya.
	)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, "failed create rabbitmq consumer")
		return err
	}

	go func() {
		err := consumer.Run(func(d rabbitmq.Delivery) rabbitmq.Action {
			messageCtx := extractTraceContext(ctx, d.Headers)

			messageCtx, messageSpan := r.tracer.Start(
				messageCtx,
				"RabbitMQ.Consumer.Handle",
				trace.WithSpanKind(trace.SpanKindConsumer),
			)
			defer messageSpan.End()

			messageSpan.SetAttributes(
				attribute.String("messaging.system", "rabbitmq"),
				attribute.String("messaging.operation", "consume"),
				attribute.String("messaging.destination.name", queue),
				attribute.String("messaging.rabbitmq.routing_key", d.RoutingKey),
				attribute.Int("messaging.message.body.size", len(d.Body)),
				attribute.String("messaging.message.content_type", d.ContentType),
				attribute.String("messaging.rabbitmq.delivery_tag", strconv.FormatUint(d.DeliveryTag, 10)),
				attribute.Bool("messaging.rabbitmq.redelivered", d.Redelivered),
			)

			err := handler(messageCtx, broker.Message{
				Key:   d.RoutingKey,
				Value: d.Body,
			})
			if err != nil {
				log.Println("consumer error:", err)

				messageSpan.RecordError(err)
				messageSpan.SetStatus(codes.Error, "rabbitmq message handler failed")

				return rabbitmq.NackRequeue
			}

			messageSpan.SetStatus(codes.Ok, "rabbitmq message handled")

			return rabbitmq.Ack
		})

		if err != nil {
			log.Println("consumer stopped:", err)
		}
	}()

	span.SetStatus(codes.Ok, "rabbitmq consumer started")

	return nil
}

func extractTraceContext(
	parent context.Context,
	headers map[string]interface{},
) context.Context {
	carrier := propagation.MapCarrier{}

	for key, value := range headers {
		switch v := value.(type) {
		case string:
			carrier[key] = v
		case []byte:
			carrier[key] = string(v)
		}
	}

	return otel.GetTextMapPropagator().Extract(parent, carrier)
}
