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
cli-silly read foundations/cli-intro | grep -F 'With GNU `ls`' > /dev/null
cli-silly read core-tools/node-cli | grep -F 'process.argv.slice(2)' > /dev/null
TERM=dumb cli-silly do foundations/cli-intro | grep -F 'welcome-copy.txt' > /dev/null
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
printf '#!/usr/bin/env bash\nprintf "cli-silly old\\n"\n' > "$prefix/old"
chmod 755 "$prefix/old"
mv -fT -- "$prefix/old" "$installed"
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

cli-silly course install "$output/markdown-basics-v1.tar.gz"
cli-silly course list | grep -F markdown-basics
cli-silly course use markdown-basics
TERM=dumb cli-silly continue | grep -F 'Write a project heading'
pack_workspace=$(cli-silly path markdown-basics/project-heading)
printf '# Team Notes\n' > "$pack_workspace/README.md"
(cd "$pack_workspace" && cli-silly check)
cli-silly course use builtin
[[ $(cli-silly progress shell/make-it-runnable) == "$score" ]]

bash "$installer" --bin-dir "$prefix" --uninstall
printf 'Verified install, upgrade, checksum rejection, uninstall, and reinstall with saved progress.\n'
