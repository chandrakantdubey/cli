Apply the skill from `shell/make-it-runnable` to a fresh workspace. This exercise has its own files and score, so your original lesson remains available.

## Before you practice

**Outcome:** Repair a script owner’s execute permission and capture the output of a direct run without editing the script.

Read the linked original lesson first if the commands below are unfamiliar. You can return to it from **Do → More practice**. Bash and standard Linux utilities are sufficient.

## Commands and operations

Linux checks permission bits when executing a file by path. chmod u+x adds execute permission for the owner without changing unrelated permission bits. The first line chooses Bash through env. The checker measures permissions and output; it does not record your shell history.

```bash
chmod u+x bin/show-build
```

Add just the owner execute bit, then run the script by its path.

## Your task

Give the owner execute permission on bin/release-version, run it directly as ./bin/release-version, and save stdout as version.txt. Keep the script contents unchanged.

The criteria in Do show the points for each result. Work from the root shown as `lab/` in your prompt. You may use any approach that produces the requested state.

## Common mistakes and recovery

Running bash bin/release-version bypasses the need for the owner execute bit. A successful direct run is the useful verification. Broad permissions such as 777 are unnecessary.

Use `cli-silly inspect` to see the actual file layout. `cli-silly hint` gives a starting point, and later hints show concrete commands. `cli-silly reset` archives this attempt if you need fresh starter files.

## Verify and reflect

Inspect the reports and their input files before running `cli-silly check`. If it passes, explain which state changed and which inspection proves the result. A score records this exercise result; applying the idea to another problem takes additional practice.
