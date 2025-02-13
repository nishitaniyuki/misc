package main

import (
	"context"
	"log/slog"

	otel_trace "go.opentelemetry.io/otel/trace"
)

type LogWithTraceHandler struct {
	inner slog.Handler
}

func NewLogWithTraceHandler(inner slog.Handler) *LogWithTraceHandler {
	return &LogWithTraceHandler{inner}
}

func (h *LogWithTraceHandler) Handle(ctx context.Context, r slog.Record) error {
	sc := otel_trace.SpanContextFromContext(ctx)
	if sc.IsValid() {
		r.AddAttrs(
			slog.String("trace_id", sc.TraceID().String()),
			slog.String("span_id", sc.SpanID().String()),
		)
	}
	return h.inner.Handle(ctx, r)
}

func (h *LogWithTraceHandler) Enabled(ctx context.Context, level slog.Level) bool {
	return h.inner.Enabled(ctx, level)
}

func (h *LogWithTraceHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return &LogWithTraceHandler{h.inner.WithAttrs(attrs)}
}

func (h *LogWithTraceHandler) WithGroup(name string) slog.Handler {
	return &LogWithTraceHandler{h.inner.WithGroup(name)}
}
