"""Exercise the full-screen UI, lab handoff, and saved reader state in a Linux pseudo-terminal."""

import errno
import fcntl
import os
import pty
import select
import signal
import struct
import subprocess
import sys
import tempfile
import termios
import time
from contextlib import contextmanager
from pathlib import Path


@contextmanager
def terminal_session(binary, env, *args):
    child, terminal = pty.fork()
    if child == 0:
        fcntl.ioctl(0, termios.TIOCSWINSZ, struct.pack("HHHH", 27, 96, 0, 0))
        os.execve(binary, [binary, *args], env)

    buffer = b""

    def expect(text):
        nonlocal buffer
        target = text.encode()
        deadline = time.monotonic() + 10
        while target not in buffer and time.monotonic() < deadline:
            if select.select([terminal], [], [], 0.1)[0]:
                try:
                    buffer += os.read(terminal, 8192)
                except OSError as error:
                    if error.errno != errno.EIO:
                        raise
                    break
        if target not in buffer:
            raise AssertionError(f"Expected {text!r}, got {buffer[-1800:]!r}")
        buffer = buffer.split(target, 1)[1]

    def send(text):
        os.write(terminal, text.encode())

    try:
        yield expect, send
        send("q")
        _, code = os.waitpid(child, 0)
        child = 0
        assert code == 0
    finally:
        os.close(terminal)
        if child:
            os.kill(child, signal.SIGKILL)
            os.waitpid(child, 0)


def verify(binary, state):
    env = {
        **os.environ,
        "XDG_DATA_HOME": f"{state}/data",
        "XDG_STATE_HOME": f"{state}/state",
        "TERM": "xterm-256color",
    }
    lesson_count = len(subprocess.check_output([binary, "list"], env=env, text=True).splitlines())
    with terminal_session(binary, env) as (expect, send):
        expect("COURSE PROGRESS")
        expect("$ ls -l bin/show-version")
        send("/")
        expect("Type to filter")
        send("log-report")
        expect("Count server errors by route")
        send("\r")
        expect("WORKED EXAMPLE")
        send("l")
        expect("A pipe connects")
        send("?")
        expect("cli-silly read streams/log-report")
        send("?")
        send("\r")
        expect("WORKED EXAMPLE")
        send("t")
        expect("YOUR TASK")
        send("h")
        expect("Hint 1/3")
        send("\r")
        expect("cli-silly> ")
        send("cli-silly task\n")
        expect("Use access.log")
        send("cli-silly check\n")
        expect("0/100")
        send("awk '$3 >= 500 && $3 < 600 {print $2}' access.log > errors.txt; cli-silly check\n")
        expect("45/100")
        send("sort errors.txt | uniq -c | awk '{print $1, $2}' > counts.txt; cli-silly check\n")
        expect("100/100")
        send("exit\n")
        expect("Best 100/100")
        send("m")
        expect("COURSE PROGRESS")
        expect(f"1/{lesson_count} complete")
    scores = subprocess.check_output([binary, "progress"], env=env, text=True)
    assert "streams/log-report" in scores and f"Course: 1/{lesson_count} labs complete" in scores

    env["NO_COLOR"] = "1"
    with terminal_session(binary, env, "learn", "shell/make-it-runnable") as (expect, send):
        expect("WORKED EXAMPLE")
        send("l")
        expect("SECTION 1/3")
        send("\x1b[C")
        expect("SECTION 2/3")
        expect("Owner permission")
        expect("chmod u+x bin/show-version")
        send("b")
        expect("Bookmarked")

    with terminal_session(binary, env, "resume") as (expect, send):
        expect("SECTION 2/3")
        expect("Bookmarked")
        send("\x1b[C")
        expect("SECTION 3/3")
        send("B")
        expect("BOOKMARKS")
        expect("Read the permission bits")
        send("\r")
        expect("SECTION 2/3")


if __name__ == "__main__":
    cli = str(Path(sys.argv[1] if len(sys.argv) > 1 else "bin/cli-silly").resolve())
    with tempfile.TemporaryDirectory(prefix="cli-silly-terminal-") as directory:
        verify(cli, directory)
    print("Verified course navigation, reading, command help, Bash handoff, saved lab scoring, and reading resume with bookmarks.")
