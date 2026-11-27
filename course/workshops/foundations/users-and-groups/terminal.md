## Try it in the terminal

### Read an account record like a table

Inspect `captures/passwd.txt` with `head -n 2`. Split fields with `awk -F:` and print the account name, UID, and primary GID together. Compare the similarly named `runner-old` and `runner` records. A substring search would select both, so match the first field exactly in your UID report.

Inspect the groups capture separately. The primary GID from the account record maps to field 3 of the group record. A name listed in field 4 can indicate supplementary membership rather than the primary group. Try `id -u` in the lab to see what your real shell reports, then use the supplied captures for the scored task.
