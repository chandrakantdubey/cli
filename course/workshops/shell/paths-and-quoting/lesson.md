The shell splits unquoted words at spaces. Double quotes keep a literal path together while still allowing variable expansion.

A path that starts with `/` begins at the filesystem root. A relative path begins in your current working directory. Run `pwd` before moving files in an unfamiliar project. Use `ls -a` to include names that start with a dot.

`mv SOURCE DESTINATION` changes the original file's location. Check both directories afterward so you can tell a move from a copy.

## Commands and operations

Run `pwd` to locate yourself, then `ls -a "Draft Files"` and `ls -a config` to inspect the two directories. The quotes pass `Draft Files` as one argument. Without them, the shell splits that name into two words before `ls` or `mv` sees it. The `-a` flag includes hidden entries such as `.keep`.

- `mv "Input Files/db settings.txt" settings/db.conf` takes two arguments. The first is the existing file and the second is its new path. Giving a destination directory alone keeps the old filename. Giving the full new path moves and renames the file in one operation.
- `cat settings/db.conf` checks the file contents after the move. Listing the old directory separately tells you whether the original was moved or copied.
- A plain `mv` may replace a destination that already exists. Inspect both paths before running it. For an interactive confirmation before replacement, add `-i`.

A missing-file error usually means the working directory or filename is different than you thought. Run `pwd` and list the quoted directory again. Move the requested file itself so unrelated files, including hidden markers, stay in place.
