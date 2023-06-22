package hashivault

import "go.opentelemetry.io/otel/trace"

const defaultTracerName = "go.opentelemetry.io/otel"

var tracerName string

func traceError(span trace.Span, err error) {
	if err != nil {
		span.RecordError(err)
	}
}
