A process normally receives three streams. File descriptor 0 is stdin, 1 is stdout, and 2 is stderr. A pipe sends stdout forward by default. Use `2>` to record diagnostic messages without mixing them into data.

The shell opens redirections left to right. Use two clear filenames when a caller needs to keep the streams separate.

## Commands and operations

Run `bash sync-check.sh > data.txt 2> errors.txt` to route both output streams in one invocation. A plain `>` redirects file descriptor 1, standard output. The `2>` redirects file descriptor 2, standard error. Each destination is opened by Bash and starts empty on a repeated run.

- `2>&1` sends standard error wherever standard output currently points. That is useful for a combined log, but this task calls for two separate reports.
- Use `printf '%s\n' "$?" > status.txt` immediately after the command to record its exit status. Running `cat`, `echo`, or another command first replaces the value in `$?`.
- Inspect the files separately after saving the status. The file descriptor numbers tell the shell where data goes. The exit status reports how the command ended. One does not change the other.

This fixture deliberately writes to both streams and exits with a nonzero value. In a shell with `set -e`, run it in an explicit error-handling construct such as `status=0; bash sync-check.sh > data.txt 2> errors.txt || status=$?` so the shell can record the failure.
