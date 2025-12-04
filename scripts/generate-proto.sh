#!/bin/bash
# Generate protobuf code for all modules

set -e

REPO_ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
MODULES_DIR="$REPO_ROOT/internal/modules"

# Check if protoc is installed
if ! command -v protoc &> /dev/null; then
    echo "Error: protoc is not installed. Please install it first:"
    echo "  macOS: brew install protobuf"
    echo "  Linux: apt-get install protobuf-compiler"
    exit 1
fi

# Check if protoc-gen-go is installed
if ! command -v protoc-gen-go &> /dev/null; then
    echo "Installing protoc-gen-go..."
    go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
fi

# Check if protoc-gen-go-grpc is installed
if ! command -v protoc-gen-go-grpc &> /dev/null; then
    echo "Installing protoc-gen-go-grpc..."
    go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
fi

echo "Scanning for proto files in modules..."
echo ""

GENERATED_COUNT=0

# Find all modules with proto definitions
for MODULE_DIR in "$MODULES_DIR"/*; do
    if [ ! -d "$MODULE_DIR" ]; then
        continue
    fi
    
    MODULE_NAME=$(basename "$MODULE_DIR")
    PROTO_DIR="$MODULE_DIR/domain/proto"
    
    # Skip if module doesn't have proto definitions
    if [ ! -d "$PROTO_DIR" ]; then
        continue
    fi
    
    echo "📦 Processing module: $MODULE_NAME"
    
    # Find all .proto files in this module
    PROTO_FILES=$(find "$PROTO_DIR" -name "*.proto" 2>/dev/null)
    
    if [ -z "$PROTO_FILES" ]; then
        echo "   ⚠️  No .proto files found"
        continue
    fi
    
    # Generate code for each proto file
    for PROTO_FILE in $PROTO_FILES; do
        echo "   🔨 Generating: $(basename "$PROTO_FILE")"
        
        protoc \
            --go_out="$REPO_ROOT" \
            --go-grpc_out="$REPO_ROOT" \
            --go_opt=module=github.com/kamil5b/go-pste-monolith \
            --go-grpc_opt=module=github.com/kamil5b/go-pste-monolith \
            -I"$PROTO_DIR" \
            "$PROTO_FILE"
        
        GENERATED_COUNT=$((GENERATED_COUNT + 1))
    done
    
    echo "   ✓ Generated code: internal/modules/$MODULE_NAME/proto/"
    echo ""
done

if [ $GENERATED_COUNT -eq 0 ]; then
    echo "⚠️  No proto files found in any module"
    echo ""
    echo "To add proto definitions, create:"
    echo "  internal/modules/<module>/domain/proto/v1/<module>.proto"
    exit 0
fi

echo "✅ Successfully generated protobuf code for $GENERATED_COUNT proto file(s)!"
echo ""
echo "Generated files are located in:"
echo "  internal/modules/<module>/proto/*.pb.go"
