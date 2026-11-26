# Curriculum

The course serves developers who need to navigate a Linux project, run and diagnose an app, ship changes, and operate a small service. The course contains scored labs that run in local workspaces. Application, Docker, and release labs use their own starter projects and fixtures.

| Workshop | Related web references | Learner outcome |
| --- | --- | --- |
| Shell and files | 1–7, 23 | Navigate paths, quote filenames, use patterns, and repair a script's permissions. |
| Streams and data | 8–10, 19 | Combine search, streams, text tools, and JSON to produce a report. |
| Linux app debugging | 11–13, 15, 26 | Follow configuration into a process and diagnose a failed localhost request. |
| Git | 17 | Stage selectively, merge a feature, resolve a conflict, and push to a local remote. |
| Bash automation | 16, 31, 33 | Write a repeatable script that handles arguments and failure. |
| Docker and Compose | 28, 32 | Build an image, connect services, and prove data persists. |
| Service operations | 14, 18, 24, 25, 27, 29 | Diagnose a service startup failure using config, journals, and permissions. |
| Release rescue | 30, 34 | Restore a local release and record the evidence of recovery. |

The 34 original pages are the companion reference library. Modules 20–22 are available as references. Review the corresponding page and preserve its URL when promoting reference material into a scored lab.

Platform changes are tracked in the [roadmap](../plan.md). This document tracks the built-in developer course.

## Built-in practice sequence

- **Shell and files** covers script modes and quoting paths.
- **Search, streams, and data** covers recursive search, pipelines, structured JSON, and separating stdout from stderr.
- **Linux application debugging** runs a local API and checks a real HTTP response.
- **Git** covers selective staging, merging, a local remote, and conflict resolution.
- **Bash automation** checks argument handling, exit codes, and report generation on new fixtures.
- **Docker and Compose** covers images, host ports, health checks, named volumes, and service-name networking.
- **Service operations** uses unit and journal fixtures, tar archives, and effective SSH settings.
- **Release rescue** asks the learner to diagnose and correct a broken local API release.

The optional worker-incident challenge revisits status filtering and counting from the log-report lab, then adds an ownership lookup from JSON. It has its own workspace and grade.
