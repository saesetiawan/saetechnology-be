package tracing

import (
	"bytes"
	"io"
	"net/http"
	"strings"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

type responseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}

type Middleware struct {
	handler http.Handler
	tracer  trace.Tracer
}

func NewMiddleware(
	handler http.Handler,
	tracerProvider trace.TracerProvider,
) http.Handler {
	return &Middleware{
		handler: handler,
		tracer:  tracerProvider.Tracer("HttpServer"),
	}
}

func (m *Middleware) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	operationName := r.Method + " " + r.URL.Path

	ctx, span := m.tracer.Start(r.Context(), operationName)
	defer span.End()

	r = r.WithContext(ctx)

	body := ""
	contentType := r.Header.Get("Content-Type")
	if !strings.HasPrefix(contentType, "multipart/form-data") {
		body = readBody(r, 1024*10)
	}

	span.SetAttributes(
		attribute.String("http.method", strings.ToValidUTF8(r.Method, "_")),
		attribute.String("http.path", strings.ToValidUTF8(r.URL.Path, "_")),
		attribute.String("http.query", strings.ToValidUTF8(r.URL.RawQuery, "_")),
		attribute.String("http.user_agent", strings.ToValidUTF8(r.UserAgent(), "_")),
		attribute.String("http.remote_addr", strings.ToValidUTF8(r.RemoteAddr, "_")),
	)

	if body != "" {
		span.SetAttributes(
			attribute.String("http.request.body", strings.ToValidUTF8(maskSensitiveBody(body), "_")),
		)
	}

	rw := &responseWriter{
		ResponseWriter: w,
		statusCode:     http.StatusOK,
	}

	m.handler.ServeHTTP(rw, r)

	span.SetAttributes(
		attribute.Int("http.status_code", rw.statusCode),
	)

	if rw.statusCode >= 500 {
		span.SetStatus(codes.Error, "internal server error")
	} else if rw.statusCode >= 400 {
		span.SetStatus(codes.Error, "client error")
	} else {
		span.SetStatus(codes.Ok, "success")
	}
}

func readBody(r *http.Request, limit int64) string {
	if r.Body == nil {
		return ""
	}

	bodyBytes, err := io.ReadAll(io.LimitReader(r.Body, limit))
	if err != nil {
		return ""
	}

	r.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))

	return strings.ToValidUTF8(string(bodyBytes), "_")
}

func maskSensitiveBody(body string) string {
	sensitiveWords := []string{
		"password",
		"token",
		"access_token",
		"refresh_token",
		"authorization",
		"signature",
		"secret",
	}

	lower := strings.ToLower(body)

	for _, word := range sensitiveWords {
		if strings.Contains(lower, word) {
			return "[masked_sensitive_body]"
		}
	}

	return body
}
