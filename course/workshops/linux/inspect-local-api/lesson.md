Processes inherit exported environment variables from their parent. A variable can also be set for one command with `NAME=value command`. Changing a variable in the shell after a process starts cannot update that process.

This API binds to a free loopback port and writes it to server.port. A shell trailing `&` starts a background job, and `$!` gives its PID. Use `ps -p` to inspect the process and `curl -i` to inspect its HTTP status. Keep the workspace shell open while grading and stop the job with `kill` when you finish.
