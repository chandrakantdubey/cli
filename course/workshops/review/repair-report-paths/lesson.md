Apply the skill from `foundations/filesystem-navigation` to a fresh workspace. This exercise has its own files and score, so your original lesson remains available.

## Before you practice

**Outcome:** Move misplaced reports to the workspace root while preserving their contents and removing the old locations.

Read the linked original lesson first if the commands below are unfamiliar. You can return to it from **Do → More practice**. Bash and standard Linux utilities are sufficient.

## Commands and operations

A move changes a file's location. From the lab root, the full relative source path names the misplaced report. You can alternatively cd into its directory and move it to ../../. Both approaches produce the same required state.

```bash
mv services/web/version.txt version.txt
```

From the root, the source is nested and the destination is beside services/.

## Your task

Two reports were saved in the wrong place. Move services/api/api-owner.txt to api-owner.txt at lab/. Move services/worker/worker-mode.txt to worker-mode.txt at lab/. Preserve their contents and remove the misplaced originals.

The criteria in Do show the points for each result. Work from the root shown as `lab/` in your prompt. You may use any approach that produces the requested state.

## Common mistakes and recovery

cp would leave the misplaced originals behind. Verify the root reports with cat and the source directories with ls. If a path fails, run pwd and inspect the directory before retrying.

Use `cli-silly inspect` to see the actual file layout. `cli-silly hint` gives a starting point, and later hints show concrete commands. `cli-silly reset` archives this attempt if you need fresh starter files.

## Verify and reflect

Inspect the reports and their input files before running `cli-silly check`. If it passes, explain which state changed and which inspection proves the result. A score records this exercise result; applying the idea to another problem takes additional practice.
