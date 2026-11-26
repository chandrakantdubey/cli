# Terminal UI

Bubble Tea handles key presses, resize events, and Bash handoff. Lip Gloss draws the panels and score bars. When Bash starts, the UI releases the terminal and restores it on exit.

| Screen | Controls |
| --- | --- |
| Course map | Up and down select a lesson. Enter opens its demo. / searches the catalog, L reads the lesson, T opens its task, S shows its score, and lowercase g resumes reading. |
| Worked example | Left and right move between frames. Space advances. X toggles the owner execute bit when the frame contains permissions. L opens the reading and T opens the task. |
| Reading | Left and right change sections. Up and down scroll, and Page Up or Page Down moves a page. Press lowercase b to toggle the section bookmark. Enter opens the demo, and T opens the task. |
| Bookmarks | Up and down select a saved section. Enter opens it. |
| Task and results | Enter opens Bash. C checks the workspace, H reveals a hint, and R replays the demo. |
| All screens | ? lists terminal commands for the selected lesson. Uppercase B opens bookmarks, M opens the map, and Q exits. |

The wide map shows the learning paths alongside a preview of the selected lesson. The preview includes its worked command and, where space allows, the associated diagram. Search matches lesson titles, workshop names, tool names, and demo commands. The map groups lessons into workshops and shows the best score for each attempted lesson. Score feedback shows the latest result, individual criteria, and a next step. The best score records completion across retries. Shell commands resolve the lesson ID from the workspace marker.

The execute-bit diagram previews changes within the worked example. Learners perform the real operation in the lab. Other diagrams trace pipelines, Git commits, the HTTP request path, and Docker services.

Glamour renders headings, lists, fenced commands, and tables from the embedded Markdown. Top-level Markdown headings define reading sections. The CLI saves a lesson ID and section anchor, so a different terminal width can restore the same section. A removed heading falls back to the overview. Use text labels alongside color. Allow narrow terminals to scroll and honor NO_COLOR. Run `python3 scripts/test-terminal.py` after building the binary to test navigation, Bash handoff, and returning to a saved score in a different lab.
