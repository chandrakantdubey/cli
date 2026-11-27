# Terminal UI

Bubble Tea handles key presses, resize events, and Bash handoff. Lip Gloss draws the panels and score bars. When Bash starts, the UI releases the terminal and restores it on exit.

| Screen | Controls |
| --- | --- |
| Lessons | Up and down select a lesson. Enter or r opens Read, d opens Do, and / searches by topic, tool, or command. |
| Read | Left and right change sections. Up and down scroll, and Page Up or Page Down moves a page. Press e for the worked example, b to toggle a bookmark, and d to start Do. |
| Worked example | Left and right move between frames. Space advances. X toggles the owner execute bit when the frame contains permissions. Press r to return to the reading or d for the task. |
| Bookmarks | Up and down select a saved section. Enter opens it. |
| Do and results | Enter opens Bash after the task and rubric. Press e to edit a starter file, c to check the workspace, w to toggle file-based automatic checks, h for hints, or r to return to reading. From Bash, run `cli-silly check` and type `exit` to return. |
| All screens | ? lists terminal commands for the selected lesson. Uppercase B opens bookmarks, m returns to the lessons, and q exits. |

The catalog opens with Foundations and Core tools in order. The wide layout shows the lessons beside the selected lesson's preview. Read includes command explanations and worked examples. Do includes the workspace, checking, and results. Search matches lesson titles, source reference IDs, workshop names, tool names, and demo commands. Score feedback shows the latest result, individual criteria, and a next step. The best score records completion across retries. Shell commands resolve the lesson ID from the workspace marker.

The execute-bit diagram previews changes within the worked example. Learners perform the real operation in the lab. Other diagrams trace pipelines, Git commits, the HTTP request path, and Docker services.

Glamour renders headings, lists, fenced commands, and tables from the embedded Markdown. Top-level Markdown headings define reading sections. The CLI saves a lesson ID and section anchor, so a different terminal width can restore the same section. A removed heading falls back to the overview. Use text labels alongside color. Allow narrow terminals to scroll and honor NO_COLOR. Run `python3 scripts/test-terminal.py` after building the binary to test navigation, Bash handoff, and returning to a saved score in a different lab.

The editor uses `VISUAL` or `EDITOR`. It opens the selected starter file from its managed workspace. Automatic checks run for simple workspace file criteria and save a score when the result changes. They stop when leaving the score view or pausing with W. Use a manual check for Git, script, Docker, and other runtime exercises.

**z** opens the selected lesson’s skill check. Choose with arrows or number keys and press Enter to submit. **V** opens the review queue. Early practice remains available through the optional challenge workshop on the map.

Press **n** to open reference examples. Search with **/**, preview with Enter, and press **p** to print and exit. Pass `name=value` to `cli-silly snippet` for custom parameter values.
