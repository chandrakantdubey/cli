Use `du -h` to inspect the size of directories and `df -h` for free space on the containing filesystem. Before archiving, inspect the files you mean to include. `tar -czf archive.tar.gz paths` writes a compressed archive and `tar -tzf` lists its entries.

## Commands and operations

Use `du -h logs` to inspect the directory's disk usage. `df -h .` reports how much space is available on the filesystem containing the current directory. One describes the selected files, and the other describes the storage volume.

- `tar -czf handoff.tar.gz logs/api.txt` creates an archive, compresses it with gzip, and writes it to the filename after `-f`. Supplying a single source file keeps other logs out.
- `tar -tzf handoff.tar.gz` lists member names without extracting files. Save the list to the requested report. Use the member name exactly as stored, including any directory prefix.
- `tar -xOzf handoff.tar.gz logs/api.txt` extracts the named member to standard output. The capital `O` lets you inspect or redirect its content without unpacking it over your original files.

If a member cannot be found, list the archive and compare the full stored path. Check both the archive's contents and the extracted evidence before grading. An archive of the whole logs directory would include material the task did not ask for.
