Combine a structured deployment inventory and a plain-text incident log. Query each in its native format and compare the outputs without losing track of order or duplicate rows.

## Commands and operations

In `jq -r '.services[] | select(.enabled) | .name' deploy.json`, the JSON pipe passes each matching object to the name selector. `-r` prints strings without JSON quotes, ready for line-based tools.

For the second report, `awk '$2 == "ERROR" {print $1}' events.log` prints a service once per matching event. `LC_ALL=C sort -u` sorts and removes duplicates. `comm -12` prints only the lines common to two inputs that are both sorted in the same locale. The first and second columns (lines exclusive to either input) are suppressed.

In Bash, `<(command)` offers a temporary readable path connected to command output. This lets `comm` compare the output of two pipelines. Intermediate files work as well and are easier to debug. Inspect both sorted inputs before trusting an empty intersection.

## Learning goals

By the end of this module, you should be able to:

- pick `grep`, `sed`, `awk`, `sort`, `uniq`, `jq`, or Python based on input structure;
- build pipelines that are readable and safe with unusual filenames;
- aggregate logs without loading entire files into memory;
- avoid locale, quoting, and delimiter traps;
- know when text processing should become a small script.

## Start with the data shape

Before writing a command, inspect representative input:

```bash
head -n 5 access.log
file access.log
wc -l access.log
```

Ask whether the input is line-oriented, field-oriented, nested, binary, compressed, or untrusted. A command that works on a tiny sample may fail when fields contain spaces, tabs, quotes, or missing values.

## awk: records, fields, and aggregation

`awk` is a small language for record-oriented data. It reads one record at a time, splits it into fields, runs pattern/action blocks, then continues streaming.

```bash
awk '{print $1}' access.log
awk -F: '{print $1, $3, $7}' /etc/passwd
awk 'NR == 1 {print "first line:", $0}' file.txt
```

Useful built-ins:

```text
NR  current input record number
NF  number of fields in the current record
$0  full current record
$1  first field
FS  input field separator
OFS output field separator
```

Aggregation example:

```bash
awk '{count[$9]++} END {for (status in count) print status, count[status]}' access.log | sort -n
```

That works for common web logs where field 9 is the HTTP status. If logs use quoted fields or JSON, switch to a parser that understands the format.

## sed: precise stream edits

Use `sed` for simple line selection and substitutions:

```bash
sed -n '1,20p' file.txt
sed -n '/ERROR/p' app.log
sed -E 's/[[:space:]]+/ /g' file.txt
```

Use `-n` when you want selective output. Prefer `sed -E` for extended regular expressions. Be cautious with `sed -i`: in-place edits are convenient, but a bad expression can damage many files quickly. Test without `-i` first.

## sort, uniq, comm, and diff

Sorting changes how later tools behave:

```bash
sort names.txt
sort -n numbers.txt
sort -k2,2 -k1,1 data.tsv
sort access.log | uniq -c | sort -nr | head
```

`comm` compares two sorted files:

```bash
comm -3 <(sort old.txt) <(sort new.txt)
```

Process substitution is Bash-specific. For portable scripts, create temporary files with `mktemp` and clean them up.

Use `diff -u` for human review and patches:

```bash
diff -u old.conf new.conf
```

## xargs and NUL-delimited input

`xargs` turns stdin into arguments. Use NUL delimiters when filenames may contain spaces, newlines, quotes, or leading dashes:

```bash
find logs -type f -name '*.log' -print0 | xargs -0 wc -l
```

Limit command size or parallelize carefully:

```bash
find logs -type f -print0 | xargs -0 -n 50 grep -H 'ERROR'
find images -type f -print0 | xargs -0 -P 4 -n 1 identify
```

Parallel execution can reorder output and increase load, so use it deliberately.

## Structured data: jq beats regex

JSON is not safely parsed with `grep` and `sed`. Use `jq`:

```bash
jq '.services[] | .name' config.json
jq -r '.users[] | select(.active == true) | .email' users.json
jq '[.items[].bytes] | add' report.json
```

For line-delimited JSON logs:

```bash
jq -r 'select(.level == "error") | [.time, .message] | @tsv' app.jsonl
```

For YAML, use a YAML-aware tool such as `yq` or a short Python/Ruby script. YAML features such as anchors, multiline strings, and nested structures make indentation-only parsing unreliable.

## Locale and reproducibility

Sorting and character classes can vary by locale. For repeatable pipelines, set `LC_ALL=C` when byte-wise ordering is acceptable:

```bash
LC_ALL=C sort huge.txt | uniq -c
```

Document this in scripts so future readers know the ordering choice was intentional.

## When to switch to a script

A pipeline is good when each stage is obvious. Move to `awk`, Python, Node, or another language when you need named variables, tests, nested parsing, complex error handling, or reusable behavior.
