Apply the skill from `foundations/users-and-groups` to a fresh workspace. This exercise has its own files and score, so your original lesson remains available.

## Before you practice

**Outcome:** Resolve an account UID and primary group by joining captured passwd and group records on the numeric GID.

Read the linked original lesson first if the commands below are unfamiliar. You can return to it from **Do → More practice**. Bash and standard Linux utilities are sufficient.

## Commands and operations

In a passwd record, the third colon-delimited field is the UID and the fourth is the primary GID. Group records put the group name first and the numeric GID third. Use that shared number to relate the two files. awk -v passes a shell value to an awk variable.

```bash
awk -F: '$1 == "builder" {print $4}' captures/passwd.txt
```

Select one account and read its primary GID before resolving the group name.

## Your task

Use captures/passwd.txt to find the UID and primary GID of appsvc. Save the UID in uid.txt and its primary group name from captures/group.txt in group.txt. One value and newline per report. Preserve both captures.

The criteria in Do show the points for each result. Work from the root shown as `lab/` in your prompt. You may use any approach that produces the requested state.

## Common mistakes and recovery

appsvc appears in the operators member list, but that is supplemental membership. The primary group comes from its passwd GID. You do not need sudo or changes to your computer's users.

Use `cli-silly inspect` to see the actual file layout. `cli-silly hint` gives a starting point, and later hints show concrete commands. `cli-silly reset` archives this attempt if you need fresh starter files.

## Verify and reflect

Inspect the reports and their input files before running `cli-silly check`. If it passes, explain which state changed and which inspection proves the result. A score records this exercise result; applying the idea to another problem takes additional practice.
