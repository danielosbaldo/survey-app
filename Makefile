.PHONY: run compose-up compose-down compose-prod deploy deploy-local clean help

# Development
run:
	go run ./cmd/server

# Docker Compose - Development
compose-up:
	docker compose up --build

compose-down:
	docker compose down -v

# Docker Compose - Production
compose-prod:
	docker compose -f docker-compose.prod.yml up --build -d

compose-prod-down:
	docker compose -f docker-compose.prod.yml down

compose-prod-logs:
	docker compose -f docker-compose.prod.yml logs -f

# Deployment
deploy:
	@echo "Running deployment script..."
	./deploy.sh

deploy-remote:
	@echo "Remote deployment - Usage: make deploy-remote SERVER=user@host"
	@if [ -z "$(SERVER)" ]; then \
		echo "ERROR: SERVER variable required"; \
		echo "Example: make deploy-remote SERVER=daniel@192.168.1.100"; \
		exit 1; \
	fi
	./deploy-remote.sh $(SERVER)

deploy-local:
	@echo "Testing deployment locally..."
	DEPLOY_DIR=./deploy_test ./deploy.sh

# Cleanup
clean:
	docker compose down -v
	docker compose -f docker-compose.prod.yml down -v
	docker system prune -f

# Help
help:
	@echo "Available targets:"
	@echo "  run              - Run the application locally with Go"
	@echo "  compose-up       - Start development environment with Docker Compose"
	@echo "  compose-down     - Stop development environment"
	@echo "  compose-prod     - Start production environment"
	@echo "  compose-prod-down - Stop production environment"
	@echo "  compose-prod-logs - View production logs"
	@echo "  deploy           - Deploy to production (run on server)"
	@echo "  deploy-remote    - Deploy to remote server (SERVER=user@host)"
	@echo "  deploy-local     - Test deployment locally"
	@echo "  clean            - Clean up Docker resources"
	@echo "  help             - Show this help message"
	@echo ""
	@echo "Examples:"
	@echo "  make deploy-remote SERVER=daniel@192.168.1.100"
	@echo "  make deploy-remote SERVER=root@myserver.com"

