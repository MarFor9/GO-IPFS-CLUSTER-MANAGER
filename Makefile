LOCAL_DEV_PATH = $(shell pwd)/infrastructure/local
DOCKER_COMPOSE_FILE := $(LOCAL_DEV_PATH)/docker-compose.yml
DOCKER_COMPOSE_CMD := docker compose -p nbo-ipfs-cluster-manager -f $(DOCKER_COMPOSE_FILE)

GO?=$(shell which go)

.PHONY: build
build: ## Build the docker image.
	docker build --no-cache --tag IPFS-CLUSTER-MANAGER .

.PHONY: run
run: ## Run the docker image.
	$(DOCKER_COMPOSE_CMD) up -d api

.PHONY: down
down: ## Stop the docker image.
	$(DOCKER_COMPOSE_CMD) down --remove-orphans

.PHONY: tests
tests:
	$(GO) test -v ./... --count=1

.PHONY: api
api: ## Generate API files.
	oapi-codegen -o ./internal/api/api.gen.go -config api/config-oapi-codegen.yaml api/api.yaml