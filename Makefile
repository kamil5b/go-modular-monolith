.PHONY: help run server worker test test-unit test-internal lint proto proto-% migrate deps-check mocks-gen clean-mocks mock

# Display available targets
help:
	@echo "Available targets:"
	@echo "  Application:"
	@echo "    run              - Run the application (migrations + server)"
	@echo "    server           - Run server only (skip migrations)"
	@echo "    worker           - Run worker only"
	@echo "    build            - Build the application"
	@echo ""
	@echo "  Testing:"
	@echo "    test             - Run all tests with race detection"
	@echo "    test-unit        - Run unit tests only"
	@echo "    test-internal    - Run tests for internal packages only"
	@echo ""
	@echo "  Code Quality:"
	@echo "    lint             - Run golangci-lint"
	@echo "    deps-check       - Check module dependency violations"
	@echo ""
	@echo "  Protocol Buffers:"
	@echo "    proto            - Generate protobuf code for all modules"
	@echo "    proto-<module>   - Generate protobuf code for specific module"
	@echo "                       (e.g., make proto-product)"
	@echo ""
	@echo "  Database:"
	@echo "    migrate-up       - Apply SQL migrations"
	@echo "    migrate-down     - Rollback SQL migrations"
	@echo "    migrate-mongo    - Apply MongoDB migrations"
	@echo ""
	@echo "  Mocks:"
	@echo "    mocks-gen        - Generate mocks from source"
	@echo "    clean-mocks      - Remove generated mocks"
	@echo "    mock             - Clean and regenerate all mocks"
	@echo ""
	@echo "  CI/CD:"
	@echo "    pre-commit       - Run all checks before committing"
	@echo "    clean            - Clean build artifacts"

# Run the application
run:
	go run .

# Run server only (skip migrations)
server:
	go run . server

worker:
	go run . worker

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

# Generate protobuf code for all modules
proto:
	@bash scripts/generate-module-protobuf.sh

# Generate protobuf code for specific module (e.g., make proto-product)
proto-%:
	@bash scripts/generate-module-protobuf.sh $*

# Check module dependencies
deps-check:
	go run cmd/lint-deps/main.go

# Database migrations
migrate-up:
	go run . migration sql up

migrate-down:
	go run . migration sql down

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

# All checks before commit (run proto generation manually if .proto files changed)
pre-commit: deps-check lint test