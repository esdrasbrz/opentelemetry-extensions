package tracebatchprocessor

import "go.opentelemetry.io/collector/model/pdata"

func splitByTrace(rs pdata.ResourceSpans) []pdata.ResourceSpans {
	var result []pdata.ResourceSpans

	for i := 0; i < rs.InstrumentationLibrarySpans().Len(); i++ {
		// the batches for this ILS
		batches := map[string]pdata.ResourceSpans{}

		ils := rs.InstrumentationLibrarySpans().At(i)
		for j := 0; j < ils.Spans().Len(); j++ {
			span := ils.Spans().At(j)
			key := span.TraceID().HexString()

			// initialize for the first traceID in the ILS
			if _, ok := batches[key]; !ok {
				newRS := pdata.NewResourceSpans()
				rs.Resource().CopyTo(newRS.Resource())

				newILS := newRS.InstrumentationLibrarySpans().AppendEmpty()
				ils.InstrumentationLibrary().CopyTo(newILS.InstrumentationLibrary())
				batches[key] = newRS

				result = append(result, newRS)
			}

			// there is only one instrumentation library per batch
			span.CopyTo(batches[key].InstrumentationLibrarySpans().At(0).Spans().AppendEmpty())
		}
	}

	return result
}
