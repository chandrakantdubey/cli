A shell script sees its first argument as `$1`. Double-quote it to preserve spaces in a path. `[ -f "$path" ]` checks whether a regular file is present.

Use stdout for successful output, stderr for error messages, and an exit status for the caller. Shell automation can branch on that status without parsing an error sentence.
