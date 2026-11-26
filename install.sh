#!/usr/bin/env bash

main() (
  set -euo pipefail

  fail() {
    printf 'cli-silly installer: %s\n' "$*" >&2
    exit 1
  }

  repo=https://github.com/chandrakantdubey/cli
  tag=
  options=()
  while (($#)); do
    case "$1" in
      --help|-h)
        printf '%s\n' 'Install the latest cli-silly Linux release.' \
          'Options: --version vVERSION, --bin-dir DIRECTORY, --replace'
        return
        ;;
      --version)
        (($# >= 2)) && [[ -n $2 ]] || fail '--version needs a published tag'
        tag=$2
        shift
        ;;
      --bin-dir)
        (($# >= 2)) && [[ -n $2 ]] || fail '--bin-dir needs a directory'
        options+=("$1" "$2")
        shift
        ;;
      --replace) options+=("$1") ;;
      *) fail "unknown option $1" ;;
    esac
    shift
  done

  [[ $(uname -s) == Linux ]] || fail 'this installer supports Linux'
  case "$(uname -m)" in
    x86_64) arch=amd64 ;;
    aarch64|arm64) arch=arm64 ;;
    *) fail 'this architecture has no Linux release archive' ;;
  esac

  command -v curl >/dev/null || fail 'curl is required'
  curl_options=(--fail --silent --show-error --location --proto '=https' --proto-redir '=https' --connect-timeout 15 --retry 2)
  if [[ -z $tag ]]; then
    url=$(curl "${curl_options[@]}" --max-time 30 --output /dev/null --write-out '%{url_effective}' "$repo/releases/latest")
    [[ $url == "$repo/releases/tag/"* ]] || fail 'could not find a published release'
    tag=${url##*/}
  fi
  [[ $tag =~ ^v[0-9]+\.[0-9]+\.[0-9]+(-[A-Za-z0-9.-]+)?$ ]] || fail 'expected a published tag such as v0.1.1'

  version=${tag#v}
  base="$repo/releases/download/$tag"
  archive="cli-silly-$version-linux-$arch.tar.gz"
  scratch=$(mktemp -d)
  trap 'rm -f -- "$scratch/install.sh" "$scratch/SHA256SUMS-$version.txt" "$scratch/$archive"; rmdir -- "$scratch"' EXIT
  printf 'Downloading cli-silly %s for Linux %s…\n' "$tag" "$arch"
  for file in "SHA256SUMS-$version.txt" "$archive" install.sh; do
    curl "${curl_options[@]}" --max-time 240 --output "$scratch/$file" "$base/$file"
  done
  bash "$scratch/install.sh" "${options[@]}" "$scratch/$archive"
)

main "$@"
