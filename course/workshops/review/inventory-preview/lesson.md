Apply the skill from `foundations/core-terms` to a fresh workspace. This exercise has its own files and score, so your original lesson remains available.

## Before you practice

**Outcome:** Choose a row count and extract one field from every inventory record while retaining the original row order.

Read the linked original lesson first if the commands below are unfamiliar. You can return to it from **Do → More practice**. Bash and standard Linux utilities are sufficient.

## Commands and operations

head is the command, -n is an option, 3 is its value, and services.csv is an operand. cut uses -d to choose a delimiter and -f to choose a field. Each input row remains a separate output row.

```bash
head -n 2 tasks.csv
```

This previews two rows. Choose the row count required by your new inventory task.

## Your task

Save the first three rows of services.csv to preview.csv. Save only the second comma-separated field from every row as owners.txt. Preserve services.csv.

The criteria in Do show the points for each result. Work from the root shown as `lab/` in your prompt. You may use any approach that produces the requested state.

## Common mistakes and recovery

Do not sort or deduplicate owners.txt: the task asks for one owner per original row. These records have no quoted commas; a full CSV parser is needed for more complex CSV.

Use `cli-silly inspect` to see the actual file layout. `cli-silly hint` gives a starting point, and later hints show concrete commands. `cli-silly reset` archives this attempt if you need fresh starter files.

## Verify and reflect

Inspect the reports and their input files before running `cli-silly check`. If it passes, explain which state changed and which inspection proves the result. A score records this exercise result; applying the idea to another problem takes additional practice.
