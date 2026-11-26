A shell script sees its first argument as `$1`. Double-quote it to preserve spaces in a path. `[ -f "$path" ]` checks whether a regular file is present.

Use stdout for successful output, stderr for error messages, and an exit status for the caller. Shell automation can branch on that status without parsing an error sentence.

## Commands and operations

Invoke your script with `bash check-config.sh "fixtures/demo config"`. Bash makes the quoted path available as the first argument, `$1`. Inside the script, use `"$1"` whenever it stands for a filename so the embedded space survives both invocations.

- `[ -f "$1" ]` returns success for a regular file. Use it as the condition of an `if` block. In the success branch, `cat -- "$1"` prints the contents. The `--` prevents a filename starting with a dash from being parsed as an option.
- In the failure branch, `printf '%s\n' "missing config" >&2` writes the message to standard error. `exit 2` lets the caller distinguish that failure using its status.
- Redirect the good output with `>` and the bad output with `2>`. Save `$?` after each run, before inspecting the files.

The grader runs the script again against another input. Test a path with spaces and a missing path in your own shell, and inspect both the redirect target and the status. A script that always prints the sample output will fail on the new fixture.
