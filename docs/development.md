# Local development

Use a Linux checkout with Bash. Build tooling and learner dependencies are separate. The CLI and its bundled materials can be used offline after installation.

| Work | Prerequisites |
| --- | --- |
| Build the CLI | Go 1.26 or newer, Git, Bash, and GNU Make for the provided targets |
| Build the website | Node 22.12 or newer on the Node 22 line in `.node-version`, npm, and access to the locked packages on the first install |
| Run the CLI checks | Bash, Git, Python 3, curl, jq, tar, ps, Node.js, and the OpenSSH client. The API labs bind to ephemeral loopback ports |
| Exercise Docker labs | Docker Engine, the Compose plugin, an accessible daemon, and the lab image available locally or from the registry |
| Package releases | Go, Bash, tar, sha256sum, and realpath |

## Fresh checkout

From the repository root:

```bash
make setup
make cli
make doctor
make run-dev
```

Use `make cli` by itself for CLI work. It fetches Go dependencies as needed. For website-only work, run `make web-deps` and `make web-dev`. The website build consumes bundled course files and the shared installation guide in `docs/install.md`. It needs no server credentials.

`make doctor` shows whether each workshop has the tools it needs. A Docker setup problem affects its workshop. Start with shell, streams, or Git while Docker is being set up.

## Files and cleanup

| Path | Purpose |
| --- | --- |
| `bin/` | Local CLI build |
| `node_modules/` | Packages installed by `npm ci` |
| `dist/` | Vite website output |
| `artifacts/releases/` | Linux archives and checksums |
| `.local/data/cli-silly/` | Development workspaces created with `make run-dev` |
| `.local/state/cli-silly/` | Development scores and locks created with `make run-dev` |
| `${XDG_DATA_HOME:-$HOME/.local/share}/cli-silly/` | Normal learner workspaces |
| `${XDG_STATE_HOME:-$HOME/.local/state}/cli-silly/` | Normal learner scores and locks |

Run `make clean` to clear generated binaries, web output, and release archives. Rebuild with `make cli`, `make web-check`, or `make package VERSION=0.1.0`. Use `make web-deps` to refresh dependencies from the committed lockfile.

Treat `.local/` and the normal XDG directories as learner data. Use an isolated test run for smoke tests. Use the CLI's targeted reset to start a specific exercise over. Stop background lab processes and run `docker compose down` from the Compose lab before archiving an active workspace.

## Verification

| Change | Run |
| --- | --- |
| Go commands, course validation, rubric or persistence | `make go-check` |
| Terminal keys, routing, rendering, or shell handoff | `make terminal-check` |
| Companion website or shared published lessons | `make web-deps web-check` and load the affected routes in a browser. Verify `/install`, the affected `/workshop/…` page, and legacy links when changing navigation |
| Compose lifecycle or grading | `make docker-check` |
| Packaging or startup path | `make package VERSION=0.1.0` and `bash scripts/test-install.sh 0.1.0` |
| Full checkout | `make check` after setup, then the live Docker check on a host with a daemon |

The terminal smoke test creates isolated state in a temporary directory. It exercises Bash handoff and verifies a Markdown table, saved bookmarks, and reading resume across separate processes. The Go suite covers starter, partial, and completed scenarios. Read verbose Go output if you need to verify which optional cases were skipped. The default suite checks Compose file syntax if Compose is installed. The live Docker target starts and cleans up its own projects.

If the build environment restricts network access, populate dependency caches through its approved mechanism. Local API tests require only loopback access. A missing daemon can still allow Compose file validation. Record that distinction when reporting results.

After building the companion, run a production preview in one terminal and the browser check in another:

```bash
npm run preview -- --host 127.0.0.1 --port 4193 --strictPort
# In the second terminal:
CLI_SILLY_TEST_BROWSER=brave-browser node scripts/test-browser.mjs http://127.0.0.1:4193
```

Use an installed Chromium-family executable for `CLI_SILLY_TEST_BROWSER`. The script creates an isolated profile under ignored `.local/browser-checks/` and closes its browser afterward. It checks every workshop and legacy route, follow-up links, search, drill disclosure, copy feedback, installation, and mobile navigation. Pass the public site origin to verify a deployment. Generated profiles and release output are excluded from linting.

### Dependency review for v0.6.0

The locked dependency update removes the high- and moderate-severity findings. The typography package's nested selector parser is overridden to patched 6.1.5; remove the override when its upstream dependency is patched. See the [selector parser advisory](https://github.com/advisories/GHSA-rj75-hqrm-r3gf).

The 2026-11-28 audit still reports two low-severity entries for one [KaTeX prototype-pollution gadget](https://github.com/advisories/GHSA-238p-pmpm-9mq7), through Mermaid. Exploitation requires existing prototype pollution and attacker-controlled expressions. This companion renders repository-authored content and has no user-submitted diagram input; Mermaid uses strict rendering. This limits exposure but does not make the installed dependency patched. npm proposes a major Mermaid upgrade, which needs a separate diagram compatibility review; no forced major update is part of v0.6.0. Re-run `npm audit` as advisories change. Mermaid remains lazy-loaded; its large optional chunks are not part of the homepage's initial download.

## Release rehearsal

```bash
make package VERSION=0.1.0
bash scripts/test-install.sh 0.1.0
```

The test verifies the archive checksums, installs into a temporary directory, completes a lab, and exercises upgrade, uninstall, and reinstall. Run it on arm64 hardware or a suitable emulator to establish arm64 runtime coverage.

The release workflow builds assets from a version tag such as `v0.1.0` and tests the same archives on native amd64 and arm64 runners before publishing them on GitHub Releases. See the [installation guide](install.md) for the user-facing flow.

## Product completion checks

The product target remains native Linux amd64/arm64. `make go-check` includes `make course-check`, which validates published teaching structure without running commands from Markdown. Keep draft lessons under `.local/lesson-drafts/` until they have real grader coverage. See `CONTRIBUTING.md` for review and `docs/learner-sessions.md` for human testing. Usability and screen-reader claims need recorded sessions; automated PTY checks are separate evidence.
