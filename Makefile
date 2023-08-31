importdb:
	./scripts/import_bear.sh

run:
	go run main.go

all:
	docker compose up -d
	go run main.go