package tracebatchprocessor

import (
	"context"

	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/config"
	"go.opentelemetry.io/collector/consumer"
	"go.opentelemetry.io/collector/processor/processorhelper"
)

const (
	// type key in configuration
	typeStr = "tracebatch"
)

func NewFactory() component.ProcessorFactory {
	return processorhelper.NewFactory(
		typeStr,
		createDefaultConfig,
		processorhelper.WithTraces(createTraceProcessor),
	)
}

func createDefaultConfig() config.Processor {
	return &Config{
		ProcessorSettings: config.NewProcessorSettings(
			config.NewComponentID(typeStr),
		),
	}
}

func createTraceProcessor(
	_ context.Context,
	_ component.ProcessorCreateSettings,
	_ config.Processor,
	nextConsumer consumer.Traces,
) (component.TracesProcessor, error) {
	return newTraceBatch(nextConsumer), nil
}
