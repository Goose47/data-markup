.PHONY: migrate

migrate:
	@docker build -t migrator-image --build-arg BUILD_PATH=./cmd/migrator/migrator.go .
	@docker run --rm -v ./migrations:/app/migrations --network data-markup_markup migrator-image ./app --config=./config/local.yml --migrations-path=migrations