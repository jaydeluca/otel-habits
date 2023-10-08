importdb:
	./scripts/import_bear.sh

run:
	export OTEL_EXPORTER_OTLP_METRICS_ENDPOINT=http://localhost:4317 \
	go run main.go

run-dummy:
	export GENERATE_DATA_DAY_COUNT=10; \
	export OTEL_EXPORTER_OTLP_METRICS_ENDPOINT=http://localhost:4317; \
	go run main.go

docker:
	docker compose up -d

docker-reset:
	./scripts/import_bear.sh
	docker compose down
	docker compose up -d --build habits-app