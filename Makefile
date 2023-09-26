importdb:
	./scripts/import_bear.sh

run:
	export OTEL_EXPORTER_OTLP_METRICS_ENDPOINT=http://localhost:4317
	go run main.go

docker:
	docker compose up -d

docker-reset:
	./scripts/import_bear.sh
	docker compose down
	docker compose up -d --build habits-app