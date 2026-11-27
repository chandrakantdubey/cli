A rollout inventory mixes disabled and staging targets with production endpoints. Prepare a narrow report for the release operator, then derive a second report of distinct port numbers.

## Commands and operations

`jq '.targets[]' rollout.json` emits each JSON object. Add `select(.enabled and .environment == "production")` to require both properties. In jq, the `|` between filters passes JSON values from one filter to the next, while the Bash pipe outside the quotes passes bytes to another program.

Convert an array such as `[.name, .port]` to a row with `@tsv`. The `-r` option prints the result as raw text with literal tabs instead of JSON string escapes. Use a stable string sort on the name rows.

`cut -f2` reads the second tab-separated field. With `sort -nu`, `-n` orders values numerically and `-u` keeps one value per distinct port.
