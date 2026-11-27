Git maintains a working tree, an index, and a committed snapshot. `git diff` shows unstaged edits. `git diff --cached` shows what the next commit will contain. `git commit` records the index.

Inspect the staged diff before committing. If you stage too much, `git restore --staged PATH` removes the path from the index and preserves its working-tree edits.

## Commands and operations

Use `git status --short` to see which files are changed. Its first status column describes the index, and the second describes the working tree. An ` M` entry has an unstaged modification. An `M ` entry already has a staged modification.

- `git diff -- notes.md` previews the working-tree changes to your draft. The `--` separates the path from revisions or options.
- `git add app.cfg` stages the current contents of the fix. `git diff --cached` compares the index with `HEAD` and shows what the next commit will change. Use this check before committing.
- `git commit -m "Fix app configuration"` records the staged snapshot. `git show --stat HEAD` lists the paths in the new commit. Run `git status --short` again to verify the draft is still unfinished in the working tree.

If you staged the draft accidentally, `git restore --staged notes.md` puts its index entry back to the last commit and keeps the edits on disk. This exercise grades both the contents of the new commit and where the draft remains.

## Repository state

```bash
git status
git branch --show-current
git log --oneline --decorate -10
```

Think of Git as three major local areas:

```text
working tree -> staging area -> commit
```

## Basic workflow

```bash
git add README.md
git diff --cached
git commit -m "docs: improve README"
```

Use `git diff` for unstaged changes and `git diff --cached` for staged changes.

## Branches

```bash
git switch -c feature/health-check
git switch main
git branch -vv
```

Branches are lightweight references to commits.

## Remote work

```bash
git remote -v
git fetch origin
git switch -c feature-x origin/main
git push -u origin feature-x
```

`fetch` downloads remote state without merging it into your current branch.

## Merge and rebase

Merge preserves branch topology:

```bash
git merge feature-x
```

Rebase replays commits on a new base:

```bash
git rebase main
```

Rebase rewrites commit IDs, so avoid rebasing shared history unless your team explicitly uses that workflow.

## Undoing changes

Discard unstaged changes carefully:

```bash
git restore path/to/file
```

Unstage:

```bash
git restore --staged path/to/file
```

Create a new commit that reverses an earlier commit:

```bash
git revert <commit>
```

Avoid `git reset --hard` until you understand exactly which uncommitted work it will discard.

## Inspect history

```bash
git show <commit>
git log --graph --oneline --all
git blame path/to/file
```

## .gitignore

Keep generated files, local environments, secrets, and other non-source artifacts out of the repository where appropriate.

Never commit credentials just because `.gitignore` was missing.
