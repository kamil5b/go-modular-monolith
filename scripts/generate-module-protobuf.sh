#!/bin/bash
# Generate protobuf code for a specific module or all modules
# Usage: 
#   ./scripts/generate-module-protobuf.sh [module-name]
#   ./scripts/generate-module-protobuf.sh           # Generate for all modules
#   ./scripts/generate-module-protobuf.sh product   # Generate only for product module

set -e

REPO_ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
MODULES_DIR="$REPO_ROOT/internal/modules"
MODULE_NAME="$1"

# Check if protoc is installed
if ! command -v protoc &> /dev/null; then
    echo "❌ Error: protoc is not installed. Please install it first:"
    echo "  macOS: brew install protobuf"
    echo "  Linux: apt-get install protobuf-compiler"
    exit 1
fi

# Check if protoc-gen-go is installed
if ! command -v protoc-gen-go &> /dev/null; then
    echo "📦 Installing protoc-gen-go..."
    go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
fi

# Check if protoc-gen-go-grpc is installed
if ! command -v protoc-gen-go-grpc &> /dev/null; then
    echo "📦 Installing protoc-gen-go-grpc..."
    go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
fi

# Function to generate proto for a single module
generate_module_proto() {
    local module_name=$1
    local module_dir="$MODULES_DIR/$module_name"
    local domain_proto_dir="$module_dir/domain/proto"
    local output_proto_dir="$module_dir/proto"
    
    # Check if module exists
    if [ ! -d "$module_dir" ]; then
        echo "❌ Module '$module_name' not found in $MODULES_DIR"
        return 1
    fi
    
    # Check if module has proto definitions
    if [ ! -d "$domain_proto_dir" ]; then
        echo "⚠️  Module '$module_name' has no proto definitions (expected: $domain_proto_dir)"
        return 0
    fi
    
    # Find all .proto files
    local proto_files=$(find "$domain_proto_dir" -name "*.proto" 2>/dev/null)
    
    if [ -z "$proto_files" ]; then
        echo "⚠️  No .proto files found in $domain_proto_dir"
        return 0
    fi
    
    echo "📦 Generating protobuf code for module: $module_name"
    
    # Create output directory if it doesn't exist
    mkdir -p "$output_proto_dir"
    
    local file_count=0
    
    # Generate code for each proto file
    for proto_file in $proto_files; do
        local proto_basename=$(basename "$proto_file")
        echo "   🔨 Generating: $proto_basename"
        
        # Get the relative path from domain/proto to maintain version structure
        local rel_path=$(realpath --relative-to="$domain_proto_dir" "$proto_file")
        local version_dir=$(dirname "$rel_path")
        
        # Ensure version directory exists in output
        if [ "$version_dir" != "." ]; then
            mkdir -p "$output_proto_dir/$version_dir"
        fi
        
        # Generate with custom output path to maintain version structure
        protoc \
            --go_out="$output_proto_dir" \
            --go-grpc_out="$output_proto_dir" \
            --go_opt=paths=source_relative \
            --go-grpc_opt=paths=source_relative \
            -I"$domain_proto_dir" \
            "$proto_file"
        
        file_count=$((file_count + 1))
    done
    
    echo "   ✅ Generated $file_count proto file(s) → $output_proto_dir/"
    echo ""
    
    return 0
}

echo "🚀 Protobuf Code Generator"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo ""

GENERATED_MODULES=0
FAILED_MODULES=0

if [ -n "$MODULE_NAME" ]; then
    # Generate for specific module
    if generate_module_proto "$MODULE_NAME"; then
        GENERATED_MODULES=1
    else
        FAILED_MODULES=1
    fi
else
    # Generate for all modules
    echo "Scanning modules in: $MODULES_DIR"
    echo ""
    
    for module_dir in "$MODULES_DIR"/*; do
        if [ ! -d "$module_dir" ]; then
            continue
        fi
        
        module_name=$(basename "$module_dir")
        
        if generate_module_proto "$module_name"; then
            if [ -f "$module_dir/proto/"*.pb.go ]; then
                GENERATED_MODULES=$((GENERATED_MODULES + 1))
            fi
        else
            FAILED_MODULES=$((FAILED_MODULES + 1))
        fi
    done
fi

echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"

if [ $GENERATED_MODULES -eq 0 ] && [ $FAILED_MODULES -eq 0 ]; then
    echo "ℹ️  No proto files found in any module"
    echo ""
    echo "To add proto definitions for a module:"
    echo "  1. Create: internal/modules/<module>/domain/proto/v1/<module>.proto"
    echo "  2. Run: ./scripts/generate-module-protobuf.sh <module>"
    echo ""
    exit 0
fi

if [ $FAILED_MODULES -gt 0 ]; then
    echo "⚠️  $FAILED_MODULES module(s) failed"
    exit 1
fi

echo "✅ Successfully generated protobuf code for $GENERATED_MODULES module(s)!"
echo ""
echo "Generated structure:"
echo "  internal/modules/<module>/domain/proto/v1/*.proto  ← Source"
echo "  internal/modules/<module>/proto/v1/*.pb.go         ← Generated"
