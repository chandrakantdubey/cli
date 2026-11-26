Git maintains a working tree, an index, and a committed snapshot. `git diff` shows unstaged edits. `git diff --cached` shows what the next commit will contain. `git commit` records the index.

Inspect the staged diff before committing. If you stage too much, `git restore --staged PATH` removes the path from the index and preserves its working-tree edits.
