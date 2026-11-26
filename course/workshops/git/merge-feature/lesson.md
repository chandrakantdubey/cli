A branch name points at a commit. When your current branch is an ancestor of a feature, Git can fast-forward by moving the branch pointer. With a merge commit, the new commit has both histories as parents.

Use `git log --oneline --graph --all` to inspect the branches. Switch to the branch receiving the change before you merge.

## Commands and operations

Start with `git status --short` and `git log --oneline --graph --all`. The abbreviated commit messages show the history, and the graph shows where the branch pointers sit. The branch you have checked out is the receiving branch.

- `git switch main` selects the receiving branch. It updates the working tree to the files at that branch's tip.
- `git merge feature` incorporates the feature into the current branch. If `main` is an ancestor of `feature`, Git can move `main` to the existing feature commit in a fast-forward. There may be no new merge commit.
- `git branch --show-current` reports the active branch. `git log --oneline --graph --all` lets you compare the branch pointers afterward. Check the resulting file as well as the graph.

Switching to `feature` lets you inspect the file but does not integrate it into `main`. Finish on `main` and verify its committed contents. A standalone edit that looks like the feature is different from retaining the feature's history.
