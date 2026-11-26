In Git, a remote is a named repository address. A bare repository stores objects and refs without a working directory. This lab uses one in the workspace, so the push needs no account or network.

`git remote -v` shows where pushes go. `git push -u origin main` updates the branch on the remote and sets its upstream.
