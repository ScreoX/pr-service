.PHONY: lint lint-install build up down down-clean start start-with-build test test-unit test-integration test-down

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

start: up
	@echo "Ready: http://localhost:$(APP_PORT)"

start-with-build: build up
	@echo "Ready: http://localhost:$(APP_PORT)"

test: test-unit test-integration

test-unit:
	go test ./internal/... -v

test-integration:
	docker-compose -f docker-compose.test.yml --env-file .env.test up -d

	@until docker-compose -f docker-compose.test.yml --env-file .env.test exec test-postgres pg_isready -U $(TEST_DB_USER) -d $(TEST_DB_NAME); do \
		sleep 2; \
	done

	go test -v -tags=integration ./tests/integration/...

	docker-compose -f docker-compose.test.yml --env-file .env.test down -v

test-down:
	docker-compose -f docker-compose.test.yml --env-file .env.test down -v