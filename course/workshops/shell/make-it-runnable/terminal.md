## Try it in the terminal

### Verify the script at three layers

First read the interpreter line with `head -n 1 bin/health-check`. Then check the file mode with `ls -l bin/health-check`. Finally use `test -x bin/health-check` and inspect the command status. A script may be valid text even if it cannot yet run directly.

After enabling the owner's execute bit, compare `bash bin/health-check` and `./bin/health-check` from the workspace root. Both can show the same result once the executable and interpreter are valid. The scored task needs the direct path. If the report is empty, run the command once without redirection so you can read its diagnostic before trying again.
