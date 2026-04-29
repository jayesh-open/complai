.PHONY: dev down db-up migrate migrate-all seed lint test build clean fmt help

DOCKER_COMPOSE := docker compose -f docker-compose.dev.yml
GO_SERVICES := $(shell find services/go -name 'go.mod' -exec dirname {} \;)

help: ## Show this help
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'

dev: ## Start all dev services (docker-compose)
	$(DOCKER_COMPOSE) up -d
	@echo "Waiting for services to be healthy..."
	@$(DOCKER_COMPOSE) exec -T postgres pg_isready -U complai -q && echo "Postgres: ready" || echo "Postgres: waiting..."
	@echo "Dev environment started. Services:"
	@echo "  Postgres:     localhost:5432"
	@echo "  Redis:        localhost:6379"
	@echo "  LocalStack:   localhost:4566"
	@echo "  OpenSearch:   localhost:9200"
	@echo "  Keycloak:     http://localhost:8080"
	@echo "  Temporal:     localhost:7233"
	@echo "  Temporal UI:  http://localhost:8088"
	@echo "  Mailpit:      http://localhost:8025"
	@echo "  Jaeger:       http://localhost:16686"

down: ## Stop all dev services
	$(DOCKER_COMPOSE) down

down-clean: ## Stop all dev services and remove volumes
	$(DOCKER_COMPOSE) down -v

db-up: ## Start only Postgres
	$(DOCKER_COMPOSE) up -d postgres

migrate: ## Run database migrations for all services
	@for svc in $(GO_SERVICES); do \
		if [ -d "$$svc/migrations" ] && [ "$$(ls -A $$svc/migrations 2>/dev/null)" ]; then \
			echo "Migrating: $$svc"; \
			cd $$svc && go run cmd/server/main.go migrate 2>/dev/null || true; \
			cd $(CURDIR); \
		fi; \
	done

MIGRATE_ORDER := \
	services/go/identity-service \
	services/go/tenant-service \
	services/go/user-role-service \
	services/go/master-data-service \
	services/go/document-service \
	services/go/notification-service \
	services/go/audit-service \
	services/go/workflow-service \
	services/go/rules-engine-service \
	services/go/gst-service \
	services/go/vendor-compliance-service \
	services/go/recon-service \
	services/go/einvoice-service \
	services/go/ewb-service \
	services/go/tds-service

migrate-all: ## Run migrations for all services in dependency order (stops on failure)
	@failed=0; \
	for svc in $(MIGRATE_ORDER); do \
		if [ -d "$$svc/migrations" ] && [ "$$(ls -A $$svc/migrations 2>/dev/null)" ]; then \
			echo "Migrating: $$svc"; \
			(cd $$svc && go run cmd/server/main.go migrate) || { echo "FAILED: $$svc"; failed=1; break; }; \
		fi; \
	done; \
	if [ $$failed -eq 0 ]; then echo "All migrations applied successfully."; else exit 1; fi

seed: ## Seed development data
	@if [ -f scripts/seed-dev.sh ]; then \
		bash scripts/seed-dev.sh; \
	else \
		echo "No seed script found"; \
	fi

lint: ## Run linters across all workspaces
	@echo "=== Go lint ==="
	@for svc in $(GO_SERVICES); do \
		echo "Linting: $$svc"; \
		cd $$svc && go vet ./... && cd $(CURDIR); \
	done
	@echo "=== Node lint ==="
	pnpm -r run lint

test: ## Run tests across all workspaces
	@echo "=== Go tests ==="
	@for svc in $(GO_SERVICES); do \
		echo "Testing: $$svc"; \
		cd $$svc && go test ./... && cd $(CURDIR); \
	done
	@echo "=== Node tests ==="
	pnpm -r run test

test-go: ## Run Go tests only
	@for svc in $(GO_SERVICES); do \
		echo "Testing: $$svc"; \
		cd $$svc && go test ./... && cd $(CURDIR); \
	done

test-node: ## Run Node tests only
	pnpm -r run test

typecheck: ## TypeScript type checking
	pnpm -r run typecheck

build: ## Build all services
	@echo "=== Go build ==="
	@for svc in $(GO_SERVICES); do \
		echo "Building: $$svc"; \
		cd $$svc && go build ./... && cd $(CURDIR); \
	done
	@echo "=== Node build ==="
	pnpm -r run build

build-go: ## Build Go services only
	@for svc in $(GO_SERVICES); do \
		echo "Building: $$svc"; \
		cd $$svc && go build ./... && cd $(CURDIR); \
	done

build-node: ## Build Node packages only
	pnpm -r run build

fmt: ## Format all code
	@echo "=== Go fmt ==="
	@for svc in $(GO_SERVICES); do \
		cd $$svc && gofmt -w . && cd $(CURDIR); \
	done
	@echo "=== Prettier ==="
	pnpm format

clean: ## Clean build artifacts
	@for svc in $(GO_SERVICES); do \
		cd $$svc && rm -rf bin/ && cd $(CURDIR); \
	done
	pnpm -r run clean
	rm -rf .turbo

health-probe: ## Run health-probe-service locally
	cd services/go/health-probe-service && go run cmd/server/main.go

localstack-queues: ## List SQS queues in LocalStack
	docker exec complai-localstack awslocal sqs list-queues --region ap-south-1

localstack-topics: ## List SNS topics in LocalStack
	docker exec complai-localstack awslocal sns list-topics --region ap-south-1

localstack-buckets: ## List S3 buckets in LocalStack
	docker exec complai-localstack awslocal s3 ls --region ap-south-1
