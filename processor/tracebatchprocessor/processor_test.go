package tracebatchprocessor

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.opentelemetry.io/collector/consumer/consumertest"
	"go.opentelemetry.io/collector/model/pdata"
)

func TestSplitDifferentTracesIntoDifferentBatches(t *testing.T) {
	inBatch := pdata.NewTraces()
	rs := inBatch.ResourceSpans().AppendEmpty()

	// the first ILS has two spans
	ils := rs.InstrumentationLibrarySpans().AppendEmpty()
	library := ils.InstrumentationLibrary()
	library.SetName("first-library")

	firstSpan := ils.Spans().AppendEmpty()
	firstSpan.SetName("first-batch-first-span")
	firstSpan.SetTraceID(pdata.NewTraceID([16]byte{1, 2, 3, 4}))

	secondSpan := ils.Spans().AppendEmpty()
	secondSpan.SetName("first-batch-second-span")
	secondSpan.SetTraceID(pdata.NewTraceID([16]byte{2, 3, 4, 5}))

	// test
	next := new(consumertest.TracesSink)
	processor := newTraceBatch(next)
	err := processor.ConsumeTraces(context.Background(), inBatch)

	// verify
	assert.NoError(t, err)
	assert.Len(t, next.AllTraces(), 2)

	// first batch
	firstOutILS := next.AllTraces()[0].ResourceSpans().At(0).InstrumentationLibrarySpans().At(0)
	assert.Equal(t, library.Name(), firstOutILS.InstrumentationLibrary().Name())
	assert.Equal(t, firstSpan.Name(), firstOutILS.Spans().At(0).Name())

	// second batch
	secondOutILS := next.AllTraces()[1].ResourceSpans().At(0).InstrumentationLibrarySpans().At(0)
	assert.Equal(t, library.Name(), secondOutILS.InstrumentationLibrary().Name())
	assert.Equal(t, secondSpan.Name(), secondOutILS.Spans().At(0).Name())
}
