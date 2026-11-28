Use `du -h` to inspect the size of directories and `df -h` for free space on the containing filesystem. Before archiving, inspect the files you mean to include. `tar -czf archive.tar.gz paths` writes a compressed archive and `tar -tzf` lists its entries.

## Commands and operations

Use `du -h logs` to inspect the directory's disk usage. `df -h .` reports how much space is available on the filesystem containing the current directory. One describes the selected files, and the other describes the storage volume.

- `tar -czf handoff.tar.gz logs/api.txt` creates an archive, compresses it with gzip, and writes it to the filename after `-f`. Supplying a single source file keeps other logs out.
- `tar -tzf handoff.tar.gz` lists member names without extracting files. Save the list to the requested report. Use the member name exactly as stored, including any directory prefix.
- `tar -xOzf handoff.tar.gz logs/api.txt` extracts the named member to standard output. The capital `O` lets you inspect or redirect its content without unpacking it over your original files.

If a member cannot be found, list the archive and compare the full stored path. Check both the archive's contents and the extracted evidence before grading. An archive of the whole logs directory would include material the task did not ask for.

## tar

Create an archive:

```bash
tar -cf project.tar project/
```

Extract:

```bash
tar -xf project.tar
```

List without extracting:

```bash
tar -tf project.tar
```

Create gzip-compressed archive:

```bash
tar -czf project.tar.gz project/
```

Modern `tar` can infer compression from common extensions in some workflows, but explicit flags are useful for teaching.

## bzip2 and xz

```bash
tar -cjf project.tar.bz2 project/
tar -cJf project.tar.xz project/
```

Compression trades CPU time for size; the best choice depends on workload.

## Zip

```bash
zip -r project.zip project/
unzip project.zip
unzip -l project.zip
```

ZIP is common for cross-platform interchange. `tar` plus compression is common in Unix/Linux workflows.

## Verify before extraction

Treat archives as untrusted input. Inspect their contents before extracting:

```bash
tar -tf archive.tar
unzip -l archive.zip
```

Be cautious of path traversal or unexpected absolute paths in archives received from untrusted sources.

## Backups

An archive is not automatically a backup. A backup strategy also needs:

- retention;
- integrity checks;
- off-host copies;
- access controls;
- restoration testing.

## Reproducibility

For automated packaging, be aware that timestamps, ownership, ordering, and metadata can make archives differ even when source content is the same.

## Before you practice

**Outcome:** Create an archive and verify both its member names and recovered contents.

Know relative paths and redirection. tar is required.

## Common mistakes and recovery

The path stored inside an archive depends on where tar starts. List members before extracting and keep extraction inside the lab.

Inspect the relevant input and output before repeating a command. Use `cli-silly hint` for a starting point and `cli-silly check` for criterion-by-criterion feedback. If the starter is damaged, `cli-silly reset` archives this attempt before recreating it.

## Use it in development

Package support evidence so another developer can inspect it reliably.

After passing, explain which files or state changed and how you verified them. Open related practice from **Do** to apply the idea to different data.
