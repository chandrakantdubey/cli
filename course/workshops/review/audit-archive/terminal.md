## Try it in the terminal

Print the intended input and a neighboring file with `ls -l reports`. Compare a file's host path with the member name printed by tar. If the listing shows a prefix or more than one file, create a fresh archive with the single quoted member.

Try reading the member to the terminal before redirecting it. Use `cmp` on your extracted report and original log after both exist to check that their bytes match.
