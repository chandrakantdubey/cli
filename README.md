# cli-silly

cli-silly is a terminal-first learning platform for developers. Its built-in course teaches Linux, Git, Bash, and Docker by stepping through terminal demos, working in a real shell, and getting a score for what you accomplish. The full course is bundled with the CLI. The companion site uses the same lessons and rubrics.

## Install and start on Linux

Use the [installation guide](docs/install.md) to install from a clone or source ZIP, or from a prebuilt GitHub release when one is available. The CLI bundles the course for offline reading and practice.

After installing, run:

```bash
"$HOME/.local/bin/cli-silly" doctor
"$HOME/.local/bin/cli-silly"
```

Use `cli-silly` directly when its install directory is on your PATH.

Use ↑ and ↓ to select a lesson or **/** to search by lesson, workshop, tool, or command. On a wide terminal the map shows a lesson preview beside the learning paths. Enter opens its visual demo. Use ← and → to move through the examples, **L** to read the lesson, **T** for the task, and **S** for your score. Press **?** to see the available CLI commands with the selected lesson ID. **M** returns to the course map.

In the reader, **← / →** changes sections and **↑ / ↓** scrolls. Press **b** to toggle a bookmark and **B** to see saved sections. Open the last-read section with **g** on the map or run `cli-silly resume`.

On the task screen, press **H** for hints and **Enter** to open Bash. Do the task and run `cli-silly check`. Type `exit` to return to the score screen. **C** checks again from the UI. Each lab shows your latest and best score. The course map counts completed labs.

### Tools by workshop

| Workshop | Tools |
| --- | --- |
| Shell, streams, Bash | Standard Linux utilities and Bash. The JSON lab also needs `jq`. |
| Git | Git with `git switch` support. Exercises use local repositories and need no account. |
| Linux app debugging and capstone | Python 3 and `curl`. The exercise starts a service on a free loopback port. |
| Docker | Docker Engine, the Compose plugin, and `curl`. Confirm access to the daemon with `cli-silly doctor`. |
| Service operations | `tar` and OpenSSH's `ssh` client for their respective lessons. The systemd exercise uses a local journal and unit file. |

Docker labs use `python:3.12-alpine`. Pull it during setup with `docker pull python:3.12-alpine` to prepare for offline practice. The labs allocate separate Compose project names and bind published ports to loopback. Run `docker compose down` in the lab when finished. For the volume exercise, `docker compose down -v` also removes the exercise's saved data. Stop background API processes from their practice shell when you finish those labs.

## Use it from your own shell

```bash
cli-silly list
cli-silly read streams/log-report
cli-silly learn streams/log-report
cd "$(cli-silly path streams/log-report)"
cli-silly task
cli-silly hint 2
cli-silly check
cli-silly progress
```

Inside a lab, `task`, `hint`, `check`, `progress`, and `reset` identify the lesson from the workspace, including from a subdirectory. Outside the lab, use a lesson ID to target it. Run `cli-silly progress` outside a lab for all scores or `cli-silly continue` for the first unfinished lesson.

`check` returns 0 for completion, 1 for work remaining, and 2 for an environment or checker error. Add `--json` for structured results. Reopening a lab preserves your files. `reset` asks for confirmation, archives the previous attempt, and retains your records.

Workspaces are stored under `${XDG_DATA_HOME:-~/.local/share}/cli-silly/workspaces`. Grades are tracked under `${XDG_STATE_HOME:-~/.local/state}/cli-silly/progress`. Redirected output and `TERM=dumb` use plain output. Set `NO_COLOR=1` for a color-free full-screen UI.

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
