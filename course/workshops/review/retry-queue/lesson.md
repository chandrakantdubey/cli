A dispatcher printed usable results before failing. Recover a clean retry queue and preserve the status for the next operator. Three channels carry different information. stdout contains successful dispatches, stderr carries warnings and retry messages, and the exit code indicates the overall run failed.

## Commands and operations

`> sent.out` connects file descriptor 1 to a file. `2> messages.log` does the same for file descriptor 2. Inspect each file separately to see why the terminal originally mixed the lines.

Store `$?` before inspecting either file. If the dispatch returns nonzero, the right side of `|| status=$?` runs immediately and preserves that code even when using `set -e`.

`grep '^retry ' messages.log` copies only diagnostic lines starting with the word retry and a space. It keeps their order. Write the status with `printf '%s\n' "$status"` rather than redirecting a message that contains the number.
