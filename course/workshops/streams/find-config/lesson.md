Use `grep -r` to search files recursively. Add `-l` when you want only filenames. Anchor a pattern with `^` and `$` to match a whole setting line. Restrict the search to the directory and file extension you care about.

Shell redirection with `>` saves stdout to a file. Use `cat` to inspect the saved result.

## Commands and operations

Use `grep -rl '^debug=true$' config` to scan beneath `config`. The `-r` flag searches recursively, and `-l` prints only the names of files that contain a match. The pattern anchors `^` and `$` require the entire line to match. This excludes comments and similar settings with extra text.

- Remove `-l` and give `grep` one known filename when you want the matching line instead of the filename.
- Run the search in the terminal first. Then use `>` to send its standard output to the requested report file. A separate command is useful here because the path and setting reports contain different things.
- Inspect both reports with `cat`. The first should be a relative path under `config/`. The second should contain the setting itself on one line.

Search the named directory. An example file outside `config` can contain the same value and still be irrelevant. If nothing matches, inspect the file and pattern before weakening the anchors. `grep` returns status 1 for no matches, which is different from a search error.
