# cli-silly development roadmap

The [shared design discussion](https://chatgpt.com/share/6b07f372-6888-83ee-ac0e-27ed1001fe1c) sets the product direction. The first screen should show the learner what they can study, what the selected task looks like, and which keys or commands will let them act. cli-silly is the platform. The existing developer CLI course is its first course.

This file is the source of truth for planned work, delivery order, and completion status. The [architecture](docs/architecture.md) describes code and data boundaries. The [development guide](docs/development.md) covers contributor setup and the [installation guide](docs/install.md) covers learner installation.

## Status at a glance

Checked boxes record completed work with evidence. An implementation and its environment-specific verification can have separate entries. Record the remaining test and the environment it needs when verification is pending.

- [x] **BASE-01** Establish the local setup, shared Make targets, isolated development state, Git ignore rules, and CI checks. Verified a clean local rebuild and confirmed that cleanup preserves course files and learner data.
- [x] **COURSE-01** Ship 20 built-in labs across eight workshops with offline reading, worked examples, starter workspaces, hints, and scored rubrics. Verified the 17 non-Docker labs with starter, partial, and completed attempts.
- [x] **TUI-01** Provide the searchable split-pane catalog, lesson preview, keyboard help, Bash handoff, and score views. Verified search, reading, practice, and saved scores in a pseudo-terminal.
- [x] **WEB-01** Connect the companion site to the built-in catalog and preserve the legacy reader routes. Verified all 20 companion routes and the legacy redirect in a local browser. Public deployment is tracked under REL-03.
- [x] **WEB-02** Align the companion with the installed CLI. Share the Linux install guide with the site and release archives, and show the installed practice command on each lesson. Verified install navigation and reload, all 20 lab commands, and the shared guide inside both Linux archives.
- [x] **REL-01** Implement Linux packaging, checksums, and per-user install, upgrade, and uninstall. Built amd64 and arm64 archives and verified the complete install lifecycle from the packaged amd64 binary, including progress after reinstall.
- [x] **REL-02** Add a version-tagged GitHub Releases workflow that builds and tests assets before publishing.
- [x] **REL-03** Publish and verify [v0.1.1](https://github.com/chandrakantdubey/cli/releases/tag/v0.1.1). Both native installer jobs passed. The public downloads passed checksum, installation, lab, and offline course-pack checks. Verified the [production install page](https://cli-eight-phi.vercel.app/install), the core and review lab commands, and the legacy reader routes.
- [x] **REL-04** Verify the installed CLI on native Linux arm64. The [ARM64 installer job passed](https://github.com/chandrakantdubey/cli/actions/runs/45799569549). Tagged releases test the same candidate artifacts on both architectures before publication.
- [x] **REL-06** Add a single GitHub installer URL. Verified latest-release discovery, architecture selection, checksum verification, installation to a temporary user bin directory, and offline reading from the installed binary.
- [x] **REL-05** Verify all three live Compose labs. The GitHub Actions Docker grading step [passed](https://github.com/chandrakantdubey/cli/actions/runs/45796479619).
- [x] **READ-01** Build the rich terminal reader, section bookmarks, and resume. Verified formatted tables and fences at narrow and wide widths, color-free reading, saved state across launches, the packaged CLI, and companion rendering. The plain `read` command and first-unfinished behavior of `continue` are preserved.
- [x] **LAB-01** Add contextual editing, workspace feedback, and opt-in file validation. Verified with real terminal editor handoff, live score changes, concurrent checks, and reset/cancellation tests.
- [x] **REVIEW-01** Add quizzes, an independent log-analysis challenge, and a review queue. Verified quiz feedback and challenge navigation in a pseudo-terminal, review scheduling with a fixed test clock, and starter, partial, and complete challenge grading.
- [x] **REF-01** Add searchable snippets and starter templates. Verified terminal browsing, command output with quoted parameters, workspace-local file generation, and collision handling.
- [x] **PACK-01** Prove offline installation and upgrades with the Markdown Basics pack. Tested the installed binary, grading, preserved progress after upgrading, incompatible manifests, and invalid archive paths.
- [x] **POLISH-01** Rehearse full learner onboarding and polish the companion site. Verified installed terminal practice through the review challenge and all 21 companion lesson routes, including the link back to the original lesson.
- [x] **WEB-03** Add the terminal-prompt brand, explain commands and operations in every built-in lab, and load the content for the active route. Verified all 21 workshop lessons and 34 reference routes in a browser. With the cache disabled, gzip-compressed assets fell from 213 kB to 108 kB on the homepage and from 313 kB to 157 kB on the first shell lesson in local production builds.
- [ ] **POLISH-02** Observe a first-time learner completing the install, practice, and review flow. Record feedback and fixes.
- [x] **COURSE-02** Bring all 20 companion Foundations and Core tools topics into the offline CLI.
  - Outcome: Learners can choose these topics in the same order and read command explanations, see a worked example, and do a scored local exercise for each.
  - Acceptance: Read and Do are the two primary lesson actions. All existing lesson IDs, reader URLs, and saved grades remain usable. The companion links each source topic to the matching offline lesson.
  - Verification: Passed starter, partial, and complete grading for the ten new labs and the Read and Do flow in a pseudo-terminal. Browser checks passed for all 31 workshop pages and 34 legacy URLs.
- [x] **REL-07** Deliver the Foundations and Core tools course through the direct installer.
  - Acceptance: The new GitHub release includes Linux amd64 and arm64 archives. Upgrading preserves an existing workspace and saved score.
  - Verification: Test the native archive locally and on both GitHub runners, then run the published installer with isolated learner state.
  - Evidence: [v0.2.0](https://github.com/chandrakantdubey/cli/releases/tag/v0.2.0) passed installation on [native amd64 and arm64](https://github.com/chandrakantdubey/cli/actions/runs/45947538125). The public one-line installer upgraded v0.1.1 with its workspace and score intact. The published binary read the Node lesson offline and graded the new intro lab at 100/100.
- [x] **TUI-02** Make Read and Do a guided terminal learning flow.
  - Acceptance: Read opens the full companion lesson as a scrollable document, with a visible section index and additional hands-on notes for all 20 Foundations and Core tools lessons. Do explains where reports belong, carries the correct lesson into Bash, and gives clear feedback on return.
  - Verification: Go and terminal checks passed for continuous reading, contents navigation, resume, narrow and color-free layouts, and a full intro attempt with a misplaced report, partial credit, and automatic scoring on return. Verified all 20 embedded chapters and the companion routes.
- [x] **REL-08** Publish the guided reader and practice flow in v0.3.0.
  - Acceptance: The GitHub installer upgrades v0.2.0 and preserves an existing lab and its best score.
  - Verification: [Native install jobs passed](https://github.com/chandrakantdubey/cli/actions/runs/45950944250). The public one-line installer upgraded v0.2.0 to [v0.3.0](https://github.com/chandrakantdubey/cli/releases/tag/v0.3.0), preserving the intro workspace and its 100/100 score. Reading from inside the workspace opened the correct chapter and terminal drills.
- [x] **TUI-03** Make lesson navigation consistent from the catalog through practice.
  - Outcome: Learners can move between topics and switch Read and Do without losing their place.
  - Acceptance: Tab switches Read and Do, brackets move between lessons, and Escape closes the current view or overlay. The catalog shows progress and workshop shortcuts. Read, worked examples, practice, and results keep their next actions visible at narrow and wide widths.
  - Verification: Go and pseudo-terminal tests passed for workshop jumps, returning from search and examples, preserving a reading position, switching Read and Do, and Bash handoff with scoring. Reviewed the color layout in a real terminal and verified the narrow layout.
- [x] **REL-09** Publish the navigation update in v0.4.0.
  - Verification: [Both native installers passed](https://github.com/chandrakantdubey/cli/actions/runs/45960488705). The public installer upgraded v0.3.0 to [v0.4.0](https://github.com/chandrakantdubey/cli/releases/tag/v0.4.0) and preserved the intro workspace and its 100/100 score. The companion install page shows the new controls.

## Scope

Keep the built-in course, its IDs, and the current public reader routes throughout the migration. The terminal is the primary practice interface. The companion site is a public handbook with exact CLI entry points. Use local content and progress for the first platform release. A second small course should prove the package boundary before the project commits to a large new curriculum.

### Installation paths

| Learner starts with | Supported path |
| --- | --- |
| A clone or an extracted GitHub source ZIP | Run `GOBIN="$HOME/.local/bin" go install ./cmd/cli-silly` from the source root with Go 1.26 or newer. |
| A prebuilt release | Put the installer, matching checksum file, and architecture-specific archive together. Run `bash install.sh <archive>`. |
| A direct GitHub release link | Download the same three files from the versioned release URLs and run the installer. The [installation guide](docs/install.md) includes the terminal commands. |

The installer defaults to `~/.local/bin`. Once on PATH, the learner runs `cli-silly doctor` and `cli-silly`. The CLI bundles the course. Tools such as Git, jq, and Docker are workshop prerequisites. Upgrades and uninstall keep learner workspaces and progress.

### Product decisions

| Decision | Initial direction |
| --- | --- |
| Catalog | Preserve the two-pane path and preview design. On smaller terminals keep the same actions in a stacked view. |
| Teaching | Each lesson teaches an observable skill, demonstrates a real command and state change, and offers practice. |
| Progress | Show per-criterion feedback, current and best score, and workshop completion. Track quiz answers separately from lab points. |
| Navigation | Use Read and Do as the two lesson actions. Preserve `/` for search, `?` for help, `m` for the catalog, and the existing Bash handoff. |
| Tools | Keep Go, Bubble Tea, Lip Gloss, and React/Vite. Evaluate a terminal Markdown renderer against representative lessons during the reader slice. |
| Storage | Keep current XDG locations. Add versioned records only when an implemented feature needs them. |
| Course distribution | Start with a built-in course and manual local installation of a second declarative pack. Add discovery and online distribution after the format works end to end. |

The interface should show the current lesson, prerequisite tools, an example with actual input and output, the task and rubric, and a visible next action. For every check, give the learner the observed result and what to inspect next. The map should show which workshops they can already practice when optional tools are missing.

## Delivery sequence

### 0. Release the verified baseline · REL-03–REL-05

1. Verify the installer and checksum files from a clean download directory on Linux amd64. Use the exact filenames and commands in the installation guide.
2. Publish through a version tag after the release content is committed. Check that the published assets and direct links match the guide. Verify installation again from the published release.
3. Test the arm64 archive on native hardware or an emulator before claiming arm64 runtime support. Record the tested environment.
4. Run the Docker lab suite against a working daemon. Verify the three labs, including scoped setup, grading, and cleanup.
5. Check the public companion site's legacy routes and its links to installation and the corresponding practice command.

**Verification:** A learner without the source checkout or a Go toolchain can install the published build, read a lesson offline, finish a starter lab, and reopen with the same score. Report platform and lab coverage separately.

### 1. Reader and lesson navigation · READ-01

1. Add a reader component that renders headings, lists, inline code, fenced commands, and tables within the viewport. Keep the Go lesson body as source and preserve the plain output of `read`.
2. Keep worked frames distinct from prose and show the actual command, output or diagram, and explanation. Show keyboard actions consistently in course map, reading, demo, and task screens.
3. Add a resume pointer and bookmarks. Store a stable lesson ID and section anchor. An unavailable anchor can fall back to the lesson overview. Keep grade records independent.
4. Update the companion's matching lessons and keep direct links, refresh, and the legacy redirect working.

**Verification:** Start at a fresh terminal, find a lesson by its command, read formatted Markdown at narrow and wide widths, bookmark a section, restart, and resume. Check keyboard-only and color-free use. A fixture lesson with tables, code, and long commands should render in both interfaces.

### 2. Practice workspace and feedback · LAB-01

1. Add a contextual editor action for declared starter files. Open the selected file through the learner's configured editor, using the same terminal release pattern as Bash. Keep the file within the managed workspace.
2. Surface the workspace path, relevant tools, full practice objective, hint controls, and the commands that open and grade the lab.
3. Bring learner results back to the current lesson. Keep manual checks usable from the TUI and shell. Add opt-in watch validation for file-based labs after cancellation and duplicate-save behavior are covered. Runtime and script checks can remain explicit actions.
4. Make reset and runtime cleanup understandable. A reset archives the prior attempt. The learner can see which project or process a lab owns and how to stop it.

**Verification:** Edit a named exercise file, earn partial credit, repair the output, and reach completion. Confirm shell and editor handoff, simultaneous shell and TUI checks, cancellation, missing tools, and reset while using checkout-local state.

### 3. Skill checks, challenges, and review · REVIEW-01

1. Define a declarative quiz item with a stable ID, prompt, choices or expected answer, and explanation. Show answer explanations after the learner submits an answer or requests feedback.
2. Add an independent task or fixture variant for one completed skill. Use explicit lesson and variant IDs, separate workspaces, and a rubric tied to the target outcome.
3. Track completion and review dates in a separate local record. Show due reviews and missed topics. Start with a simple review interval and adjust using learner sessions.
4. Add the first cumulative workshop mission. Use prior CLI skills in a realistic scenario with several independently graded outcomes.

**Verification:** The learner can answer a quiz, understand the explanation, find a due review after restarting, and complete a changed challenge. A review attempt leaves the original grade intact. Tests control time and verify repeat attempts, missing variants, and resets.

### 4. Snippets and starter templates · REF-01

1. Add searchable reference records with a stable ID, title, related tools, parameters, example input, expected output, and related lesson links.
2. Support terminal browsing and an explicit copy or print operation. Provide normal text output for piping. Show parameter substitutions before emitting the command.
3. Add file generation for authored templates. Resolve destinations inside the selected workspace. A collision should present the destination and let the learner choose a new path or explicitly replace it.

**Verification:** Find a command by problem or tool, review what it will do, obtain the snippet through both interfaces, and generate a lab file. Cover paths with spaces, invalid parameters, and overwrite behavior.

### 5. Installable courses · PACK-01

1. Define a versioned package manifest, course namespace, content version, engine compatibility, lesson registry, typed criteria, and relative asset paths. Reuse the catalog, reader, TUI, and grader.
2. Validate a local pack before making it visible. Reject duplicate IDs, missing entries, invalid references, unsupported check types, paths outside the pack, and incompatible versions. Limit archive size and extraction scope.
3. Start with declarative content and checker types the app implements. Run only documented, explicitly invoked actions. Define an explicit trust and execution model before adding third-party checker code.
4. Keep the built-in IDs valid. Namespace installed course IDs and progress. Persist the selected course and last activity. Install or upgrade atomically, with a clear recovery path if validation fails.
5. Deliver a tiny second course as the integration fixture. Use a lesson with a different skill and an ordinary file to show that no subject-specific code is needed.

**Verification:** Install the fixture pack offline, read and complete it, restart, upgrade the pack, and see preserved compatible grades. Exercise a duplicate namespace, broken manifest, failed upgrade, malicious archive path, and existing state created by the current built-in course.

### 6. Release and companion polish · POLISH-01

1. Generate the companion's public catalog from published course metadata. Preserve existing deep links and ensure each practice command is available in the matching CLI release.
2. Rehearse clean Linux installation and offline reading. Test all advertised architecture/runtime combinations and live Docker grading. Verify artifacts and checksums.
3. Observe a learner starting from install, locating a relevant lesson, reading it, completing a task, and returning for a review. Fix problems that prevent them from finding the next action.

**Verification:** The website and installation instructions agree with the packaged course, stable links work, the binary completes the onboarding path without an account, and checks distinguish current state from historical progress.

## How new work enters the roadmap

Add or update a roadmap item before implementing a new feature or changing a public contract. Use the closest existing item when it already covers the work. Routine fixes can reference the item they support. Contributor docs explain how to build or verify a change. Keep priorities and completion status here.

Use this format for a new item:

```markdown
- [ ] **AREA-NN** Feature name
  - Outcome: What the learner or contributor can do.
  - Acceptance: Observable behavior that demonstrates the outcome.
  - Dependencies: Other items, decisions, or environments this needs.
  - Verification: Commands or learner flow to exercise.
```

Place it in the appropriate phase, or in the backlog if its priority is undecided. When work starts, record the agreed scope. When it lands, mark completed work, link to the implementation or test evidence where useful, and keep any open verification as a separate checkbox. When priorities change, move the item instead of copying it to a second list.

### Backlog

Add new ideas here until there is a clear outcome and a place in the delivery sequence.

## Shared contracts and dependencies

| Contract | Before changing it |
| --- | --- |
| Lesson manifest and checker kinds | Add validation cases, preserve existing IDs, and confirm both consumers resolve the same content. |
| Workspace marker | Test opening current workspaces, reset and rollback, stale shell context, and a workspace containing dotfiles or spaces. |
| Score and quiz storage | Provide a migration or a compatible read path. Test a legacy record and ensure a fresh write preserves the earlier score. |
| Shell and editor handoff | Verify the terminal restores state after success, error, and interruption. Confirm the process starts in the selected workspace. |
| CLI output and routes | Characterize documented exit codes, JSON output, direct commands, browser reload, and the legacy redirect. |
| Release boundary | Record the currently supported platforms, optional external tools, and any checks skipped on the build machine. |

Use the existing lab suite for terminal-independent grading. Add focused Go tests for new model transitions and storage. Use pseudo-terminal tests for keyboard navigation and handoff. Add browser tests for routes when the companion lifecycle changes. A course package change should exercise the loader against both valid and rejected content, then exercise at least one valid pack through the normal learner entry points.

## Decisions to confirm during implementation

- Choose the Markdown renderer after checking rendering, terminal resizing, dependency weight, and color-free output.
- Choose the first review schedule and second course after the reading and practice flow has been tried by learners.
- Decide release distribution and licensing before publishing downloadable third-party course material. The local build and test workflow already packages the embedded course.
