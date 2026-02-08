# Makefile

.PHONY: install-deps wasm server build run dev

install-deps:
	@echo "Installing dependencies..."
	go get -u github.com/fsnotify/fsnotify
	go get -u github.com/gorilla/websocket
	go get -u github.com/maxence-charriere/go-app/v10

# Build the Frontend (WebAssembly)
# Note: GOOS=js and GOARCH=wasm are required for go-app to run in the browser
wasm:
	@echo "Building WebAssembly..."
	GOOS=js GOARCH=wasm go build -o web/app.wasm ./cmd/wasm

# Build the Backend (Server)
server:
	@echo "Building Server..."
	go build -o server ./cmd/server

build: install-deps wasm server

# Run the application
run: build
	@echo "Starting server at http://localhost:8080"
	./server

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