#!/usr/bin/env bash
set -euo pipefail

version=${1:?Usage: bash scripts/test-install.sh VERSION [ARTIFACTS_DIRECTORY]}
output=$(realpath "${2:-artifacts/releases}")
scratch=$(mktemp -d)
trap 'rm -r -- "$scratch"' EXIT
prefix="$scratch/local bin"
export XDG_DATA_HOME="$scratch/data" XDG_STATE_HOME="$scratch/state"

case "$(uname -m)" in
  x86_64) arch=amd64 ;;
  aarch64|arm64) arch=arm64 ;;
  *) printf 'Unsupported test architecture\n' >&2; exit 1 ;;
esac
filename="cli-silly-$version-linux-$arch.tar.gz"
archive="$output/$filename"
installer="$output/install.sh"
installed="$prefix/cli-silly"

(cd "$output" && sha256sum -c "SHA256SUMS-$version.txt")
bash "$installer" --bin-dir "$prefix" "$archive"
export PATH="$prefix:$PATH"
[[ $(command -v cli-silly) == "$installed" ]]
[[ $(cli-silly --version) == "cli-silly $version" ]]
cli-silly read streams/log-report > /dev/null
if [[ -x bin/cli-silly ]]; then
  diff <(bin/cli-silly list) <(cli-silly list)
fi
"$installed" list | grep -F 'shell/make-it-runnable' > /dev/null
bash "$installer" --bin-dir "$prefix" "$archive"
python3 "$(dirname "$0")/test-terminal.py" "$installed"

workspace=$("$installed" path shell/make-it-runnable)
touch "$workspace/learner-notes.txt"
chmod u+x "$workspace/bin/health-check"
(cd "$workspace" && ./bin/health-check > status.txt && "$installed" check)
score=$("$installed" progress shell/make-it-runnable)

# A corrupt archive must leave a working installation in place.
mkdir "$scratch/corrupt"
cp "$output/SHA256SUMS-$version.txt" "$archive" "$scratch/corrupt/"
printf 'corrupt' >> "$scratch/corrupt/$filename"
if bash "$installer" --bin-dir "$prefix" --replace "$scratch/corrupt/$filename"; then
  printf 'Installer accepted a corrupt archive\n' >&2
  exit 1
fi
[[ $("$installed" --version) == "cli-silly $version" ]]

# An existing, different executable needs an explicit upgrade.
printf '#!/usr/bin/env bash\nprintf "cli-silly old\\n"\n' > "$installed"
if bash "$installer" --bin-dir "$prefix" "$archive"; then
  printf 'Installer replaced a different executable without --replace\n' >&2
  exit 1
fi
[[ $("$installed" --version) == 'cli-silly old' ]]
bash "$installer" --bin-dir "$prefix" --replace "$archive"
[[ $("$installed" --version) == "cli-silly $version" ]]

bash "$installer" --bin-dir "$prefix" --uninstall
[[ ! -e $installed && -f $workspace/learner-notes.txt ]]
bash "$installer" --bin-dir "$prefix" "$archive"
[[ $("$installed" progress shell/make-it-runnable) == "$score" ]]
bash "$installer" --bin-dir "$prefix" --uninstall
printf 'Verified install, upgrade, checksum rejection, uninstall, and reinstall with saved progress.\n'
