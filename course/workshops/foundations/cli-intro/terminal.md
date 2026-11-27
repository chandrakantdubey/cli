## Try it in the terminal

### See where a command will write

When the prompt says `cli-silly:lab/>`, you are beside the `docs` directory. Use `ls docs` to see inside it without moving. A command that writes `> some-file.txt` creates the file in your current directory. If your prompt says `lab/docs`, that same redirection creates `docs/some-file.txt`.

```text
lab/
├── docs/                   ← inputs
│   ├── checklist.txt
│   └── welcome.txt
├── welcome-copy.txt        ← your output
└── docs-list.txt           ← your output
```

### Try it in Bash

1. Run `pwd`, then `ls docs`. Run `pwd` again. Its output stays the same because `ls` inspects a path without moving your shell.
2. Run `cd docs` and `pwd`. Notice the prompt and the last directory name. Return with `cd ..` or `cd "$CLI_SILLY_LAB_ROOT"`.
3. Try `printf '%s\n' 'first line' > scratch.txt` at the root, then `cat scratch.txt`. Repeat with a different message. The second `>` replaces the original contents. The scratch file is outside `docs/` and has no effect on your score.

Start the scored task from `lab/`. Inspect each report with `cat` and run `cli-silly check` any time. If you see `Found lab/docs/welcome-copy.txt`, Bash used `docs` as the starting directory for the redirection. Return to `lab/` and move the copy out with `mv -- docs/welcome-copy.txt welcome-copy.txt`. Keep both starter files. If you removed one, `cli-silly reset` can archive the attempt and recreate the starter after you confirm.
