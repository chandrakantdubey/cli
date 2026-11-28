# Teaching a lesson

Open with the job the learner will be able to do. Explain the state the relevant command reads or changes. Show a short worked example and let the learner step through it at their own pace. Follow with a task in a real Bash workspace.

For each demo frame, include the command, the resulting state, and a sentence explaining the link between them. Use a diagram when it makes a change easier to see. In the permissions lesson the tree stays in place while the owner execute bit changes. In the Git lesson a commit should visibly take the snapshot from the index.

Make the task explicit about the deliverable and show how each criterion contributes to the score. When a check fails, point the learner to something they can inspect. Let them retry as often as they like. The walkthrough and task can use different filenames so learners practice applying the same idea.

Watch a first-time learner install the CLI, find the lesson, follow the example, complete the exercise, and resume after closing the CLI. Revise the lesson where they misread an instruction or cannot tell what to do with the feedback.

Use the same skill, example, and rubric on the terminal and companion site. Use descriptive, stable Markdown headings because bookmarks refer to those sections. Verify fenced commands and tables at a narrow terminal width. Keep installation steps in the [shared guide](install.md).

Use a **Commands and operations** section in each built-in workshop reading. Explain the command arguments, which files or state change, how to inspect the result, and what a common failure looks like. Reuse the manifest's worked example steps, which the companion shows together and the terminal advances one frame at a time.

## Release teaching contract

Every core lesson must include Before you practice (an outcome and prerequisite guidance), Commands and operations, Common mistakes and recovery, and Use it in development. `make course-check` verifies these sections alongside examples, hints, links, and scoring totals. This is a structural check; a human content review still checks correctness and whether the task teaches the intended skill.

Avoid equating a perfect report with command history or mastery. Teach multiple valid approaches and validate the resulting state. If a task asks learners to preserve starter inputs, enforce that requirement. Do not reduce earned scores for requesting a hint. Use the learner session guide in `docs/learner-sessions.md` to collect usability evidence.
