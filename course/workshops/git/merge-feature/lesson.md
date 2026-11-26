A branch name points at a commit. When your current branch is an ancestor of a feature, Git can fast-forward by moving the branch pointer. With a merge commit, the new commit has both histories as parents.

Use `git log --oneline --graph --all` to inspect the branches. Switch to the branch receiving the change before you merge.
