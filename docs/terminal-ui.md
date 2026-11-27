# Terminal UI

Bubble Tea handles key presses, resize events, and Bash handoff. Lip Gloss draws the panels and score bars. When Bash starts, the UI releases the terminal and restores it on exit.

| Screen | Controls |
| --- | --- |
| Lessons | Up and down select a lesson. Left and right jump workshops. Enter opens Read, d opens Do, and / searches by topic, tool, or command. Escape clears the current filter. |
| Read | Up and down scroll the full chapter. Left and right jump sections, and Page Up or Page Down moves a page. Press i for the contents, e for a worked example, or b to bookmark. |
| Worked example | Left and right move between frames. Enter advances and opens Do from the final frame. Press a to replay motion or x to toggle the permission preview when available. Escape returns to your reading position. |
| Bookmarks | Up and down select a saved section. Enter opens it. |
| Do and results | Enter opens Bash. Press p for available follow-up labs, e to edit a starter file, c to check, w to toggle automatic file checks, or h for hints. Escape returns from results to the task. |
| Lesson views | Tab switches Read and Do, and [ and ] open the previous or next lesson. Your reading position is remembered for each visit within the session. |
| All screens | ? opens help, B opens bookmarks, m returns to lessons, and q exits. Escape closes help, search, and the file picker. |

The catalog opens with Foundations and Core tools in order. A circle marks a new lesson, a half-filled circle marks an attempted lab, and a checkmark marks a completed lab. The wide layout shows the lessons beside the selected lesson's preview. Search matches lesson titles, source reference IDs, workshop names, tool names, and demo commands. Canceling search returns to the original view. Do puts the rubric and workspace side by side on a wide terminal. The results show the latest score, best score, and the next action.

The execute-bit preview highlights the permission that changes. A moving marker traces redirected output and a highlight steps through pipelines, Git commits, HTTP requests, and other example flows. The checker shows a spinner while working and fills the score bar when the result changes. The score number always shows the measured result. Run `CLI_SILLY_NO_MOTION=1 cli-silly` for a static display with the same information. Color-free terminals also use the static display.

Glamour renders headings, lists, fenced commands, and tables from embedded Markdown. The reader scrolls through the shared chapter and terminal drills as one document. On a wide terminal the contents stay visible beside it. The contents can also take focus for quick jumps. The CLI saves a lesson ID and section anchor as you move through the document. A removed heading falls back to the overview. The palette adapts to light and dark terminals and honors `NO_COLOR`.

The Bash prompt shows a path relative to `lab/`, the workspace root. Return there with `cd "$CLI_SILLY_LAB_ROOT"`. In the workspace, `cli-silly read` prints the current chapter, drills, and worked example. Run `cli-silly check` for feedback while practicing. Returning to the interface refreshes file-only checks against the same attempt. Git, runtime, and script labs use explicit checks.

The editor uses `VISUAL` or `EDITOR`. It opens the selected starter file from its managed workspace. Automatic checks run for simple workspace file criteria and save a score when the result changes. They stop when leaving the score view or pausing with W. Use a manual check for Git, script, Docker, and other runtime exercises.

**z** opens the selected lesson’s skill check. Choose with arrows or number keys and press Enter to submit. **V** opens the review queue. Early practice remains available through the optional challenge workshop on the map.

Press **n** to open reference examples. Search with **/**, preview with Enter, and press **p** to print and exit. Pass `name=value` to `cli-silly snippet` for custom parameter values.
