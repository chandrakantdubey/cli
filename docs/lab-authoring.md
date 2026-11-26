# Lab authoring

1. State one observable skill. Teach it with a worked example that uses different filenames or data from the learner task.
2. Give the learner ordinary starter files they can inspect and edit. Declare dependencies and set intended file modes during setup.
3. Write the task and scoring criteria before writing the checker. Criteria should sum to 100 and describe outcomes.
4. Give each check useful feedback for the common wrong state. Hints should move from observation to a concrete next action.
5. Verify an untouched starter, a partial solution, a complete solution, and an alternate valid result where applicable.
6. Change the grading revision if expectations or weights change. Change the starter revision when the starter state changes. Copy edits can retain both revisions.

Use the shell permissions lesson as the reference for a file-based checker. Its checker reads a regular file and compares report bytes. Tasks that assess script behavior should use multiple inputs and document when the checker runs learner code. Git and Docker checks should scope inspection to the disposable repository or lab project.

## Catalog and checkers

Add new lesson IDs to `course/catalog.json` in teaching order. The CLI map and website sidebar load titles, tasks, and criteria from that catalog. Declare additional tools in the lesson manifest and include the corresponding starter layout. Add an untouched, partial, and completed scenario to `internal/lab/course_test.go`.

File, Git, archive, HTTP, script, and Compose checks evaluate different observable results. Provide a criterion-specific feedback message so a learner sees the relevant command or file. Script checks run the learner's script under their current user in a temporary directory with alternate fixtures. Explain this in the lesson task. Compose checks scope inspection to the project's generated name. Document how to stop a live service and clean up disposable resources.

Use the [roadmap](../plan.md) for a new workshop or changes to the lesson or checker contract. Review the [teaching guide](teaching.md) when editing lesson prose or demos.

Reference examples live in `course/references.json`. Use declared `{{name}}` parameters as shell arguments. The application quotes values and previews the result. Templates can be written inside the active lab at a new relative path.
