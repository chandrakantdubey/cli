Learn how to read a command line and use its help to choose an option. A deployment inventory gives you two independent ways to practice. First choose the lines to read, then choose a field to print.

## Commands and operations

In `head -n 2 inventory.csv`, `head` is the command, `-n` is an option that takes `2` as its value, and `inventory.csv` is the file operand. `head --help` shows options on GNU Linux. Bash's `type head` shows how it resolves the command and `help cd` documents the Bash built-in `cd`.

`cut -d ',' -f 2 inventory.csv` splits each line on commas and prints the second field, counting from one. This works here because the fixture has no quoted commas. Use a CSV parser for general CSV. Read an option's documentation before combining it with another option.

An ordinary command can accept filenames that begin with a dash. With utilities that support it, put `--` before the filename to end option parsing. For example `cat -- -draft.txt` treats `-draft.txt` as an operand.

## Anatomy of a command

```bash
grep -n "ERROR" app.log
```

- command: `grep`
- option: `-n`
- pattern argument: `ERROR`
- input file argument: `app.log`

Not every option is a simple boolean. Some options consume values:

```bash
head -n 20 app.log
```

## Discover before guessing

```bash
type cd
command -v grep
grep --help
man grep
```

For shell built-ins, use:

```bash
help cd
```

## Exit status

Do not assume every non-zero status means the same thing. Programs define their own status meanings.

For example, search tools may use one non-zero status to mean "no match" rather than "fatal error."

## Root and sudo

`root` is UID 0. `sudo` allows an authorized user to execute a command with elevated privileges.

Prefer:

```bash
sudo systemctl restart nginx
```

over starting a permanent root shell when you only need one privileged action.

## Before you practice

**Outcome:** Identify a command, its options, and its operands in a data-inspection task.

Read CLI & Linux Intro and practice reading an ordinary file.

## Common mistakes and recovery

The -n option consumes a number. The filename is a separate operand. CSV here is deliberately simple; cut does not parse quoted CSV fields.

Inspect the relevant input and output before repeating a command. Use `cli-silly hint` for a starting point and `cli-silly check` for criterion-by-criterion feedback. If the starter is damaged, `cli-silly reset` archives this attempt before recreating it.

## Use it in development

Preview generated inventories and select a field before automating a larger report.

After passing, explain which files or state changed and how you verified them. Open related practice from **Do** to apply the idea to different data.
