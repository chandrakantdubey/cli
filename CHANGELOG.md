# Release notes

## 0.6.0 — Linux learning and recovery

- Added a welcome/resume screen and searchable Ctrl+k action menu while keeping Read and Do as the primary modes.
- Expanded all 30 core lessons with outcomes, prerequisites, common mistakes, and development use cases. Corrected the release-report lesson's contradictory ordering advice.
- Added seven independently scored follow-ups. The course now has 43 exercises, including related practice for every Foundations lesson.
- Added five connected developer projects for incidents, Git releases, API recovery, Compose operations, and Bash release preparation.
- Improved feedback for empty, misplaced, differently ordered, and incorrectly formatted reports. The intro and archive follow-up now verify required starter preservation; older rubric scores remain available separately.
- Added spaced review intervals that grow after due successful attempts. Repeated checks in one session do not postpone review.
- Added read-only file/Git inspection, archived-attempt restoration, and bounded local backup export/import. Import requires empty destinations and preserves existing learner data.
- Improved Linux dependency guidance and added static ASCII line art alongside existing no-color and reduced-motion options.
- Exposed terminal drills on the companion site, added lesson search and copyable practice commands, and preserved existing lesson URLs.
- Added course quality checks, lesson scaffolding, contributor guidance, and a learner-session protocol.
- Updated locked website dependencies to remove high- and moderate-severity audit findings, and enabled strict diagram rendering. Two low-severity KaTeX/Mermaid audit entries remain; see the dependency review in `docs/development.md`.

Supported application targets remain native Linux amd64 and arm64. The repository remains unlicensed at the owner's request. Automated checks do not replace the still-pending first-time learner and screen-reader sessions tracked in `plan.md`.

## 0.5.0

Added replayable worked-example motion and five follow-up exercises. See the [published release](https://github.com/chandrakantdubey/cli/releases/tag/v0.5.0) for artifacts and history.
