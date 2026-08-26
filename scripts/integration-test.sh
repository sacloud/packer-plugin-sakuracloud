#!/usr/bin/env bash
set -euo pipefail

OUTPUT_FILE="integration-test-$(date +%Y%m%d-%H%M%S).log"

# Move to repository root
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
cd "$SCRIPT_DIR/.."

# Output both stdout/stderr to terminal and log file
exec > >(tee -a "$OUTPUT_FILE")
exec 2>&1

run() {
    echo "==> RUN: $*"
    "$@"
}

echo "==> Integration test started at $(date '+%Y-%m-%d %H:%M:%S')"

# Check required environment variables
if [[ -z "${SAKURA_ACCESS_TOKEN:-}" ]] || [[ -z "${SAKURA_ACCESS_TOKEN_SECRET:-}" ]]; then
    echo "Error: SAKURA_ACCESS_TOKEN and SAKURA_ACCESS_TOKEN_SECRET must be set." >&2
    exit 1
fi

echo "==> Step 1: Confirm archive list is empty before test"
before=$(usacloud archive list -o table --scope user --zone is1b -o json)
echo "$before"
if [[ -n "$before" ]] && [[ "$before" != "[]" ]]; then
    echo "Error: Expected empty array before test, but got: $before" >&2
    exit 1
fi

echo "==> Step 2: Build and install plugin"
run make tools
run make
run make install-plugin

echo "==> Step 3: Verify packer plugin installation"
run packer plugins installed
while IFS= read -r plugin_path; do
    [[ -n "$plugin_path" ]] || continue
    run ls -lah "$plugin_path"
done < <(packer plugins installed)

echo "==> Step 4: Validate and build packer template"
run packer validate examples/ubuntu/template.pkr.hcl
run packer build examples/ubuntu/template.pkr.hcl

echo "==> Step 5: Confirm archive list is not empty after test"
after=$(usacloud archive list -o table --scope user --zone is1b -o json)
echo "$after"
if [[ -z "$after" ]] || [[ "$after" == "[]" ]]; then
    echo "Error: Expected non-empty array after test, but got empty array." >&2
    exit 1
fi

echo "==> Step 6: Cleanup created archive"
run usacloud archive delete --zone is1b -y packer-example-ubuntu

echo "==> Integration test completed successfully at $(date '+%Y-%m-%d %H:%M:%S')"
echo "==> Log saved to: $OUTPUT_FILE"
