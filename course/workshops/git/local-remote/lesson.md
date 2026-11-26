In Git, a remote is a named repository address. A bare repository stores objects and refs without a working directory. This lab uses one in the workspace, so the push needs no account or network.

`git remote -v` shows where pushes go. `git push -u origin main` updates the branch on the remote and sets its upstream.

## Commands and operations

Before pushing, run `git remote -v`. Each entry lists a remote name and its fetch or push destination. In this lab, the target is a bare repository on disk. Use the supplied `origin` so the checker can read the expected ref.

- Edit `release.txt`, then use `git diff -- release.txt` to see the change. Stage the file with `git add release.txt` and commit the staged change. Editing alone leaves the change outside the commit history.
- `git push -u origin main` uploads commits reachable from `main` and updates the remote's branch. The `-u` registers an upstream for later pushes and pulls. It does not include unstaged edits.
- `git log -1 --oneline` shows your latest local commit. `git ls-remote origin refs/heads/main` displays the object ID stored for the branch in the remote repository. Compare that with `git rev-parse HEAD`.

If the remote state is behind, first check whether there is a new local commit to push. This lab uses a local repository, so GitHub authentication and network access are unnecessary.
