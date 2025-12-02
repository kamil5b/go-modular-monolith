.PHONY: run server test test-unit test-internal lint migrate deps-check mocks-gen clean-mocks mock

# Run the application
run:
	go run .

# Run server only (skip migrations)
server:
	go run . server

# Run all tests
test:
	go test -v -race -cover ./...

# Run unit tests only
test-unit:
	go test -v -short ./...

# Run unit tests for internal packages only
test-internal:
	go test -v -short ./internal/...

# Run linter
lint:
	golangci-lint run ./...

# Check module dependencies
deps-check:
	go run cmd/lint-deps/main.go

# Run SQL migrations
migrate-up:
	go run . migration sql up

migrate-down:
	go run . migration sql down

# Run MongoDB migrations
migrate-mongo:
	go run . migration mongo up

# Generate mocks (requires mockery)
mocks:
	mockery --all --output=internal/mocks --outpkg=mocks

# Generate mocks for all interfaces under internal/ (uses mockgen + our helpers)
# Use source-mode generator to avoid import resolution issues.
mocks-gen:
	./scripts/generate_mocks_from_source.sh

# Remove generated mocks
clean-mocks:
	find internal -type d -name "mocks" -print0 | xargs -0 rm -rf || true

# Clean and (re)create mocks
mock:
	$(MAKE) clean-mocks
	$(MAKE) mocks-gen

# Build the application
build:
	go build -o bin/app .

# Clean build artifacts
clean:
	rm -rf bin/

# All checks before commit
pre-commit: deps-check lint test