Read a list with `while IFS= read -r name` to preserve spaces and backslashes. Quote expansions that form paths. Preserve the source list order in this report. Redirect the command group that writes the report so a second run replaces the old report.

## Commands and operations

The pattern `while IFS= read -r service; do ...; done < services.txt` runs the loop once per input line. An empty `IFS` keeps leading and trailing whitespace, and `-r` tells `read` to keep backslashes as written. The redirect after `done` gives the entire loop its input file.

- Test each service with `[ -f "config/$service.conf" ]`. Quoting the complete path keeps a name like `background worker` together. Write one line in each branch using the same output shape.
- `printf 'present %s\n' "$service"` uses a fixed format. The service name is a separate value, so percent signs in it would remain ordinary text.
- Redirect the loop's output with `> report.txt` so every run replaces the report. An append redirect, `>>`, will duplicate the rows when run again.

Inspect the result with `cat report.txt`. Keep the same order as the source list and verify it still has one line per service after a second run. The grader runs your script again with a different list and configuration files.

## Before you practice

**Outcome:** Generate a repeatable report from an input list and configuration files.

Practice Bash arguments, loops, quoting, and file tests.

## Common mistakes and recovery

Appending on every run duplicates rows. Preserve input order here; sorting would change the required report.

Inspect the relevant input and output before repeating a command. Use `cli-silly hint` for a starting point and `cli-silly check` for criterion-by-criterion feedback. If the starter is damaged, `cli-silly reset` archives this attempt before recreating it.

## Use it in development

Automate a release preflight check across several services.

After passing, explain which files or state changed and how you verified them. Open related practice from **Do** to apply the idea to different data.
