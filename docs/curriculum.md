# Curriculum

Each built-in lesson has an offline reading, a worked example, and a scored exercise. Read explains what the commands do and how to inspect their effects. Do opens a lab with starter files, hints, and a rubric. The [roadmap](../plan.md) tracks changes to the platform.

## Foundations

| Companion topic | CLI lesson ID | Practice |
| --- | --- | --- |
| CLI & Linux Intro | `foundations/cli-intro` | Inspect a project and save a note and directory listing. |
| Mental Model | `foundations/mental-model` | Capture stdout and the exit status of a failing command. |
| CLI & Linux Core Terms | `foundations/core-terms` | Choose options and operands to preview an inventory. |
| Filesystem & Navigation | `foundations/filesystem-navigation` | Navigate sibling directories and save reports at the lab root. |
| Users, Groups & Identity | `foundations/users-and-groups` | Resolve an account's UID and primary group using captured records. |
| Permissions & Access Control | `shell/make-it-runnable` | Enable direct script execution and capture its output. |
| File Operations | `shell/paths-and-quoting` | Move a configuration file whose path includes spaces. |
| I/O, Redirection & Pipelines | `streams/split-streams` | Save standard output, diagnostics, and status separately. |
| Text Manipulation | `streams/log-report` | Filter server errors and count them by route. |

## Core tools

| Companion topic | CLI lesson ID | Practice |
| --- | --- | --- |
| Finding Things | `streams/find-config` | Locate a setting by path and exact matching line. |
| Environment Variables & Shell State | `core-tools/environment-variables` | Change a child process's environment in two ways. |
| Package Management | `core-tools/package-management` | Use offline package inventories to find a version and an executable's owner. |
| Networking from the CLI | `linux/inspect-local-api` | Debug an HTTP service running on loopback. |
| SSH | `operations/ssh-config` | Verify effective settings with `ssh -G` using a local config. |
| Processes, Jobs & Signals | `core-tools/processes-and-signals` | Inspect a child worker by PID and let it handle TERM. |
| Shell Scripting | `bash/arguments-and-exit` | Write a script that validates a file argument and distinguishes success from failure. |
| Git from the CLI | `git/selective-staging` | Commit only the intended change in a local repository. |
| Archiving & Compression | `operations/disk-and-archives` | Create, list, and verify a small compressed archive. |
| Advanced Text Processing | `core-tools/advanced-text-processing` | Compare a JSON inventory with a plain-text event log. |
| Node.js CLI Development | `core-tools/node-cli` | Build a command with distinct stdout, stderr, and exit codes. |

## Follow-on labs

After the two main groups, the catalog offers structured JSON, Git branching and remotes, Bash report automation, Docker and Compose, service operations, and release rescue. The website also has the original Systems & operations reference pages.

Press **p** in Do for follow-up challenges linked to the current lesson. Challenges use different starter files and receive their own score. The catalog's **Review challenges** workshop lists them together.

| Challenge | Practice focus |
| --- | --- |
| Trace a worker incident | Filter error events and look up ownership in JSON. |
| Publish a release note safely | Move one quoted filename and preserve the neighboring drafts. |
| Build a retry queue from a failed run | Separate streams, filter retry messages, and preserve the exit status. |
| Compare isolated app settings | Test inherited defaults and per-command environment overrides. |
| Package one build log | Create a one-member archive and verify the member name and contents. |
| Extract production deployment targets | Filter JSON, build a TSV, and sort distinct ports numerically. |

Use `cli-silly doctor` to see optional prerequisites for practice. Identity and package exercises use lab files, SSH stays local, and the Node exercise uses built-in libraries with no npm install. Read remains available offline even when a tool needed for Do is missing.

## Complete Foundations practice

All nine Foundations lessons now have linked scored practice. The additional tasks cover project handoffs, preserving failure status, changed inventories, repairing misplaced reports, resolving account/group records, and restoring script execution. Each uses independent starter files and its own grade. The Git release follow-up combines committed configuration, release notes, and publishing to a local remote.

The built-in course has 30 core labs and 13 follow-ups. Five connected projects in `course/projects.json` organize incident investigation, Git release repair, local API recovery, Compose operations, and Bash release preparation. Open them through Ctrl+k → Connected developer projects. Projects reuse lesson identities and do not duplicate grades.
