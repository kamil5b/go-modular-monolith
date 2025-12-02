#!/usr/bin/env bash
set -euo pipefail

# Generate mocks using mockgen -source mode for each interface found under internal/
echo "Using mockgen from: $(go env GOPATH)/bin/mockgen"
PATH="$(go env GOPATH)/bin:$PATH"

count=0
processed_files=0
while IFS= read -r -d $'\0' file; do
  # detect if file contains any interface declarations
  iface_count=$(awk '/^[[:space:]]*type[[:space:]]+[A-Za-z_][A-Za-z0-9_]*[[:space:]]+interface/ {count++} END{print count+0}' "$file")
  if [[ "$iface_count" -eq 0 ]]; then
    continue
  fi
  out_dir="$(dirname "$file")/mocks"
  mkdir -p "$out_dir"
  base=$(basename "$file" .go)
  out_file="$out_dir/mock_${base}.go"
  echo "Generating mocks from $file -> $out_file (interfaces: $iface_count)"
  mockgen -source="$file" -destination="$out_file" -package="mocks"
  ((count++))
  ((processed_files++))
done < <(find internal -type f -name '*.go' ! -name '*_test.go' -print0)

echo "Generated $count mock files from $processed_files source files."
