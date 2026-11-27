A pipe connects the first command's stdout to the next command's stdin. You can check a pipeline one stage at a time.

The access log has three whitespace-separated fields. The route is field 2 and the HTTP status is field 3. Select records using the status field. A route can contain the text 500 while its response status is 200.

`uniq -c` counts adjacent equal lines. Sort the routes before counting them.

## Commands and operations

The expression `awk '$3 >= 500 && $3 < 600 { print $2 }' sample.log` reads the log one line at a time. By default, `awk` splits each line on whitespace. The condition tests status field `$3`, and the action prints only route field `$2` for matching lines. Single quotes let `awk` see `$2` and `$3` without Bash expanding them first.

- Run the filter on its own and inspect the route list. It preserves input order and keeps repeats, one line for each failed request.
- `sort routes.txt | uniq -c` groups equal routes and counts each group. `uniq` only combines adjacent matches, so the order of the commands matters.
- Add `| awk '{print $1, $2}'` when the deliverable requires a count and route separated by one space. It removes the padding that `uniq -c` adds for display.

Check one route manually against the original log. If the counts are too high, compare the status field with the route field. A successful request to a URL containing `500` belongs outside the error report.

## cat, less, head, tail

```bash
cat app.log
less app.log
head -n 20 app.log
tail -n 50 app.log
tail -f app.log
```

`less` is often preferable to `cat` for large files.

## grep

```bash
grep 'ERROR' app.log
grep -n 'ERROR' app.log
grep -R --line-number 'TODO' src/
grep -E 'ERROR|WARN' app.log
```

For fixed strings, avoid regex complexity:

```bash
grep -F 'literal[brackets]' file.txt
```

## cut and sort

```bash
cut -d: -f1 /etc/passwd
sort names.txt
sort -u names.txt
```

## uniq

`uniq` removes adjacent duplicate lines, so sort first when you need global duplicate removal:

```bash
sort names.txt | uniq -c
```

## tr

```bash
printf '%s\n' 'hello world' | tr '[:lower:]' '[:upper:]'
```

`tr` translates characters; it is not a general-purpose substring replacement tool.

## sed

```bash
sed -n '1,20p' file.txt
sed 's/old/new/g' file.txt
```

By default, `sed` prints transformed output; it does not modify the original file unless you use an in-place option.

## awk

```bash
awk '{print $1}' access.log
awk -F: '{print $1, $7}' /etc/passwd
```

`awk` is excellent for record/field-oriented text processing.

## jq

For JSON, prefer a JSON parser:

```bash
jq '.items[] | .name' data.json
```

Do not parse structured JSON with fragile `grep`/`sed` pipelines when `jq` is available.
