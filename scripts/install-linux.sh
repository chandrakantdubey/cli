#!/usr/bin/env bash
set -euo pipefail

usage() {
  cat <<'HELP'
Install:   bash install.sh [--bin-dir DIRECTORY] [--replace] ARCHIVE
Uninstall: bash install.sh [--bin-dir DIRECTORY] --uninstall

Place the release's SHA256SUMS file next to ARCHIVE. The default install
directory is ~/.local/bin. Uninstall removes the executable and keeps labs
and scores. Use --replace to upgrade an existing installation.
HELP
}

fail() {
  printf 'cli-silly installer: %s\n' "$*" >&2
  exit 1
}

bin_dir="$HOME/.local/bin"
archive=
uninstall=false
replace=false
while (($#)); do
  case "$1" in
    --help|-h) usage; exit 0 ;;
    --bin-dir)
      (($# >= 2)) && [[ -n $2 ]] || fail '--bin-dir needs a directory'
      bin_dir=$2
      shift
      ;;
    --replace) replace=true ;;
    --uninstall) uninstall=true ;;
    -*) fail "unknown option $1" ;;
    *)
      [[ -z $archive ]] || fail 'specify one release archive'
      archive=$1
      ;;
  esac
  shift
done

destination="$bin_dir/cli-silly"
if [[ $uninstall == true ]]; then
  [[ -z $archive && $replace == false ]] || fail 'use --uninstall with only an optional --bin-dir'
  [[ ! -L $destination ]] || fail "remove the symlink with the tool that installed it: $destination"
  if [[ ! -e $destination ]]; then
    printf 'No cli-silly executable at %s\n' "$destination"
    exit 0
  fi
  [[ -f $destination ]] || fail "expected a regular file at $destination"
  rm -- "$destination"
  printf 'Removed %s. Your labs and scores are still available.\n' "$destination"
  exit 0
fi

[[ -n $archive ]] || { usage >&2; exit 2; }
[[ $(uname -s) == Linux ]] || fail 'this installer supports Linux'
case "$(uname -m)" in
  x86_64) arch=amd64 ;;
  aarch64|arm64) arch=arm64 ;;
  *) fail 'this architecture has no Linux release archive' ;;
esac

filename=${archive##*/}
if [[ ! $filename =~ ^cli-silly-([[:alnum:]._-]+)-linux-(amd64|arm64)\.tar\.gz$ ]]; then
  fail 'expected cli-silly-VERSION-linux-ARCH.tar.gz'
fi
version=${BASH_REMATCH[1]}
[[ ${BASH_REMATCH[2]} == "$arch" ]] || fail "choose the linux-$arch archive for this machine"
[[ -f $archive ]] || fail "archive not found: $archive"
sums="$(dirname -- "$archive")/SHA256SUMS-$version.txt"
[[ -f $sums ]] || fail "put SHA256SUMS-$version.txt next to the archive"

expected=$(awk -v file="$filename" '$2 == file {print $1}' "$sums")
[[ $expected =~ ^[[:xdigit:]]{64}$ ]] || fail "expected one checksum for $filename"
actual=$(sha256sum < "$archive")
[[ ${actual%% *} == "$expected" ]] || fail "checksum mismatch for $filename"

mkdir -p -- "$bin_dir"
bin_dir=$(cd -- "$bin_dir" && pwd -P)
destination="$bin_dir/cli-silly"
[[ ! -L $destination && ( ! -e $destination || -f $destination ) ]] || fail "expected a regular file at $destination"

stage=$(mktemp "$bin_dir/.cli-silly-install.XXXXXX")
trap 'rm -f -- "$stage"' EXIT
tar -xOzf "$archive" -- cli-silly > "$stage" || fail 'archive does not contain a readable cli-silly binary'
chmod 755 "$stage"
[[ $("$stage" --version) == "cli-silly $version" ]] || fail 'archive binary failed its version check'

if [[ $replace == true ]]; then
  mv -fT -- "$stage" "$destination"
elif ! ln -T -- "$stage" "$destination" 2>/dev/null; then
  [[ -f $destination && -x $destination && ! -L $destination ]] && cmp -s "$stage" "$destination" ||
    fail "could not install at $destination. Use --replace to upgrade an existing cli-silly executable"
fi

printf 'Installed cli-silly %s at %s\n' "$version" "$destination"
if [[ $(command -v cli-silly || true) != "$destination" ]]; then
  printf 'Run it directly with: %q\n' "$destination"
  printf 'To use cli-silly by name, put this directory first on your PATH: %s\n' "$bin_dir"
fi
