.PHONY: build run run-postgres stop test migrate logs clean help

DOCKER_COMPOSE = docker compose
APP_NAME = slugkiller

build:
	@$(DOCKER_COMPOSE) build

run:
	@$(DOCKER_COMPOSE) up -d

run-postgres:
	@STORAGE_TYPE=postgres $(DOCKER_COMPOSE) --profile postgres up -d

stop:
	@$(DOCKER_COMPOSE) down

test:
	@$(DOCKER_COMPOSE) -f docker-compose.test.yml up --build --abort-on-container-exit

migrate:
	@$(DOCKER_COMPOSE) run --rm migration

logs:
	@$(DOCKER_COMPOSE) logs -f $(APP_NAME)

clean:
	@$(DOCKER_COMPOSE) down -v
	@docker rmi -f $(shell docker images -q $(APP_NAME)) 2>/dev/null || true

help:
	@echo "Available targets:"
	@echo "  build        - Build containers"
	@echo "  run          - Start service with in-memory storage"
	@echo "  run-postgres - Start with PostgreSQL"
	@echo "  stop         - Stop all containers"
	@echo "  test         - Run unit tests"
	@echo "  migrate      - Run database migrations"
	@echo "  logs         - View service logs"
	@echo "  clean        - Remove all containers and images"