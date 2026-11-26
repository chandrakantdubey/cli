# Terminal learning platform

The [shared design discussion](https://chatgpt.com/share/6b07f372-6888-83ee-ac0e-27ed1001fe1c) describes a reusable terminal learning platform. The accompanying terminal mockup shows learning paths on the left, an exercise and its checks on the right, and keyboard controls along the bottom.

For cli-silly, the built-in developer course is the first course on that platform. Keep its existing lesson IDs and progress. Each course should provide the same complete path from reading and seeing an example to practicing in a real workspace and getting feedback.

## Learning flow

1. Search learning paths by course topic, lesson name, or command.
2. Preview the selected lesson and see its task, a worked command, and prior score.
3. Read or step through the lesson. Open the real workspace when ready to practice.
4. Check work and see the reason for each result. Resume or retry with the best score preserved.

The catalog now uses a two-pane layout on wider terminals and a stacked layout on smaller ones. The right pane previews the selected lesson. Reading, visual demonstrations, practice, and scoring remain accessible from the course map. The companion site presents the same lesson material and connects learners to the Linux install guide and the command that starts each lab.

## From the built-in course to a platform

The first course stays focused on the CLI skills full-stack developers use to build and operate applications. Grow the platform by making one learning path work end to end, then prove portability with a small second course. The [roadmap](../plan.md) defines the release sequence, prerequisites, and acceptance checks for the reader, playground, quizzes, snippets, review, and course packages.
