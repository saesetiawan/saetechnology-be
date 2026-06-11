package broker

import (
	"context"

	rabbitmq "github.com/wagslane/go-rabbitmq"

	"go-platform-core/internal/domain/broker"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/trace"
)

type rabbitPublisher struct {
	publisher *rabbitmq.Publisher
	tracer    trace.Tracer
}

func NewRabbitPublisher(
	conn *rabbitmq.Conn,
	tracerProvider trace.TracerProvider,
) broker.Publisher {
	pub, err := rabbitmq.NewPublisher(
		conn,
		rabbitmq.WithPublisherOptionsExchangeName(""),
	)
	if err != nil {
		panic(err)
	}

	return &rabbitPublisher{
		publisher: pub,
		tracer:    tracerProvider.Tracer("RabbitPublisher"),
	}
}

func (r *rabbitPublisher) Publish(
	ctx context.Context,
	topic string,
	msg broker.Message,
) error {
	ctx, span := r.tracer.Start(
		ctx,
		"RabbitMQ.Publish",
		trace.WithSpanKind(trace.SpanKindProducer),
	)
	defer span.End()

	span.SetAttributes(
		attribute.String("messaging.system", "rabbitmq"),
		attribute.String("messaging.operation", "publish"),
		attribute.String("messaging.destination.name", topic),
		attribute.Int("messaging.message.body.size", len(msg.Value)),
	)

	headers := map[string]interface{}{}
	injectTraceContext(ctx, headers)

	err := r.publisher.Publish(
		msg.Value,
		[]string{topic}, // topic = queue name kalau pakai default exchange
		rabbitmq.WithPublishOptionsContentType("application/json"),
		rabbitmq.WithPublishOptionsPersistentDelivery,
		rabbitmq.WithPublishOptionsHeaders(headers),
	)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, "rabbitmq publish failed")
		return err
	}

	span.SetStatus(codes.Ok, "rabbitmq publish success")

	return nil
}

func injectTraceContext(
	ctx context.Context,
	headers map[string]interface{},
) {
	carrier := propagation.MapCarrier{}

	otel.GetTextMapPropagator().Inject(ctx, carrier)

	for key, value := range carrier {
		headers[key] = value
	}
}
