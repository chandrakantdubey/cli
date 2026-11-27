# Terminal UI

Bubble Tea handles key presses, resize events, and Bash handoff. Lip Gloss draws the panels and score bars. When Bash starts, the UI releases the terminal and restores it on exit.

| Screen | Controls |
| --- | --- |
| Lessons | Up and down select a lesson. Enter or r opens Read, d opens Do, and / searches by topic, tool, or command. |
| Read | Up and down scroll the full chapter. Left and right jump sections, and Page Up or Page Down moves a page. Press i for the contents, e for a worked example, b to bookmark, or d to start Do. |
| Worked example | Left and right move between frames. Space advances. X toggles the owner execute bit when the frame contains permissions. Press r to return to the reading or d for the task. |
| Bookmarks | Up and down select a saved section. Enter opens it. |
| Do and results | Enter opens Bash after the task and rubric. Press e to edit a starter file, c to check the workspace, w to toggle file-based automatic checks, h for hints, or r to return to reading. After passing, press ] to read the next lesson. |
| All screens | ? lists terminal commands for the selected lesson. Uppercase B opens bookmarks, m returns to the lessons, and q exits. |

The catalog opens with Foundations and Core tools in order. The wide layout shows the lessons beside the selected lesson's preview. Read includes command explanations and worked examples. Do includes the workspace, checking, and results. Search matches lesson titles, source reference IDs, workshop names, tool names, and demo commands. Score feedback shows the latest result, individual criteria, and a next step. The best score records completion across retries. Shell commands resolve the lesson ID from the workspace marker.

The execute-bit diagram previews changes within the worked example. Learners perform the real operation in the lab. Other diagrams trace pipelines, Git commits, the HTTP request path, and Docker services.

Glamour renders headings, lists, fenced commands, and tables from embedded Markdown. The reader scrolls through the shared chapter and terminal drills as one document. On a wide terminal the contents stay visible beside it. The contents can also take focus for quick jumps. The CLI saves a lesson ID and section anchor as you move through the document. A removed heading falls back to the overview. The palette adapts to light and dark terminals and honors `NO_COLOR`.

The Bash prompt shows a path relative to `lab/`, the workspace root. Return there with `cd "$CLI_SILLY_LAB_ROOT"`. In the workspace, `cli-silly read` prints the current chapter, drills, and worked example. Run `cli-silly check` for feedback while practicing. Returning to the interface refreshes file-only checks against the same attempt. Git, runtime, and script labs use explicit checks.

The editor uses `VISUAL` or `EDITOR`. It opens the selected starter file from its managed workspace. Automatic checks run for simple workspace file criteria and save a score when the result changes. They stop when leaving the score view or pausing with W. Use a manual check for Git, script, Docker, and other runtime exercises.

**z** opens the selected lesson’s skill check. Choose with arrows or number keys and press Enter to submit. **V** opens the review queue. Early practice remains available through the optional challenge workshop on the map.

Press **n** to open reference examples. Search with **/**, preview with Enter, and press **p** to print and exit. Pass `name=value` to `cli-silly snippet` for custom parameter values.
