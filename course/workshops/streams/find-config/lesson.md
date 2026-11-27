Use `grep -r` to search files recursively. Add `-l` when you want only filenames. Anchor a pattern with `^` and `$` to match a whole setting line. Restrict the search to the directory and file extension you care about.

Shell redirection with `>` saves stdout to a file. Use `cat` to inspect the saved result.

## Commands and operations

Use `grep -rl '^debug=true$' config` to scan beneath `config`. The `-r` flag searches recursively, and `-l` prints only the names of files that contain a match. The pattern anchors `^` and `$` require the entire line to match. This excludes comments and similar settings with extra text.

- Remove `-l` and give `grep` one known filename when you want the matching line instead of the filename.
- Run the search in the terminal first. Then use `>` to send its standard output to the requested report file. A separate command is useful here because the path and setting reports contain different things.
- Inspect both reports with `cat`. The first should be a relative path under `config/`. The second should contain the setting itself on one line.

Search the named directory. An example file outside `config` can contain the same value and still be irrelevant. If nothing matches, inspect the file and pattern before weakening the anchors. `grep` returns status 1 for no matches, which is different from a search error.

## find

```bash
find . -type f -name '*.log'
find /var/log -type f -mtime -1
find . -type f -size +100M
```

Useful predicates include `-type`, `-name`, `-path`, `-mtime`, `-size`, `-user`, and `-perm`.

In this lab, run `find config -type f -name '*.conf'` to see the config files before searching their contents. The quotes hand the pattern to `find` so Bash leaves it intact. `find` answers which paths match, while `grep` answers which lines inside those paths match.

## Execute safely

```bash
find . -type f -name '*.tmp' -print
```

Inspect results before adding actions such as deletion.

For arbitrary filenames, prefer NUL separation:

```bash
find . -type f -print0
```

Then consume with a NUL-aware tool:

```bash
find . -type f -print0 | xargs -0 wc -l
```

## Locate commands

```bash
command -v python
type -a python
```

`which` is less reliable as a teaching primitive because shell aliases/functions and built-ins can complicate command lookup.

## Search file contents

```bash
grep -R --line-number 'timeout' .
rg 'timeout' .
```

`rg` (ripgrep) is often faster and more ergonomic for source trees.

## locate

If installed:

```bash
locate nginx.conf
```

`locate` uses an index, so results can be stale. `find` searches the filesystem directly.

## File metadata

```bash
stat file.txt
file file.txt
```

`file` examines content/signatures; extensions are not authoritative.

## Finding processes

```bash
pgrep -af nginx
pidof nginx
```

These are often cleaner than parsing `ps` output.
