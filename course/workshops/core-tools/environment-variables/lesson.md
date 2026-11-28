Run the same tool twice with different settings. This lab shows which variables a child inherits and how to set one only for a single command.

## Commands and operations

`APP_MODE=preview bash bin/show-mode.sh` supplies `APP_MODE` to this invocation of Bash. When the command returns, the setting used by your current shell stays as it was. Use one-command assignments for temporary experiments.

`export APP_MODE=release` sets the name in this shell and marks it for inheritance. Now child processes including `bash` can read it. Use `unset APP_MODE` to remove the setting when you are done. These changes last only for this practice shell.

The fixture uses `${APP_MODE:-local}`. Bash expands that to `local` when the name is unset or empty. Inspect a child with `printenv APP_MODE` after export. Avoid dumping the whole environment when it may include secrets.

## Inspect

```bash
printf '%s\n' "$HOME"
printf '%s\n' "$PATH"
env | sort
```

Shell variables and exported environment variables are related but not identical.

```bash
name="Ada"
export APP_ENV=dev
```

`name` is a shell variable. `APP_ENV` is exported to child processes.

## PATH

`PATH` is a colon-separated list of directories used to locate executable commands.

```bash
printf '%s\n' "$PATH" | tr ':' '\n'
command -v node
```

Avoid placing the current directory at the front of `PATH` unless you understand the security implications.

## Temporary environment

Set a variable for one command:

```bash
APP_ENV=production ./server
```

This does not permanently modify your shell.

## Persistent configuration

Interactive shell configuration is commonly stored in files such as:

```text
~/.bashrc
~/.bash_profile
~/.profile
~/.zshrc
```

Which files are read depends on shell and login/interactive mode.

## Working directory

A process has a current working directory:

```bash
pwd
```

Scripts should not assume they start in the project directory. Establish paths deliberately.

## Locale

Text-processing behavior can depend on locale:

```bash
locale
LC_ALL=C sort file.txt
```

For machine-oriented scripts, explicitly controlling locale can prevent surprising ordering or character-class behavior.

## Secrets

Environment variables are not a secure secret store. They can be exposed through process inspection, debugging, crash reports, or application behavior.

Do not put long-lived production credentials into examples.

## Before you practice

**Outcome:** Control a child process environment without unintentionally changing later commands.

Read the mental model and understand that Bash starts child processes.

## Common mistakes and recovery

A shell variable is not automatically exported. A prefix assignment affects one command; export affects future children.

Inspect the relevant input and output before repeating a command. Use `cli-silly hint` for a starting point and `cli-silly check` for criterion-by-criterion feedback. If the starter is damaged, `cli-silly reset` archives this attempt before recreating it.

## Use it in development

Run the same application in preview and release modes from one shell.

After passing, explain which files or state changed and how you verified them. Open related practice from **Do** to apply the idea to different data.
