A directory tree becomes easier to reason about when you can predict how a relative path is resolved. Practice moving into a service directory and reading a sibling directory through its parent.

## Commands and operations

`cd services/api` changes just this Bash process's working directory. `pwd` shows the result. From there, `..` means `services`, and `../worker` means the worker sibling. `.` means the current directory.

Try reading `../worker/owner.txt` after changing into `services/api`. Then save that output as `../../worker-owner.txt` to place the report at the lab's root. Verify with `ls ../..` before assuming where the file landed.

Quote a path that contains whitespace. When checking a report location, remember that the shell resolves redirection paths from its current directory, the same way it resolves file operands. Use `cd ../..` from `services/api` to return to the root.

## Important directories

| Path | Typical purpose |
|---|---|
| `/` | Filesystem root |
| `/home` | User home directories |
| `/etc` | System and application configuration |
| `/var` | Variable data such as logs and caches |
| `/tmp` | Temporary files |
| `/usr` | Installed user-space programs and data |
| `/bin` | Historically essential executables; often merged into `/usr/bin` |
| `/dev` | Device nodes |
| `/proc` | Process/kernel information pseudo-filesystem |
| `/sys` | Kernel/device information pseudo-filesystem |
| `/run` | Volatile runtime state |

Exact layouts vary by distribution and installation.

## Absolute vs relative paths

```bash
pwd
cd /var/log
cd ..
cd ~
cd -
```

- `/var/log` is absolute.
- `../logs` is relative to the current directory.
- `~` expands to your home directory in an interactive shell.

## Listing

```bash
ls
ls -lah
ls -ld /var/log
```

`ls -d` can be useful when you want metadata for the directory itself rather than its contents.

## Path expansion

```bash
printf '%s\n' ./*
printf '%s\n' "$HOME"/*
```

Globbing is performed by the shell before the program receives arguments.

## Creating a workspace

```bash
mkdir -p ~/cli-practice/project/{src,docs,build}
cd ~/cli-practice/project
```

## Find your way

```bash
find . -maxdepth 2 -type f
realpath .
readlink -f some-path
```

`readlink -f` is not available identically on every Unix platform, so check platform requirements in portable scripts.

## Mount points

A directory can be a mount point:

```bash
findmnt
findmnt /home
```

The same pathname can therefore refer to data backed by a different filesystem.

## Special paths

```text
.   current directory
..  parent directory
~   home directory (shell expansion)
-   previous directory for `cd`
```
