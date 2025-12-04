# Product Module - Protocol Buffers

This directory contains the Protocol Buffer definitions for the Product module's gRPC API.

## Structure

```
proto/
└── v1/
    └── product.proto    # Product service API definition (v1)
```

## Generating Code

To generate Go code from these proto files:

```bash
# From project root
bash scripts/generate-proto.sh

# Or manually
protoc \
  --go_out=../../../.. \
  --go-grpc_out=../../../.. \
  --go_opt=module=github.com/kamil5b/go-pste-monolith \
  --go-grpc_opt=module=github.com/kamil5b/go-pste-monolith \
  -I. \
  v1/product.proto
```

Generated files will be placed in:
- `internal/modules/product/proto/product.pb.go` (message types)
- `internal/modules/product/proto/product_grpc.pb.go` (service interfaces)

## Prerequisites

Install required tools:

```bash
# Install protoc compiler
# macOS
brew install protobuf

# Linux
apt-get install protobuf-compiler

# Install Go plugins
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
```

## Versioning

- **v1**: Initial version of the Product API
- Future versions (v2, v3) should be added as separate directories to maintain backward compatibility

## Design Principles

1. **Domain-Driven**: Proto files are part of the domain layer
2. **Module Isolation**: Each module owns its proto definitions
3. **Backward Compatibility**: Use proper versioning when making breaking changes
4. **Generated Code Separation**: Source proto files stay in domain, generated code goes to `proto/` at module root
