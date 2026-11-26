A script is a regular file containing commands. To run it directly by path, the caller needs execute permission and the file needs a valid interpreter line, such as `#!/usr/bin/env bash`.

## Read the permission bits

In `-rw-r--r--`, the first character identifies the file type. The next three characters describe the owner's permissions. The script owner can read and write the file. A dash in place of `x` means execute permission is missing.

| Slot | Owner permission |
| --- | --- |
| `r` | Read the file |
| `w` | Edit the file |
| `x` | Execute it directly |

Inspect the example file, grant its owner execute permission, and inspect it again:

```bash
ls -l bin/show-version
chmod u+x bin/show-version
ls -l bin/show-version
```

The owner's three slots should change from `rw-` to `rwx`.

Running `bash bin/show-version` asks Bash to read the file. Running `./bin/show-version` asks Linux to execute the file directly. The task practices the second form.

## Send the output to a report

Bash opens the file after `>` and connects the command's standard output to it. The shell creates it if necessary and truncates existing contents before the command runs. Use `cat version.txt` to inspect the result. If the report is empty, check whether the script printed an error to standard error.

## Commands and operations

Read the first ten characters from `ls -l` as the file type followed by three permission groups. The command `chmod u+x bin/show-version` changes only the owner's execute bit. The `u` selects the owner and `+` adds the named permission.

- `./bin/show-version` asks the shell to start the file by its path. The kernel uses the file's first line, the shebang, to select an interpreter. `bash bin/show-version` instead asks Bash to read the file as input and can succeed before the execute bit is set.
- `./bin/show-version > version.txt` saves standard output. The shell opens and truncates the destination before starting the script. Errors written to standard error still appear in the terminal.
- `cat version.txt` reads back the result. An empty report can mean that the shell created the file successfully but the script failed before it printed anything.

Start from the lab's project directory so relative paths such as `bin/health-check` and the configuration paths it reads resolve correctly. Use `pwd` to check the current directory and `ls -l bin/health-check` to confirm the `x` appears in the owner's group. If direct execution gives `Permission denied`, inspect the file mode before changing permissions again.
