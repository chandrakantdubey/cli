A process normally receives three streams. File descriptor 0 is stdin, 1 is stdout, and 2 is stderr. A pipe sends stdout forward by default. Use `2>` to record diagnostic messages without mixing them into data.

The shell opens redirections left to right. Use two clear filenames when a caller needs to keep the streams separate.

## Commands and operations

Run `bash sync-check.sh > data.txt 2> errors.txt` to route both output streams in one invocation. A plain `>` redirects file descriptor 1, standard output. The `2>` redirects file descriptor 2, standard error. Each destination is opened by Bash and starts empty on a repeated run.

- `2>&1` sends standard error wherever standard output currently points. That is useful for a combined log, but this task calls for two separate reports.
- Use `printf '%s\n' "$?" > status.txt` immediately after the command to record its exit status. Running `cat`, `echo`, or another command first replaces the value in `$?`.
- Inspect the files separately after saving the status. The file descriptor numbers tell the shell where data goes. The exit status reports how the command ended. One does not change the other.

This fixture deliberately writes to both streams and exits with a nonzero value. In a shell with `set -e`, run it in an explicit error-handling construct such as `status=0; bash sync-check.sh > data.txt 2> errors.txt || status=$?` so the shell can record the failure.

## The three standard streams

```text
0 stdin
1 stdout
2 stderr
```

## Redirect stdout

```bash
ls -la > listing.txt
```

`>` replaces the destination file.

Append instead:

```bash
printf '%s\n' "another line" >> listing.txt
```

## Redirect stderr

```bash
ls /does-not-exist 2> errors.txt
```

Combine stdout and stderr:

```bash
command > all.log 2>&1
```

In Bash, this is also common:

```bash
command &> all.log
```

The first form is more explicit and has broader POSIX-shell relevance.

## Read stdin

```bash
wc -l < listing.txt
```

A program does not need to know the input came from a file.

## Pipes

```bash
journalctl -b | grep -i error
```

Pipeline composition is the core Unix pattern:

```text
produce -> transform -> filter -> summarize
```

## tee

`tee` writes to a file while continuing the pipeline:

```bash
command | tee output.txt
```

## Pipeline status

In Bash:

```bash
set -o pipefail
```

Without `pipefail`, a pipeline's status is normally the status of its last command, which can hide an earlier failure.

## Here documents

```bash
cat <<'EOF'
server {
    listen 8080;
}
EOF
```

Quoting the delimiter prevents parameter and command expansion inside the document.

## Here strings

Bash provides:

```bash
grep 'ERROR' <<< "$text"
```

This is convenient but not POSIX shell syntax.
