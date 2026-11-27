A team keeps its published notes in a separate folder. Move exactly one document and record where a teammate can find it. The starter has a similar filename and a hidden placeholder that should stay in the draft folder.

## Commands and operations

In `ls -la "Draft Notes"`, `-a` includes names beginning with a dot and `-l` shows details. The double quotes keep the space in the directory name inside one shell argument.

`mv -- "Draft Notes/api plan.md" "Ready Notes/api.md"` gives `mv` two paths. The `--` ends option parsing. Because the second operand is a filename inside an existing directory, the command moves and renames the one file.

Inspect both directories after moving. Use `printf '%s\n' 'Ready Notes/api.md'` to print an example relative path with a final newline. A report filename written after `>` is resolved from the directory where you run the command.
