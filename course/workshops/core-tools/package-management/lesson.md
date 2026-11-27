Package managers can answer two different questions. What version is installed? Which package owns this executable? Use captured output from Debian and Fedora so the lab works offline on either system.

## Commands and operations

On Debian-family systems, `dpkg-query -W -f='${Package} ${Version}\n' package` queries installed records. `dpkg-query -S /path/to/file` asks which installed package recorded that path. On Fedora-family systems, compare `rpm -q package` and `rpm -qf /path/to/file`. An installed version can differ from the candidate offered by a repository.

Inspect the headers in the lab captures. The first field of `captures/debian-packages.tsv` is the package name and the second is its installed version. Whitespace-separated `awk` fields work on these tab-delimited inputs because the names and versions have no spaces. Compare the name exactly so `runner-tools` cannot match `runner-tools-dev`.

The Fedora capture has an absolute file path in the first column and a package name in the second. Match the path exactly, including its leading slash. The scored task changes only report files in the lab.

## Concepts and commands

These distro-specific examples illustrate package maintenance. Use the lab captures to practice without modifying the host.

## Debian/Ubuntu

```bash
apt update
apt search ripgrep
apt install ripgrep
apt remove ripgrep
apt upgrade
```

`apt update` refreshes package metadata. It does not itself upgrade installed packages.

## Fedora/RHEL-family

```bash
dnf search ripgrep
sudo dnf install ripgrep
sudo dnf upgrade
```

Exact package names and repository configuration vary.

## Alpine

```bash
apk add ripgrep
apk upgrade
```

## Installed files

Package managers can help answer "which package owns this file?" The exact command varies by distribution.

Do not assume a package name or command is universal across distributions.

## Repositories and trust

Packages normally come from configured repositories. Repository metadata and package signatures provide important integrity guarantees.

Do not blindly install random scripts from the internet when a trusted package source is available.

## Updating safely

For production hosts:

1. understand what is changing;
2. review package changes;
3. test where appropriate;
4. maintain backups;
5. know how to roll back or restore.

## Containers

Minimal container images may use different package managers or none at all. Do not assume `apt` exists.

## Build vs install

A package manager is not always the right answer. Some developer tools are distributed as binaries, language packages, source archives, or version managers.

Understand where the software comes from and how it will be updated.
