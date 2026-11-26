The shell splits unquoted words at spaces. Double quotes keep a literal path together while still allowing variable expansion.

A path that starts with `/` begins at the filesystem root. A relative path begins in your current working directory. Run `pwd` before moving files in an unfamiliar project. Use `ls -a` to include names that start with a dot.

`mv SOURCE DESTINATION` changes the original file's location. Check both directories afterward so you can tell a move from a copy.
