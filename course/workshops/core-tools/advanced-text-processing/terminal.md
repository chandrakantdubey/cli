## Try it in the terminal

### Inspect the inputs to a comparison

Use `jq '.services[] | [.name, .enabled]' deploy.json` to see the data shape. Then inspect what `awk '$2 == "ERROR" {print $1}' events.log` emits before sorting. Find the repeated name and the service whose error should be excluded because the service is disabled.

Try `comm` on two scratch lists sorted with the same locale. Without column switches, its output shows left-only, right-only, and shared values. Compare that view with `comm -12`. Save and inspect each intermediate list before combining, especially if a final report is unexpectedly empty.
