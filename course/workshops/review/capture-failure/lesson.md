Apply the skill from `foundations/mental-model` to a fresh workspace. This exercise has its own files and score, so your original lesson remains available.

## Before you practice

**Outcome:** Keep stdout, stderr, and the original failure status as three distinct diagnostic reports.

Read the linked original lesson first if the commands below are unfamiliar. You can return to it from **Do → More practice**. Bash and standard Linux utilities are sufficient.

## Commands and operations

A process has two output streams and a separate exit status. The shell variable $? holds the last command status. In status=0; command || status=$?, the assignment on the right runs only when command fails. It captures that failure before another command replaces it.

```bash
status=0; bash verify.sh > result.txt 2> error.txt || status=$?
```

Save the status immediately. `>` selects stdout and `2>` selects stderr. Next run `printf '%s\n' "$status" > exit.txt`: `%s` prints the saved value and `\n` adds a newline. Save it before inspecting the reports so later commands cannot replace the original result.

## Your task

Run bash bin/preflight.sh. Capture stdout in preflight.out, stderr in preflight.err, and its numeric exit status followed by a newline in preflight.status. Preserve the starter script.

The criteria in Do show the points for each result. Work from the root shown as `lab/` in your prompt. You may use any approach that produces the requested state.

## Common mistakes and recovery

Writing echo $? after cat records cat's status instead. A nonempty stdout report also does not prove success. Read all three pieces of evidence.

Use `cli-silly inspect` to see the actual file layout. `cli-silly hint` gives a starting point, and later hints show concrete commands. `cli-silly reset` archives this attempt if you need fresh starter files.

## Verify and reflect

Inspect the reports and their input files before running `cli-silly check`. If it passes, explain which state changed and which inspection proves the result. A score records this exercise result; applying the idea to another problem takes additional practice.
