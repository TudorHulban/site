build:
	@go build -o distrib/tara-works .

build-legacy:
	@go build -tags legacy_cpu -o distrib/tara-works .