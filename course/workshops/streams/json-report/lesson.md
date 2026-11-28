JSON has objects, arrays, strings, numbers, and booleans. Use `jq` to parse those types. The expression `.services[]` visits each element in the services array. Pipe those objects into `select` to keep a matching subset. The `-r` option prints a string without JSON quotes.

## Commands and operations

Start by reading the structure with `jq . deploy.json`. The dot returns the whole document, so you can see whether the data is an array or an object containing an array.

- `.services[]` emits each service object individually. Apply `select(.enabled)` while you still have that whole object, then read `.name` to emit only its name. If you extract the name first, you lose the Boolean needed for filtering.
- For the tab-separated report, construct a new array with `[.name, .port]` for each selected service and pass it through `@tsv`. The filter formats the two values as one row.
- Use `jq -r` for both reports. It emits raw strings with no surrounding JSON quotes, so the TSV output contains actual tabs rather than the JSON escape sequence for a tab.

Inspect each filter in the terminal before redirecting to a file. The order should match the source array, and disabled entries should be absent from both reports. If the output is `null`, check the key spelling and the point in the filter where you changed the input's shape.

## Before you practice

**Outcome:** Select JSON objects and produce both plain-text and TSV reports.

Know redirection and pipelines. jq is required.

## Common mistakes and recovery

JSON booleans are not strings. Preserve array order for this task and use @tsv for tab-separated fields.

Inspect the relevant input and output before repeating a command. Use `cli-silly hint` for a starting point and `cli-silly check` for criterion-by-criterion feedback. If the starter is damaged, `cli-silly reset` archives this attempt before recreating it.

## Use it in development

Transform a service inventory into inputs for deployment tooling.

After passing, explain which files or state changed and how you verified them. Open related practice from **Do** to apply the idea to different data.
