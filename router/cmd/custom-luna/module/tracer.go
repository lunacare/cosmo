package module

import (
	"context"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
)

// tracerFromContext returns a tracer scoped to instrumentationName using the
// TracerProvider attached to the current span (i.e. the same provider Cosmo
// configured at startup). Falls back to the global provider if no span is in
// the context. Mirrors rtrace.TracerFromContext but lets us pick our own
// instrumentation name for queryability in Grafana/Tempo.
func tracerFromContext(ctx context.Context, instrumentationName string) trace.Tracer {
	if span := trace.SpanFromContext(ctx); span.SpanContext().IsValid() {
		return span.TracerProvider().Tracer(instrumentationName)
	}
	return otel.Tracer(instrumentationName)
}
