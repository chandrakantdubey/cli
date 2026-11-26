JSON has objects, arrays, strings, numbers, and booleans. Use `jq` to parse those types. The expression `.services[]` visits each element in the services array. Pipe those objects into `select` to keep a matching subset. The `-r` option prints a string without JSON quotes.

## Commands and operations

Start by reading the structure with `jq . deploy.json`. The dot returns the whole document, so you can see whether the data is an array or an object containing an array.

- `.services[]` emits each service object individually. Apply `select(.enabled)` while you still have that whole object, then read `.name` to emit only its name. If you extract the name first, you lose the Boolean needed for filtering.
- For the tab-separated report, construct a new array with `[.name, .port]` for each selected service and pass it through `@tsv`. The filter formats the two values as one row.
- Use `jq -r` for both reports. It emits raw strings with no surrounding JSON quotes, so the TSV output contains actual tabs rather than the JSON escape sequence for a tab.

Inspect each filter in the terminal before redirecting to a file. The order should match the source array, and disabled entries should be absent from both reports. If the output is `null`, check the key spelling and the point in the filter where you changed the input's shape.
