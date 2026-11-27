# Terminal learning platform

The [shared design discussion](https://chatgpt.com/share/6b07f372-6888-83ee-ac0e-27ed1001fe1c) describes a reusable terminal learning platform. The accompanying terminal mockup shows learning paths on the left, an exercise and its checks on the right, and keyboard controls along the bottom.

For cli-silly, the built-in developer course is the first course on that platform. Keep its existing lesson IDs and progress. Each course should provide the same complete path from reading and seeing an example to practicing in a real workspace and getting feedback.

## Learning flow

1. Search learning paths by course topic, lesson name, or command.
2. Preview the selected lesson and see its task, a worked command, and prior score.
3. Choose **Read** for offline explanations and worked examples.
4. Choose **Do** to review the rubric, open Bash, and get scored feedback. Resume or retry with the best score preserved.

The catalog uses a two-pane layout on wider terminals and a stacked layout on smaller ones. It starts with the Foundations and Core tools topics from the companion. The right pane previews the selected lesson. Worked examples are part of Read. Bash, hints, and score feedback are part of Do. The companion presents the same reading and connects learners to the Linux install guide and the command that starts each lab.

## From the built-in course to a platform

The first course stays focused on the CLI skills full-stack developers use to build and operate applications. Grow the platform by making one learning path work end to end, then prove portability with a small second course. The [roadmap](../plan.md) defines the release sequence, prerequisites, and acceptance checks for the reader, playground, quizzes, snippets, review, and course packages.
