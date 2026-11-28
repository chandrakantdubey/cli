# Recover local work

Practice runs with your normal Linux user permissions. A workspace is an organized exercise directory. It is not a sandbox that prevents commands from accessing other files.

## Inspect before changing

Run `cli-silly inspect` inside a lab, including from a subdirectory, or pass its lesson ID from elsewhere. It shows up to 100 entries and four directory levels. Git labs also show status and recent history. Control characters in displayed filenames are escaped. Inspection neither executes learner scripts nor changes grades.

## Reset and restore an attempt

`cli-silly reset` asks before archiving the current directory and recreating the starter. Stop any processes or Compose project from that lab first. Earned grades remain available.

Run `cli-silly archives LESSON-ID` to find saved attempts. Restore one from your normal shell:

```text
cli-silly restore LESSON-ID .saved-NAME
```

Use the actual ID and `.saved-...` directory name printed by the listing. After confirmation, the current attempt is also archived. Restore does not change progress. An old open shell still belongs to its original attempt; reopen Bash from the lesson before checking work.

## Export and import a backup

Close other cli-silly sessions and stop lab services first. The archive includes workspaces, archived attempts, installed courses, progress, and reading state. It may contain whatever you wrote in your labs. Keep it private unless you have reviewed its contents.

```bash
cli-silly backup export "$HOME/cli-silly-backup.tar.gz"
```

Export refuses an existing destination, links, special files, more than 20,000 entries, or more than 256 MiB. Move large generated artifacts aside or use your usual backup tool if necessary. Save the archive outside cli-silly data/state directories.

Import requires empty application data and state directories. To review a backup without disturbing your current installation:

```bash
XDG_DATA_HOME="$HOME/cli-silly-recovery/data" \
XDG_STATE_HOME="$HOME/cli-silly-recovery/state" \
cli-silly backup import "$HOME/cli-silly-backup.tar.gz"

XDG_DATA_HOME="$HOME/cli-silly-recovery/data" \
XDG_STATE_HOME="$HOME/cli-silly-recovery/state" \
cli-silly
```

Keep using those same XDG locations to use the restored copy. Import validates and stages the archive before publishing it, rejects escaping paths and links, and refuses collisions. It does not merge records or overwrite the original directories.

## Clean up runtime resources

Use `jobs -l` inside the shell that started an API and terminate the specific job when finished. In a Compose lab, `docker compose ps` shows the project and `docker compose down` stops it. Only use `docker compose down -v` when you intend to discard that lab's volume data. Backup does not capture running processes or Docker volumes.
