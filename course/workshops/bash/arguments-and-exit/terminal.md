## Try it in the terminal

### Test the boundary with different names

Run the starter with `bash check-config.sh "fixtures/demo config"` and inspect what it currently returns. The quotes pass one filename, including the space, rather than splitting it into two arguments. Inspect `"$#"` and `"$1"` in a scratch example if that distinction is unclear.

Design the success path to keep file content on stdout. Design the failure path to keep its diagnostic on stderr. Capture the two statuses immediately and test again with a different missing filename. The grader runs your script with a fresh directory and names, so a script that simply prints a stored answer will fail the behavior check.
