Trace what Bash does after you press Enter. This lesson gives you one command that emits ordinary output, emits a diagnostic, and reports failure. Each is a separate part of the command's interface.

## Commands and operations

A front-end tool can print useful data and still report a failure. `>` writes file descriptor 1 (stdout) to a file. `2>` writes descriptor 2 (stderr) elsewhere. `$?` expands to the exit status of the most recent foreground command or pipeline.

Run the fixture with `bash bin/sync.sh`. The two lines look similar on your terminal because Bash connects both streams to it by default. Redirect stderr first to see which line remains. Capture `$?` immediately after running the script. Even `echo` or `cat` changes that value. An exit status is a number even when the error line contains words.

In `tool | head -n 1`, only stdout flows through the pipe. Bash normally reports the last command's status for a pipeline. Try `(set -o pipefail; false | true; printf '%s\n' "$?")` to see how a temporary subshell preserves the upstream failure as status 1.

## The command lifecycle

Suppose you enter:

```bash
ls -la /tmp
```

The shell:

1. reads your input;
2. performs parsing and expansions;
3. resolves the command;
4. creates or manages processes;
5. connects standard input/output/error;
6. starts the program;
7. waits for it or continues asynchronously.

The terminal is not the program doing the work.

## Commands are programs, shell built-ins, or functions

```bash
type cd
type ls
type printf
```

`cd` is normally a shell built-in because changing the current directory must change the shell's own state. `ls` is normally an external executable.

```bash
command -v ls
command -v bash
```

## Standard streams

Every process starts with conventional file descriptors:

```text
0 stdin
1 stdout
2 stderr
```

This is why these operations work:

```bash
command > output.txt
command 2> errors.txt
command > output.txt 2>&1
```

## Exit status

Commands return an integer status. By convention, `0` means success and non-zero means failure.

```bash
true
echo $?

false
echo $?
```

The exact non-zero value is program-specific.

Use status to compose commands:

```bash
mkdir build && cd build
```

The `cd` runs only if `mkdir` succeeds.

## Pipelines

A pipeline connects stdout of one process to stdin of another:

```bash
ps aux | grep nginx
```

The shell creates the connections; the individual programs still do their own work.

## Environment

A process receives environment variables:

```bash
printf '%s\n' "$HOME"
printf '%s\n' "$PATH"
env | sort
```

A child process normally inherits a copy of its parent's environment.

## Command lookup

When you type:

```bash
ls
```

the shell searches directories in `$PATH` for an executable matching the command name.

```bash
printf '%s\n' "$PATH" | tr ':' '\n'
command -v ls
```

## Before you practice

**Outcome:** Distinguish command output from the success or failure reported to the shell.

First practice pwd, cat, and writing stdout to a file.

## Common mistakes and recovery

Another command replaces $?. Capture the status immediately. A program may print useful output and still fail.

Inspect the relevant input and output before repeating a command. Use `cli-silly hint` for a starting point and `cli-silly check` for criterion-by-criterion feedback. If the starter is damaged, `cli-silly reset` archives this attempt before recreating it.

## Use it in development

Build scripts must propagate failure even when the logs look plausible.

After passing, explain which files or state changed and how you verified them. Open related practice from **Do** to apply the idea to different data.
