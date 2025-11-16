.PHONY: lint lint-install build up down down-clean migrate start test test-unit test-integration test-down

include .env
export

lint-install:
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	go install golang.org/x/tools/cmd/goimports@latest

lint:
	golangci-lint run ./... --config=.golinter.yml

build:
	docker-compose build --no-cache

up:
	docker-compose up -d
	@until docker-compose exec postgres pg_isready -U $(DB_USER) -d $(DB_NAME); do \
		sleep 2; \
	done

down:
	docker-compose down

down-clean:
	docker-compose down -v

migrate:
	docker run --rm -v $(PWD)/migrations:/migrations --network pr-service_default \
		golang:1.25-alpine \
		sh -c 'apk add --no-cache postgresql-client && \
			go install github.com/pressly/goose/v3/cmd/goose@latest && \
			/go/bin/goose -dir /migrations postgres "user=$(DB_USER) password=$(DB_PASSWORD) dbname=$(DB_NAME) host=postgres port=5432 sslmode=disable" up'

start: build up migrate
	@echo "Ready: http://localhost:$(APP_PORT)"

test: test-unit test-integration

test-unit:
	go test ./internal/... -v

test-integration:
	docker-compose -f docker-compose.test.yml --env-file .env.test up -d

	@until docker-compose -f docker-compose.test.yml --env-file .env.test exec test-postgres pg_isready -U $(TEST_DB_USER) -d $(TEST_DB_NAME); do \
		sleep 2; \
	done

	docker run --rm -v $(PWD)/migrations:/migrations --network pr-service_default \
		--env-file .env.test \
		golang:1.25-alpine \
		sh -c 'apk add --no-cache postgresql-client && \
			go install github.com/pressly/goose/v3/cmd/goose@latest && \
			sleep 2 && \
			/go/bin/goose -dir /migrations postgres "user=$$TEST_DB_USER password=$$TEST_DB_PASSWORD dbname=$$TEST_DB_NAME host=test-postgres port=5432 sslmode=disable" up'

	go test -v -tags=integration ./tests/integration/...

	docker-compose -f docker-compose.test.yml --env-file .env.test down -v

test-down:
	docker-compose -f docker-compose.test.yml --env-file .env.test down -v