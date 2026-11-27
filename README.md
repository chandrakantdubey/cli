# cli-silly

cli-silly is a terminal-first learning platform for developers. Its built-in course teaches Linux, Git, Bash, and Docker by stepping through terminal demos, working in a real shell, and getting a score for what you accomplish. The full course is bundled with the CLI. The companion site uses the same lessons and rubrics.

## Install and start on Linux

Use the [installation guide](docs/install.md) to install from a clone or source ZIP, or from a [prebuilt GitHub release](https://github.com/chandrakantdubey/cli/releases/latest). The CLI bundles the course for offline reading and practice.

To install the latest Linux release directly from GitHub:

```bash
curl -fsSL https://raw.githubusercontent.com/chandrakantdubey/cli/main/install.sh | bash
```

After installing, run:

```bash
"$HOME/.local/bin/cli-silly" doctor
"$HOME/.local/bin/cli-silly"
```

Use `cli-silly` directly when its install directory is on your PATH.

Start with **Foundations** and then **Core tools**. Use ↑ and ↓ to select a lesson or **/** to search by topic, tool, or command. Press **Enter** to **Read** or **d** to **Do**. Press **m** to return to the lessons and **?** for the selected lesson's commands.

In **Read**, scroll through the full chapter with **↑ / ↓**. Press **i** to jump through the contents and **e** for the worked example. Foundations and Core tools include a **Try it in the terminal** section with extra drills. Press **b** to bookmark a section or run `cli-silly resume` to return to your last section.

In **Do**, read the task and scoring criteria, then press **Enter** to open Bash. The prompt shows `lab/` at the workspace root and adds the current subdirectory as you navigate. Run `cli-silly read` for the current lesson, `cli-silly task` for the goal, and `cli-silly check` to see your score. Type `exit` to return. File-only exercises refresh their score on return. Press **]** after passing to read the next lesson.

Press **z** to take a quiz when the selected lesson has one. Press **V** to see scheduled review challenges. Review scores and quiz answers are stored separately from the original lab grade.

### Tools by workshop

| Workshop | Tools |
| --- | --- |
| Shell, streams, Bash | Standard Linux utilities and Bash. The JSON lab also needs `jq`. |
| Core tools | `ps` for process inspection, `jq` for advanced text processing, and Node.js for the Node lab. The identity and package labs use local captures. |
| Git | Git with `git switch` support. Exercises use local repositories and need no account. |
| Linux app debugging and capstone | Python 3 and `curl`. The exercise starts a service on a free loopback port. |
| Docker | Docker Engine, the Compose plugin, and `curl`. Confirm access to the daemon with `cli-silly doctor`. |
| Service operations | `tar` and OpenSSH's `ssh` client for their respective lessons. The systemd exercise uses a local journal and unit file. |

Docker labs use `python:3.12-alpine`. Pull it during setup with `docker pull python:3.12-alpine` to prepare for offline practice. The labs allocate separate Compose project names and bind published ports to loopback. Run `docker compose down` in the lab when finished. For the volume exercise, `docker compose down -v` also removes the exercise's saved data. Stop background API processes from their practice shell when you finish those labs.

## Use it from your own shell

```bash
cli-silly list
cli-silly read streams/log-report
cli-silly do streams/log-report
cd "$(cli-silly path streams/log-report)"
cli-silly task
cli-silly hint 2
cli-silly check
cli-silly progress
```

Inside a lab, `task`, `hint`, `check`, `progress`, and `reset` identify the lesson from the workspace, including from a subdirectory. Outside the lab, use a lesson ID to target it. Run `cli-silly progress` outside a lab for all scores or `cli-silly continue` for the first unfinished lesson.

`check` returns 0 for completion, 1 for work remaining, and 2 for an environment or checker error. Add `--json` for structured results. Reopening a lab preserves your files. `reset` asks for confirmation, archives the previous attempt, and retains your records.

Workspaces are stored under `${XDG_DATA_HOME:-~/.local/share}/cli-silly/workspaces`. Grades are tracked under `${XDG_STATE_HOME:-~/.local/state}/cli-silly/progress`. Redirected output and `TERM=dumb` use plain output. Set `NO_COLOR=1` for a color-free full-screen UI.

## Commands and starter templates

Press **n** in the terminal UI to browse the reference library. Use **/** to search and Enter to preview an example. Press **p** to print the selected example back to your shell.

```bash
cli-silly snippets awk
cli-silly snippet find-server-errors 'log=logs/worker API.log'
```

Run `cli-silly template config-check "my checker.sh"` from a lab workspace to generate a new file. If that filename is taken, choose another path.

## Add an offline course

Install the optional Markdown Basics pack from the checkout or use its release archive:

```bash
cli-silly course install examples/markdown-basics
cli-silly course use markdown-basics
cli-silly learn markdown-basics/project-heading
```

Installed packs have their own workspaces and scores. [Course pack guide](docs/course-packs.md)

## Run the companion site

Use Node.js 22.12 or newer on the Node 22 line and npm.

```bash
npm ci
npm run dev
```

The sidebar lists every core workshop and keeps existing reference pages accessible.

## Develop and verify

```bash
make setup
make run-dev
make check
```

`make run-dev` uses checkout-local data under `.local/`. Run `make docker-check` separately with a working Docker daemon. Use `make clean` to clear generated binaries, website builds, and release archives.

See [local development](docs/development.md) for prerequisites and troubleshooting, [AGENTS.md](AGENTS.md) for contributor instructions, and the [roadmap](plan.md) for the next platform phases.
