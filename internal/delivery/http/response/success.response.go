package response

import (
	"context"
	"net/http"

	"go.opentelemetry.io/otel/trace"
)

type SuccessResponse interface {
	Send(
		ctx context.Context,
		w http.ResponseWriter,
		status int,
		data interface{},
	)
}

type SuccessMessageResponse struct {
	ResponseMessage string      `json:"responseMessage"`
	ResponseCode    string      `json:"responseCode"`
	ResponseData    interface{} `json:"result"`
}

type SuccessResponseImpl struct {
	jsonResponder Responder
}

func NewSuccessResponse(
	jsonResponder Responder,
) SuccessResponse {
	return &SuccessResponseImpl{
		jsonResponder: jsonResponder,
	}
}

func (s *SuccessResponseImpl) Send(
	ctx context.Context,
	w http.ResponseWriter,
	status int,
	data interface{},
) {
	span := trace.SpanFromContext(ctx)

	if span.IsRecording() {
		span.AddEvent("http.success_response")
	}

	s.jsonResponder.JSON(
		w,
		status,
		SuccessMessageResponse{
			ResponseData:    data,
			ResponseCode:    "200",
			ResponseMessage: "successful",
		},
	)
}
