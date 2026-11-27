## Try it in the terminal

### Follow three different channels

```text
bin/sync.sh   ─── stdout ───▶  ordinary data
              ├── stderr ───▶  diagnostic
              └── status ───▶  "$?" in Bash
```

Run `bash bin/sync.sh` and immediately run `printf 'exit=%s\n' "$?"`. Compare this with running `ls` between those commands. The second experiment reports the status from `ls`, because Bash updates the previous-command status for every command it runs.

Capture the output and diagnostic in separate scratch files with `>` and `2>`. Their contents are text and neither file contains the status automatically. For the scored task, leave the diagnostic visible and save the status before using `cat` to verify your output.
