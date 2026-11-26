Git maintains a working tree, an index, and a committed snapshot. `git diff` shows unstaged edits. `git diff --cached` shows what the next commit will contain. `git commit` records the index.

Inspect the staged diff before committing. If you stage too much, `git restore --staged PATH` removes the path from the index and preserves its working-tree edits.

## Commands and operations

Use `git status --short` to see which files are changed. Its first status column describes the index, and the second describes the working tree. An ` M` entry has an unstaged modification. An `M ` entry already has a staged modification.

- `git diff -- notes.md` previews the working-tree changes to your draft. The `--` separates the path from revisions or options.
- `git add app.cfg` stages the current contents of the fix. `git diff --cached` compares the index with `HEAD` and shows what the next commit will change. Use this check before committing.
- `git commit -m "Fix app configuration"` records the staged snapshot. `git show --stat HEAD` lists the paths in the new commit. Run `git status --short` again to verify the draft is still unfinished in the working tree.

If you staged the draft accidentally, `git restore --staged notes.md` puts its index entry back to the last commit and keeps the edits on disk. This exercise grades both the contents of the new commit and where the draft remains.
