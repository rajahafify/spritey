COMPOSE ?= docker compose
SERVICE ?= dev
GO ?= go

.PHONY: ci native-ci docker-build docker-test docker-run docker-shell docker-fmt docker-ci docker-clean

ci:
	@if command -v docker >/dev/null 2>&1 && docker compose version >/dev/null 2>&1; then \
		$(MAKE) docker-ci; \
	else \
		echo "Docker Compose not found; running native Go fallback."; \
		$(MAKE) native-ci; \
	fi

native-ci:
	$(GO) test ./...
	$(GO) build -buildvcs=false -o bin/spritey ./cmd/spritey

docker-build:
	$(COMPOSE) build

docker-test:
	$(COMPOSE) run --rm $(SERVICE) go test ./...

docker-run:
	$(COMPOSE) run --rm $(SERVICE) go run ./cmd/spritey

docker-shell:
	$(COMPOSE) run --rm $(SERVICE) bash

docker-fmt:
	$(COMPOSE) run --rm $(SERVICE) gofmt -w ./cmd ./app

docker-ci: docker-build docker-test
	$(COMPOSE) run --rm $(SERVICE) go build -buildvcs=false -o bin/spritey ./cmd/spritey

docker-clean:
	$(COMPOSE) down --remove-orphans
