package tracebatchprocessor

import (
	"context"

	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/consumer"
	"go.opentelemetry.io/collector/model/pdata"
)

type traceBatch struct {
	next consumer.Traces
}

func newTraceBatch(next consumer.Traces) *traceBatch {
	tb := &traceBatch{next: next}
	return tb
}

func (tb *traceBatch) ConsumeTraces(ctx context.Context, batch pdata.Traces) error {
	for i := 0; i < batch.ResourceSpans().Len(); i++ {
		rss := splitByTrace(batch.ResourceSpans().At(i))
		for _, newBatch := range rss {
			trace := pdata.NewTraces()
			newBatch.CopyTo(trace.ResourceSpans().AppendEmpty())
			if err := tb.next.ConsumeTraces(ctx, trace); err != nil {
				// fail fast for the entire batch
				return err
			}
		}
	}
	return nil
}

func (tb *traceBatch) Start(ctx context.Context, host component.Host) error {
	return nil
}

func (tb *traceBatch) Shutdown(ctx context.Context) error {
	return nil
}

func (tb *traceBatch) Capabilities() consumer.Capabilities {
	return consumer.Capabilities{MutatesData: true}
}
