Apply the skill from `foundations/cli-intro` to a fresh workspace. This exercise has its own files and score, so your original lesson remains available.

## Before you practice

**Outcome:** Produce a handoff note and a repeatable directory inventory without changing either source note.

Read the linked original lesson first if the commands below are unfamiliar. You can return to it from **Do → More practice**. Bash and standard Linux utilities are sufficient.

## Commands and operations

Reading a path does not change the working directory. cat takes the input path as an argument. Bash opens the output path relative to its current directory before starting cat. Use pwd before creating your reports.

```bash
cat updates/today.txt > summary.txt
```

The example copies a note into the current directory. In this task use the provided notes and the requested report name.

## Your task

Read notes/next.txt and save it as handoff.txt at lab/. Save a sorted, one-name-per-line listing of notes/ as notes-list.txt. Keep both starter notes unchanged.

The criteria in Do show the points for each result. Work from the root shown as `lab/` in your prompt. You may use any approach that produces the requested state.

## Common mistakes and recovery

If the output lands inside notes/, you changed directory before redirecting. Return with cd "$CLI_SILLY_LAB_ROOT" and move the report to the root. Keep reports outside the input directory so the listing contains just the starter notes.

Use `cli-silly inspect` to see the actual file layout. `cli-silly hint` gives a starting point, and later hints show concrete commands. `cli-silly reset` archives this attempt if you need fresh starter files.

## Verify and reflect

Inspect the reports and their input files before running `cli-silly check`. If it passes, explain which state changed and which inspection proves the result. A score records this exercise result; applying the idea to another problem takes additional practice.
