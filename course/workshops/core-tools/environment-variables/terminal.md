## Try it in the terminal

### Compare unset, temporary, and inherited values

Use a subshell for this sequence so you can experiment without affecting the next Bash command after it:

```bash
(
    unset APP_MODE
    bash bin/show-mode.sh
    APP_MODE=dry-run bash bin/show-mode.sh
    bash bin/show-mode.sh
)
```

You should see the default value, the temporary value, then the default again. Now try `(export APP_MODE=dry-run; bash bin/show-mode.sh; bash bin/show-mode.sh)`. This time both child commands inherit the exported value. Use your two scored modes to compare command-local and exported settings in the real task. Keep the report filenames distinct.
