# Contributing

Read [AGENTS.md](AGENTS.md), [the roadmap](plan.md), and [local development](docs/development.md). The product supports native Linux amd64 and arm64. Record an outcome, acceptance criteria, and verification in the roadmap before implementing a feature.

Keep UI, grading, content, and persistence in their existing packages. Shared lessons belong in `course/`. Keep lesson IDs and URLs stable. Declare schema changes and test old progress/workspaces before release.

For content, use `scripts/new-lesson.py` and the [authoring guide](docs/lab-authoring.md). Add meaningful starter, partial, passing, alternate-valid, and failure coverage. Do not publish draft placeholders or claim that output inspection proves which command was used.

Run `make go-check`, `make terminal-check`, and `make web-deps web-check` for the affected surfaces. Use `make docker-check` with a daemon for live Compose work. Report unavailable checks honestly. Use checkout-local XDG directories and keep learner files untouched.

The repository is currently unlicensed at the owner's request. No open-source reuse license is granted. A future license requires the owner's explicit decision; do not add one as a routine cleanup.

Release through the existing GitHub workflow after checks pass. Verify the published Linux installers, upgrade behavior, and companion routes. Record evidence in `plan.md`; human learner sessions have their own gate.
