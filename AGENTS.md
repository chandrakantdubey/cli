# Working on cli-silly

cli-silly is a terminal-first learning platform. Its built-in course teaches Linux, Git, Bash, and Docker through offline reading, worked examples, real Bash workspaces, and scored feedback. The companion website reads the same published material.

## Where to work

- `cmd/cli-silly/` parses commands and connects the CLI to the course, lab, and terminal packages.
- `internal/tui/` handles navigation, lesson presentation, keyboard input, and Bash handoff.
- `internal/lab/` owns workspace lifecycle, checker execution, and saved progress.
- `course/` holds the catalog, lesson manifests, Markdown, and starter assets. Both products consume it. The companion also renders `docs/install.md` at `/install`.
- `src/` contains the companion site and legacy reader. Keep existing lesson URLs working during migration.
- `Makefile` and `scripts/` provide the local build, verification, and packaging commands used by CI.

## Setup and checks

- Use Go 1.26 or later for the CLI. Website development uses Node 22.12 or later on the Node 22 line declared in `.node-version`. See [local development](docs/development.md) for the tool and lab requirements.
- Run `make setup` for a full checkout. Use `make cli` for the CLI and `make web-deps` for website-only setup.
- Use `make run-dev` for interactive development with checkout-local workspaces and progress. The normal CLI uses the learner's XDG directories.
- Run `make go-check` for CLI and course changes, `make terminal-check` for interactive changes, and `make web-deps web-check` for website changes. A rubric change needs starter, partial, and passing cases. Use `make docker-check` to verify live Compose behavior when a daemon is available.
- For installer or packaging changes, run `make package VERSION=0.1.0` and `bash scripts/test-install.sh 0.1.0`. See [installing cli-silly](docs/install.md) for the release contract.
- Report skipped checks and their reason. Generated files belong in ignored output directories. Keep authored starter logs, dotfiles, and dependency lockfiles in source control.

## Product and compatibility rules

- Use [plan.md](plan.md) as the source of truth for work order and status. Add or update a roadmap item with an outcome, acceptance criteria, and verification before starting a feature. Update the item when work lands and record pending checks. The [product direction](docs/product-direction.md) describes the intended learner experience.
- Keep lesson IDs and URLs stable. Preserve existing workspaces, progress records, and public command behavior during schema or UX changes.
- Explain the skill and rubric before opening the workspace. Use a real tool and score observable results. Document tool requirements, hints, and feedback with the lesson.
- Source lesson text, demo frames, and criteria from `course/`. Review both surfaces when shared content changes.
- Keep setup, validation, reset, and cleanup scoped to the managed lab. Resets archive work and preserve earned progress. State which checks run learner-authored scripts.
- Keep UI, parsing, grading, and persistence responsibilities in their current boundaries. Record an architecture decision when changing a shared contract.

See [architecture](docs/architecture.md), [curriculum](docs/curriculum.md), [teaching](docs/teaching.md), [lab authoring](docs/lab-authoring.md), and [terminal UI](docs/terminal-ui.md).
