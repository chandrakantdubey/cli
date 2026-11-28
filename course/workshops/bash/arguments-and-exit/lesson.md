A shell script sees its first argument as `$1`. Double-quote it to preserve spaces in a path. `[ -f "$path" ]` checks whether a regular file is present.

Use stdout for successful output, stderr for error messages, and an exit status for the caller. Shell automation can branch on that status without parsing an error sentence.

## Commands and operations

Invoke your script with `bash check-config.sh "fixtures/demo config"`. Bash makes the quoted path available as the first argument, `$1`. Inside the script, use `"$1"` whenever it stands for a filename so the embedded space survives both invocations.

- `[ -f "$1" ]` returns success for a regular file. Use it as the condition of an `if` block. In the success branch, `cat -- "$1"` prints the contents. The `--` prevents a filename starting with a dash from being parsed as an option.
- In the failure branch, `printf '%s\n' "missing config" >&2` writes the message to standard error. `exit 2` lets the caller distinguish that failure using its status.
- Redirect the good output with `>` and the bad output with `2>`. Save `$?` after each run, before inspecting the files.

The grader runs the script again against another input. Test a path with spaces and a missing path in your own shell, and inspect both the redirect target and the status. A script that always prints the sample output will fail on the new fixture.

## A reliable starting point

```bash
#!/usr/bin/env bash
set -euo pipefail

main() {
    printf 'Hello, %s\n' "${1:-world}"
}

main "$@"
```

`set -euo pipefail` improves failure visibility but has important Bash semantics. Learn the exceptions rather than treating it as magic.

## Variables and quoting

```bash
name="Ada"
printf '%s\n' "$name"
```

Quote variables by default:

```bash
rm -- "$file"
```

Unquoted expansions can undergo word splitting and pathname expansion.

## Arguments

```bash
printf 'script=%s\n' "$0"
printf 'first=%s\n' "${1:-}"
printf 'count=%s\n' "$#"
printf 'all=%s\n' "$*"
printf 'each=%s\n' "$@"
```

When forwarding arguments, use:

```bash
some_command "$@"
```

## Conditionals

```bash
if [[ -f "$file" ]]; then
    printf 'regular file\n'
else
    printf 'not a regular file\n'
fi
```

`[[ ... ]]` is Bash syntax. POSIX `sh` uses `[ ... ]`.

## Loops

```bash
for file in ./*.log; do
    [ -e "$file" ] || continue
    printf '%s\n' "$file"
done
```

## Functions and exit status

```bash
check_config() {
    test -f /etc/myapp/config
}
```

Functions return the status of their last command unless you explicitly `return`.

## Temporary files

Avoid predictable filenames for sensitive temporary data. Use `mktemp` where available:

```bash
tmp=$(mktemp)
trap 'rm -f "$tmp"' EXIT
```

## ShellCheck

Validate scripts:

```bash
bash -n script.sh
shellcheck script.sh
```

## When not to use shell

Use Python, Go, Node.js, or another language when you need complex data structures, extensive error handling, APIs, concurrency, or large application logic.

## Before you practice

**Outcome:** Write a file-processing script with a documented success and failure contract.

Know quoting, stdout, stderr, and status codes.

## Common mistakes and recovery

A script must handle filenames with spaces and missing inputs. A printed error without a nonzero exit status can fool callers.

Inspect the relevant input and output before repeating a command. Use `cli-silly hint` for a starting point and `cli-silly check` for criterion-by-criterion feedback. If the starter is damaged, `cli-silly reset` archives this attempt before recreating it.

## Use it in development

Write a CI helper that other tools can trust to stop on invalid input.

After passing, explain which files or state changed and how you verified them. Open related practice from **Do** to apply the idea to different data.
