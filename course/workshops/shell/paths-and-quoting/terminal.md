## Try it in the terminal

### Make Bash show the arguments

Try `printf '[%s]\n' Draft Files/server config.txt`, then repeat with `"Draft Files/server config.txt"` quoted. The unquoted version displays three arguments and the quoted version displays one. You can inspect this behavior without moving or deleting any files.

Before using `mv`, inspect both directories with `ls -la "Draft Files" config`. The hidden `.keep` was placed there deliberately. After the move use `ls -la "Draft Files"` to check that it stayed put. Compare the content of the destination with what you saw in the draft before moving it.
