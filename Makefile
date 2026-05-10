COMPOSE ?= docker compose
SERVICE ?= dev

.PHONY: docker-build docker-test docker-run docker-shell docker-fmt docker-clean

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

docker-clean:
	$(COMPOSE) down --remove-orphans
