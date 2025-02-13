.PHONY: migrate psql

migrate:
	@docker build -t migrator-image --build-arg BUILD_PATH=./cmd/migrator/migrator.go .
	@docker run --rm -v ./migrations:/app/migrations --network data-markup_markup migrator-image ./app --config=./config/dev.yml --migrations-path=migrations

psql:
	@psql -h $(docker inspect   -f '{{range.NetworkSettings.Networks}}{{.IPAddress}}{{end}}' postgres) -p 5432 -U root markup