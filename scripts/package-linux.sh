#!/usr/bin/env bash
set -euo pipefail

version=${1:?Usage: bash scripts/package-linux.sh VERSION [OUTPUT_DIRECTORY]}
if [[ ! $version =~ ^[[:alnum:]._-]+$ ]]; then
  printf 'Invalid release version: %s\n' "$version" >&2
  exit 2
fi

cd "$(dirname "$0")/.."
output=${2:-artifacts/releases}
mkdir -p "$output"
output=$(realpath "$output")
cp scripts/install-linux.sh "$output/install.sh"
stage=$(mktemp -d)
trap 'rm -r -- "$stage"' EXIT

for arch in amd64 arm64; do
  bundle="cli-silly-${version}-linux-${arch}"
  mkdir "$stage/$bundle"
  GOOS=linux GOARCH=$arch CGO_ENABLED=0 "${GO:-go}" build -trimpath \
    -ldflags "-X main.version=$version" -o "$stage/$bundle/cli-silly" ./cmd/cli-silly
  cp README.md "$stage/$bundle/"
  cp docs/install.md "$stage/$bundle/INSTALL.md"
  tar -czf "$output/$bundle.tar.gz" -C "$stage/$bundle" cli-silly README.md INSTALL.md
done

(
  cd "$output"
  sha256sum "cli-silly-$version-linux-amd64.tar.gz" "cli-silly-$version-linux-arm64.tar.gz" > "SHA256SUMS-$version.txt"
)
