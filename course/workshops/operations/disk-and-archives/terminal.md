## Try it in the terminal

### Inspect without extracting into the workspace

Use `ls logs` and `wc -c logs/incident.txt` before creating an archive. Decide which path should be stored as a member. Inspect a created archive with `tar -tzf`. The member names depend on where you invoked `tar` and what path operands you gave it.

Try `tar -xOzf incident.tar.gz logs/incident.txt` without redirecting at first. Uppercase `-O` sends the member's contents to stdout, so the command does not overwrite the original file. Verify that both the listing and extracted evidence are what the rubric requested. The neighboring `noise.txt` is an input you should leave out of the incident archive.
