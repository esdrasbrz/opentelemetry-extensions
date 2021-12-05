build: get_builder
	builder --config builder.yaml --output-path=./dist

get_builder:
	GO111MODULE=on GOBIN=/usr/local/bin go install go.opentelemetry.io/collector/cmd/builder@latest