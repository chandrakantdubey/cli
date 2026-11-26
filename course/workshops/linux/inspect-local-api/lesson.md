Processes inherit exported environment variables from their parent. A variable can also be set for one command with `NAME=value command`. Changing a variable in the shell after a process starts cannot update that process.

This API binds to a free loopback port and writes it to server.port. A shell trailing `&` starts a background job, and `$!` gives its PID. Use `ps -p` to inspect the process and `curl -i` to inspect its HTTP status. Keep the workspace shell open while grading and stop the job with `kill` when you finish.

## Commands and operations

Run the API from the lab directory so Python can find `app.py` and write `server.port` where the checker expects it. `APP_MODE=ready python3 app.py` provides the environment variable for that one process. The trailing `&` gives you back the prompt while the process keeps listening.

- `> server.log 2>&1` sends both output streams to a single file. If the process exits early or `server.port` never appears, inspect the log first.
- Save `$!` as soon as you start the background process, or run `jobs -l` from the same shell to see the job and PID. Use `ps -p PID` to check that the process is still present.
- Build the URL using the port from `server.port`. `curl -i` shows the HTTP status and headers as well as the body. When saving just the body, use `curl -fsS` and redirect standard output. `-f` reports HTTP failures, `-s` hides progress, and `-S` keeps errors visible.

`Connection refused` means you did not reach a listener at that address. A `503` means you did reach the application, but it is in maintenance mode. Keep the process running for `cli-silly check`, then stop the job you started.
