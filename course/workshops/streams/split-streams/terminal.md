## Try it in the terminal

### Capture once and inspect separately

Run `bash sync-check.sh` and immediately record the status in a shell variable before inspecting anything else. Run the command a second time using scratch filenames to separate stdout and stderr. Compare the files with `wc -l` and `cat`. A terminal can show both streams mixed even though those lines travel over different file descriptors.

To practice recovering a status when shell options stop a script on failure, try `status=0; false || status=$?; printf '%s\n' "$status"`. The right side of `||` runs only after failure. Use the same idea when capturing the real report and save just the numeric status as the task requests.
