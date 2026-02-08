# Makefile
.PHONY: dev

dev: install-deps
    @echo "Starting development server..."
    @go run ./cmd/dev-server --port=8080 --watch

dev-with-dashboard:
    @go run ./cmd/dev-server --port=8080 --watch --dashboard

dev-profile:
    @go run ./cmd/dev-server --port=8080 --watch --profile

# Hot reload for specific components
dev-components:
    @find ./pkg/components -name "*.go" | entr -r make dev
