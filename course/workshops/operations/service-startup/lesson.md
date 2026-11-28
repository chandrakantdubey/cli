A systemd service starts a process with the user, environment, and working directory from its unit. Use a journal to distinguish failures before the process starts from failures reported by the application itself.

This exercise uses a journal and unit fixture. On a live system, `systemctl cat UNIT` shows the unit, and `journalctl -u UNIT -n 30` shows recent messages. Keep your edits inside the lab.

## Commands and operations

Read the provided unit and journal together. The unit specifies what should happen, while the journal records what actually happened. A failure such as `CHDIR` can occur before the application starts, so application logs alone may contain nothing.

- Use `grep -n 'CHDIR' service.journal` to find the matching journal lines with their line numbers. Read nearby entries to recover the directory and error details.
- Use `grep '^WorkingDirectory=' orders.service` to find the active directive. The leading anchor avoids matching text in a comment or another value.
- Compare that path with the lab's deployment record. Save the correct path and the full corrected directive to their separate output files. A report with just the directory is not a complete unit directive.

On a running system, `systemctl cat UNIT` displays a unit and its drop-ins, and `journalctl -u UNIT -n 30` shows recent service logs. In this exercise, the named files provide the same diagnostic context and all writes stay in the workspace. Check the reports for a stray prefix or trailing space before grading.

## Before you practice

**Outcome:** Connect captured startup logs with a service working-directory configuration.

Know absolute paths and basic log filtering. The exercise uses local captures and needs no systemd privileges.

## Common mistakes and recovery

Do not change your host services. Edit the provided unit fixture and use the journal evidence to justify the path.

Inspect the relevant input and output before repeating a command. Use `cli-silly hint` for a starting point and `cli-silly check` for criterion-by-criterion feedback. If the starter is damaged, `cli-silly reset` archives this attempt before recreating it.

## Use it in development

Investigate a service that works manually but fails under its service manager.

After passing, explain which files or state changed and how you verified them. Open related practice from **Do** to apply the idea to different data.
