The shell splits unquoted words at spaces. Double quotes keep a literal path together while still allowing variable expansion.

A path that starts with `/` begins at the filesystem root. A relative path begins in your current working directory. Run `pwd` before moving files in an unfamiliar project. Use `ls -a` to include names that start with a dot.

`mv SOURCE DESTINATION` changes the original file's location. Check both directories afterward so you can tell a move from a copy.

## Commands and operations

Run `pwd` to locate yourself, then `ls -a "Draft Files"` and `ls -a config` to inspect the two directories. The quotes pass `Draft Files` as one argument. Without them, the shell splits that name into two words before `ls` or `mv` sees it. The `-a` flag includes hidden entries such as `.keep`.

- `mv "Input Files/db settings.txt" settings/db.conf` takes two arguments. The first is the existing file and the second is its new path. Giving a destination directory alone keeps the old filename. Giving the full new path moves and renames the file in one operation.
- `cat settings/db.conf` checks the file contents after the move. Listing the old directory separately tells you whether the original was moved or copied.
- A plain `mv` may replace a destination that already exists. Inspect both paths before running it. For an interactive confirmation before replacement, add `-i`.

A missing-file error usually means the working directory or filename is different than you thought. Run `pwd` and list the quoted directory again. Move the requested file itself so unrelated files, including hidden markers, stay in place.

## Create and inspect

```bash
touch notes.txt
printf 'hello\n' > notes.txt
cat notes.txt
stat notes.txt
```

`stat` exposes metadata such as size, timestamps, ownership, and mode.

## Copy and move

```bash
cp notes.txt backup.txt
cp -a project/ project.backup/
mv notes.txt notes.md
```

`cp -a` preserves metadata where supported and is useful for directory copies.

## Delete

```bash
rm notes.md
rm -r old-directory
```

Be careful with recursive deletion. Prefer explicit paths and verify before using `-r` or `-f`.

## Hidden files

A filename beginning with `.` is hidden by convention:

```bash
ls -la
```

It is not a security mechanism.

## Links

Hard link:

```bash
ln original.txt second-name.txt
```

Symbolic link:

```bash
ln -s /opt/app/config.yaml config.yaml
```

A hard link references the same inode. A symbolic link stores a path to another object.

Inspect:

```bash
ls -li
readlink config.yaml
```

## Safe handling of leading dashes

A filename like `-report.txt` can be interpreted as an option.

Use:

```bash
rm -- -report.txt
```

or an explicit path:

```bash
rm ./-report.txt
```

## Batch operations

Quote shell variables:

```bash
for file in ./*.txt; do
  [ -e "$file" ] || continue
  printf '%s\n' "$file"
done
```

For robust processing of arbitrary filenames, understand NUL-delimited workflows:

```bash
find . -type f -print0
```
