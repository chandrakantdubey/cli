## Try it in the terminal

### Compare the three places Git can see a change

Run `git status --short` and inspect its two columns. Use `git diff -- app.cfg notes.md` before adding anything. After adding only `app.cfg`, run `git diff --cached` and `git diff`. The staged diff should include the configuration fix, and the working-tree diff should still show the notes.

If you stage the notes by accident, use `git restore --staged -- notes.md` to put the index back without overwriting your notes. Inspect the staged changes before committing. After the commit, check both `git show --stat --oneline HEAD` and `git status --short`. Both the commit and the remaining draft matter for the score.
