.PHONY: migrate

migrate:
	@docker build migrator-image --build-arg BUILD_PATH=./cmd/migrator/migrator.go .
	@docker run --rm migrator-image -v ./migrations /app/migrations ./app --config=./config/local.yml --migrations-path=migrations