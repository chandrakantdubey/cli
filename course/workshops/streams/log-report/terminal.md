## Try it in the terminal

### Build a pipeline one stage at a time

Start with `head -n 5 access.log`. Find the route and HTTP status columns. Use `awk '{print $2, $3}' access.log` to inspect those columns independent of any filter. Read one line and identify whether its status belongs to the 500 to 599 interval.

Run the filter to the terminal first, then write the routes to `errors.txt`. Inspect both `sort errors.txt` and `sort errors.txt | uniq -c`. Sorting moves equal values next to one another, and the counting format from `uniq -c` includes padding. The score checks the final report's single separator and ascending route order. Normalize the count and route as a final stage.
