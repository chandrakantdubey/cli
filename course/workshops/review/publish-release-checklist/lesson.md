Apply the skill from `git/local-remote` to a fresh workspace. This exercise has its own files and score, so your original lesson remains available.

## Before you practice

**Outcome:** Publish committed release configuration and notes to a local Git remote, and distinguish working files from published history.

Read the linked original lesson first if the commands below are unfamiliar. You can return to it from **Do → More practice**. Git is required. The remote repository stays local.

## Commands and operations

A release is more useful when its code/configuration and notes are committed together in the published history. git diff --cached shows the next commit. git push moves the remote branch to the commit you send; it does not publish uncommitted files.

```bash
git diff --cached; git push -u origin main
```

Inspect what will be committed before publishing. Push after the relevant commits exist.

## Your task

Set release.txt to release=ready and write release-notes.md containing Verified health and rollback. with a final newline. Commit both files on main and push main to the supplied local origin. The remote is inside this workspace; no hosted account is needed.

The criteria in Do show the points for each result. Work from the root shown as `lab/` in your prompt. You may use any approach that produces the requested state.

## Common mistakes and recovery

If a file is correct in the working tree but the checker cannot see it in HEAD, stage and commit it. If local results pass but the remote criterion fails, inspect git remote -v and push main. One or several focused commits are accepted.

Use `cli-silly inspect` to see the actual file layout. `cli-silly hint` gives a starting point, and later hints show concrete commands. `cli-silly reset` archives this attempt if you need fresh starter files.

## Verify and reflect

Inspect git status and git log --oneline --decorate before running `cli-silly check`. If it passes, explain which state changed and which inspection proves the result. A score records this exercise result; applying the idea to another problem takes additional practice.
