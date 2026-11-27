Start a disposable Bash worker, inspect exactly the child you started, and shut it down gracefully. The worker tells you when it is ready and writes proof when it handles TERM.

## Commands and operations

Run `bash worker.sh &` and immediately capture its PID with `worker=$!`. The ampersand starts the job in the background. `$!` identifies the most recently started background command in this shell. Record it before starting another one.

`ps -p "$worker" -o args=` selects that PID and prints its command line without a header. Wait until `cat ready.txt` shows `ready` before sending a signal so the trap has time to install. If you see no process, start the lab worker again and capture the fresh PID.

`kill -TERM "$worker"` requests termination of that exact worker. `wait "$worker"` asks Bash to wait for the child and collect its status. The worker's trap writes `handled TERM` to stopped.txt. The fixture times out after about 20 seconds if left alone. Use a fresh lab worker for each try.

## Inspect processes

```bash
ps aux
ps -ef
ps -eo pid,ppid,user,stat,%cpu,%mem,comm --sort=-%cpu | head
```

Important fields include PID, parent PID, user, state, CPU, and memory.

## Interactive monitoring

```bash
top
```

If installed, `htop` provides a more interactive interface.

## Find a process

```bash
pgrep -af nginx
```

Avoid killing a process merely because its name looks similar to the target.

## Signals

Common signals:

```text
SIGTERM  request graceful termination
SIGKILL  force termination; cannot be caught or handled
SIGHUP   commonly used for reload/reopen behavior
SIGINT   interrupt, often from Ctrl-C
```

Prefer graceful termination:

```bash
kill <PID>
```

Use `kill -9` only when a normal termination path is not working and you understand the consequences.

## Foreground and background jobs

```bash
sleep 100 &
jobs
fg
bg
```

`Ctrl-Z` suspends the foreground job; `bg` resumes it in the background.

## Process trees

```bash
pstree -p
```

Parent/child relationships matter when diagnosing supervisors, workers, and orphaned processes.

## Open files and sockets

```bash
lsof -p <PID>
```

This can reveal configuration files, logs, sockets, and deleted-but-open files.

## Resource limits

Processes can be constrained by limits such as open files and process counts:

```bash
ulimit -a
```

Systemd and container runtimes can impose additional controls.

## Zombies

A zombie is an exited child whose parent has not yet collected its exit status. It is not a "running process consuming CPU."
