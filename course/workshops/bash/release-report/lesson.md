Read a list with `while IFS= read -r name` to preserve spaces and backslashes. Quote expansions that form paths. Sort your final report when its order should be independent of the input. Redirect the command group that writes the report so a second run replaces the old report.

## Commands and operations

The pattern `while IFS= read -r service; do ...; done < services.txt` runs the loop once per input line. An empty `IFS` keeps leading and trailing whitespace, and `-r` tells `read` to keep backslashes as written. The redirect after `done` gives the entire loop its input file.

- Test each service with `[ -f "config/$service.conf" ]`. Quoting the complete path keeps a name like `background worker` together. Write one line in each branch using the same output shape.
- `printf 'present %s\n' "$service"` uses a fixed format. The service name is a separate value, so percent signs in it would remain ordinary text.
- Redirect the loop's output with `> report.txt` so every run replaces the report. An append redirect, `>>`, will duplicate the rows when run again.

Inspect the result with `cat report.txt`. Keep the same order as the source list and verify it still has one line per service after a second run. The grader runs your script again with a different list and configuration files.
