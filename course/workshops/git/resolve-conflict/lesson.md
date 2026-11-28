A conflict means Git needs help choosing the final content. `git status` lists unresolved paths. Conflict markers show the current branch and the incoming branch. Edit the file to the intended combined result, stage it, and complete the merge.

Look at the graph afterward. Both branch histories should be reachable from the result.

## Commands and operations

Run `git merge feature` from `main`. For this fixture, Git reports a conflict because both branches edited the same line differently. `git status --short` can show `UU`, which means both sides modified the file and the index still contains an unresolved path.

- Open `settings.txt` and inspect the sections between `<<<<<<<`, `=======`, and `>>>>>>>`. The upper and lower variants show the two versions Git could not combine automatically.
- Replace the entire marker block with the intended result. Saving the file changes the working tree. Run `git add settings.txt` to tell Git that the index now contains your chosen resolution.
- `git commit -m "Resolve theme conflict"` finishes the merge once all conflicts are staged. Check the resulting content with `git show HEAD:settings.txt` and inspect the graph with `git log --oneline --graph --all`.

Check `git status` if Git refuses the commit. A remaining unmerged path needs a resolution. Leaving conflict markers in a staged file can produce a commit while still giving the application invalid contents.

## Before you practice

**Outcome:** Resolve conflicting edits and complete a merge deliberately.

Practice branches, merging, staging, and commits.

## Common mistakes and recovery

Deleting conflict markers is only part of the task. The resolved content must satisfy the requested behavior, then be staged and committed.

Inspect the relevant input and output before repeating a command. Use `cli-silly hint` for a starting point and `cli-silly check` for criterion-by-criterion feedback. If the starter is damaged, `cli-silly reset` archives this attempt before recreating it.

## Use it in development

Reconcile two developers changing the same configuration.

After passing, explain which files or state changed and how you verified them. Open related practice from **Do** to apply the idea to different data.
