package logging

import (
	"context"
	"fmt"

	"math/rand"

	"github.com/apex/log"
)

const XTraceId = "trace-id"
const XSpanId = "span-id"
const Component = "component"

func Log(ctx context.Context, component string) *log.Entry {
	return log.WithFields(log.Fields{
		XTraceId:  ctx.Value(XTraceId),
		XSpanId:   ctx.Value(XSpanId),
		Component: component,
	})
}

func NewInitialContext() context.Context {
	return AddTraceToContext(context.Background(), "")
}

func AddTraceToContext(ctx context.Context, traceId string) context.Context {
	spanId := generateSpanId()
	if len(traceId) < 1 {
		traceId = spanId
	}

	return context.WithValue(
		context.WithValue(
			ctx, XTraceId, traceId),
		XSpanId, spanId)
}

func generateSpanId() string {
	b := make([]byte, 8)
	rand.Read(b)
	return fmt.Sprintf("%x", b)
}
