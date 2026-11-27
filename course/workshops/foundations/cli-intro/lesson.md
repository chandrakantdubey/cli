Use your first Bash session to recognize the prompt, list a project, and read a file. The terminal displays text. Bash interprets the command you type and uses your current directory to resolve relative paths.

## Commands and operations

Run `pwd` first. It prints the full path of your current directory. In a lab it ends in the lesson folder. `ls` lists names in that directory, and `ls docs` lists the contents of the relative path `docs`. Names are case-sensitive on the filesystems used by most Linux labs.

`cat docs/welcome.txt` prints the file contents. Add `> welcome-copy.txt` to send that same output to a file. A command that redirects stdout may print nothing in the terminal. Verify its result by reading the destination with `cat`. Running the command again replaces the old report.

With GNU `ls`, use `ls -1 docs` for one name per line. The flag is the digit **one**. Set `LC_ALL=C` for this command if you need a predictable A–Z sort. Never type the example's leading prompt character.

## Why learn the CLI?

- **Automation:** repeat reliable work with scripts and pipelines.
- **Remote operations:** administer servers over SSH without a graphical desktop.
- **Precision:** expose options and behavior that GUIs may hide.
- **Composability:** combine small programs into useful workflows.
- **Debugging:** inspect processes, files, sockets, logs, and system state directly.

## Linux: kernel, distribution, and user space

Linux technically refers to the kernel. A **distribution** combines the kernel with user-space programs, package management, libraries, defaults, and release tooling.

Examples include Debian, Ubuntu, Fedora, Arch, and Alpine.

Android also uses the Linux kernel, but it is not simply a conventional Linux desktop distribution.

## Terminal vs shell

A **terminal emulator** provides the text interface. A **shell** interprets commands.

Common shells:

```text
bash
zsh
fish
sh
```

You can change shells without changing the terminal application.

```text
keys    → terminal → Bash → command
output  ← terminal ← Bash ← command
                                 ↕
                            Linux kernel
```

## Your first commands

```bash
printf '%s\n' 'Hello, Linux'
pwd
whoami
uname -srm
```

`pwd` prints the current working directory. `whoami` prints the effective username. `uname` reports kernel/system information.

## Case sensitivity and paths

Linux filenames are normally case-sensitive:

```text
notes.txt
Notes.txt
NOTES.txt
```

These can be three different names.

A path containing spaces should be quoted:

```bash
cd "Project Files"
```

Avoiding spaces is a preference, not a Linux requirement.

## Safe habits from day one

- Read a command's help before using an unfamiliar destructive option.
- Treat `sudo`, `rm`, `mv`, `chmod`, `chown`, and disk tools with care.
- Use a disposable directory for experiments.
- Never paste real credentials into examples or shell history.
