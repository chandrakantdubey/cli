# Install cli-silly on Linux

The release bundles the CLI and its offline course. You need Linux on x86-64 or arm64, Bash, and standard Linux utilities such as tar and sha256sum. `cli-silly doctor` lists any additional tools required by specific workshops.

For a version that has yet to be published, use the source installation below.

## Install a release

1. Open the [GitHub releases page](https://github.com/chandrakantdubey/cli/releases) and select a published version.
2. Download `install.sh`, the matching `SHA256SUMS-<version>.txt`, and the archive for your machine into the same directory. Run `uname -m` to identify your architecture. Use the `linux-amd64` archive for `x86_64` and `linux-arm64` for `aarch64`.
3. Run the installer from that directory. The example below uses version 0.1.0 on x86-64. Substitute the name of the archive you downloaded.

```bash
bash install.sh cli-silly-0.1.0-linux-amd64.tar.gz
"$HOME/.local/bin/cli-silly" --version
"$HOME/.local/bin/cli-silly" doctor
"$HOME/.local/bin/cli-silly"
```

The installer verifies the archive checksum before installing the executable into `~/.local/bin`. To use a different directory, pass `--bin-dir /path/to/bin`. Once your install directory is on PATH, launch the app with `cli-silly` from any directory.

## Download a release directly from GitHub

From a download directory, set the published version and your architecture. This example uses version 0.1.0 and amd64. Use `arm64` for ARM64 Linux.

```bash
version=0.1.0
arch=amd64
release="https://github.com/chandrakantdubey/cli/releases/download/v$version"
archive="cli-silly-$version-linux-$arch.tar.gz"

curl -fLO "$release/install.sh" &&
curl -fLO "$release/SHA256SUMS-$version.txt" &&
curl -fLO "$release/$archive" &&
bash install.sh "$archive"
```

## Upgrade or uninstall

Download the new archive and its matching checksum file. Pass `--replace` to install the verified binary over your current copy.

```bash
bash install.sh --replace cli-silly-0.1.1-linux-amd64.tar.gz
```

To remove the executable:

```bash
bash install.sh --uninstall
```

Pass the same `--bin-dir` if you chose one during install. Workspaces and progress stay in the CLI's XDG data and state directories so they are available when you reinstall.

## Build and install from a clone or source ZIP

Clone or download the ZIP from the [GitHub source repository](https://github.com/chandrakantdubey/cli). With Go 1.26 or newer installed, open the folder containing `go.mod` and run:

```bash
GOBIN="$HOME/.local/bin" go install ./cmd/cli-silly
~/.local/bin/cli-silly
```
