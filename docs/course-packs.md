# Offline course packs

The built-in developer workshop is available as soon as cli-silly is installed. A course pack adds another learning path while keeping its workspaces and grades under its own lesson IDs.

## Install and use a pack

From this checkout:

```bash
cli-silly course install examples/markdown-basics
cli-silly course list
cli-silly course use markdown-basics
cli-silly
```

The same installer accepts a local tar.gz archive. Download `markdown-basics-v1.tar.gz` from the release assets, then run `cli-silly course install markdown-basics-v1.tar.gz`. Installation works offline once you have the file.

The course map shows installed packs beside the built-in workshops. Use `cli-silly learn markdown-basics/project-heading` to open the example directly. `cli-silly continue` prioritizes your selected course. Use `cli-silly course use builtin` to switch back. `cli-silly resume` opens the last-read section across courses.

Install a higher version of the same pack to upgrade. A stable lesson ID and grading revision retain the existing attempt and score. A failed upgrade leaves the previous version available.

## Author a declarative pack

Use the layout in `examples/markdown-basics/`. Put `pack.json` at the root with `schemaVersion: 1`, `engineSchema: 1`, a stable lowercase ID, a title, a positive integer version, and an ordered list of lesson slugs. Store each lesson under `lessons/<slug>/` with `lesson.json`, `lesson.md`, and `starter/`. A lesson ID is `<pack-id>/<lesson-slug>`.

Keep file checks declarative. Installed packs support text comparisons, absent paths, and owner-executable mode. The CLI stages and validates the pack before activation. Archives can contain regular files and directories up to the installed limits of 1 MiB per file, 8 MiB of file contents, and 1,024 entries. Symlinks, escaping paths, filenames with control characters, incompatible versions, custom setup, and executable checker types are rejected.

Use Bash for hands-on practice only after reviewing commands and files. Installation reads data and copies assets. To build an archive, put the manifest at its root:

```bash
tar -czf markdown-basics-v1.tar.gz -C examples/markdown-basics .
```
