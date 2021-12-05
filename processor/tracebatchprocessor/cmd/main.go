package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/esdrasbrz/opentelemetry-extensions/processor/tracebatchprocessor"
	"github.com/hashicorp/go-multierror"
	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/config/configmapprovider"
	"go.opentelemetry.io/collector/service"
	"go.opentelemetry.io/collector/service/defaultcomponents"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage:", os.Args[0], "CONFIG_FILE")
		return
	}

	factories, err := components()
	if err != nil {
		log.Fatalf("failed to build components: %v", err)
	}

	info := component.BuildInfo{
		Command:     "otelcol-tracebatch",
		Description: "OpenTelemetry Collector with tracebatch processor",
		Version:     "1.0.0",
	}

	app, err := service.New(
		service.CollectorSettings{
			BuildInfo:         info,
			Factories:         factories,
			ConfigMapProvider: configmapprovider.NewFile(os.Args[1]),
		},
	)
	if err != nil {
		log.Fatalf("failed to construct the application: %v", err)
	}

	err = app.Run(context.Background())
	if err != nil {
		log.Fatalf("application run finished with error: %v", err)
	}
}

func components() (component.Factories, error) {
	var errs error

	factories, err := defaultcomponents.Components()
	if err != nil {
		return component.Factories{}, err
	}
	processors := []component.ProcessorFactory{
		tracebatchprocessor.NewFactory(),
	}
	for _, pr := range factories.Processors {
		processors = append(processors, pr)
	}

	factories.Processors, err = component.MakeProcessorFactoryMap(processors...)
	if err != nil {
		errs = multierror.Append(errs, err)
	}

	return factories, errs
}
