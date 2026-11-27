## Try it in the terminal

### Test the command as another program sees it

Read the starter and the quoted sample filename. After adding the success path, run `node line-count.mjs "sample input.txt" > scratch.stdout 2> scratch.stderr`. Save the status before using `cat`. The stdout file should contain one count and its newline, and the stderr file should stay empty.

Repeat with no argument, an unknown flag, and a missing path. Invalid invocation needs status 2, while an unreadable input needs status 1. Try a scratch input without a final newline using `printf '%s' worker > scratch-input.txt`. It still has one non-empty line. That distinction is one reason to test behavior with more than the provided sample.
