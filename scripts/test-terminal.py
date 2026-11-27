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
    editor = Path(state) / "test editor"
    editor.write_text("#!/bin/sh\nprintf '%s\\n' \"$1\" > editor-opened.txt\n")
    editor.chmod(0o700)
    env["VISUAL"] = str(editor)
    progress = subprocess.check_output([binary, "progress"], env=env, text=True)
    lesson_count = int(progress.splitlines()[-1].split("/")[1].split()[0])
    with terminal_session(binary, env) as (expect, send):
        expect("COURSE PROGRESS")
        expect("CLI & Linux Intro")
        expect("Enter Read")
        send("d")
        expect("YOUR TASK")
        expect("welcome-copy.txt")
        send("m")
        expect("COURSE PROGRESS")
        send("/")
        expect("Type to filter")
        send("log-report")
        expect("Text Manipulation")
        send("\r")
        expect("A pipe connects")
        send("?")
        expect("cli-silly read streams/log-report")
        send("?")
        expect("A pipe connects")
        send("e")
        expect("WORKED EXAMPLE")
        send("d")
        expect("YOUR TASK")
        send("e")
        expect("STARTER FILES")
        expect("access.log")
        send("\r")
        expect("YOUR TASK")
        workspace = Path(subprocess.check_output(
            [binary, "path", "streams/log-report"], env=env, text=True
        ).strip())
        assert (workspace / "editor-opened.txt").read_text().strip() == "./access.log"
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
        send("w")
        expect("Watching file checks")
        report = workspace / "errors.txt"
        correct = report.read_text()
        report.write_text("/wrong-route\n")
        expect("55/100")
        report.write_text(correct)
        expect("LAB COMPLETE")
        send("w")
        expect("Automatic checks paused.")
        send("m")
        expect("COURSE PROGRESS")
        expect(f"1/{lesson_count} complete")
        send("z")
        expect("SKILL CHECK")
        expect("status field")
        send("2")
        send("\r")
        expect("Correct.")
        send("\r")
        expect("Question 2/2")
        send("\r")
        expect("Correct.")
        send("\r")
        expect("2/2 correct")
        send("V")
        expect("REVIEW QUEUE")
        expect("Trace a worker incident")
        send("\r")
        expect("YOUR TASK")
        expect("failed-workers.txt")
        send("\r")
        expect("cli-silly> ")
        send("awk '$3 >= 500 && $3 < 600 {print $2}' events.log > failed-workers.txt\n")
        send("sort failed-workers.txt | uniq -c | awk '{print $1, $2}' > failures.txt\n")
        send("jq -r '.workers[] | select(.name == \"worker-alpha\") | .team' workers.json > owner.txt; cli-silly check\n")
        expect("LAB COMPLETE")
        send("exit\n")
        expect("Best 100/100")
        send("n")
        expect("REFERENCE LIBRARY")
        send("/")
        expect("Type to filter")
        send("curl")
        expect("Check the status of a local HTTP service")
        send("\r")
        send("\r")
        expect("curl -fsS")
    scores = subprocess.check_output([binary, "progress"], env=env, text=True)
    assert "streams/log-report" in scores and f"Course: 1/{lesson_count} labs complete" in scores

    env["NO_COLOR"] = "1"
    with terminal_session(binary, env, "learn", "shell/make-it-runnable") as (expect, send):
        expect("WORKED EXAMPLE")
        send("r")
        expect("SECTION 1/")
        send("\x1b[C")
        expect("SECTION 2/")
        expect("Owner permission")
        expect("chmod u+x bin/show-version")
        send("b")
        expect("Bookmarked")

    with terminal_session(binary, env, "resume") as (expect, send):
        expect("SECTION 2/")
        expect("Bookmarked")
        send("\x1b[C")
        expect("SECTION 3/")
        send("B")
        expect("BOOKMARKS")
        expect("Read the permission bits")
        send("\r")
        expect("SECTION 2/")

    with terminal_session(binary, env, "do", "core-tools/node-cli") as (expect, send):
        expect("YOUR TASK")
        expect("Finish line-count.mjs")
        send("r")
        expect("SECTION 1/")
        expect("Build a tiny Node.js command")
        send("d")
        expect("YOUR TASK")


if __name__ == "__main__":
    cli = str(Path(sys.argv[1] if len(sys.argv) > 1 else "bin/cli-silly").resolve())
    with tempfile.TemporaryDirectory(prefix="cli-silly-terminal-") as directory:
        verify(cli, directory)
    print("Verified course navigation, reading, command help, Bash handoff, saved lab scoring, and reading resume with bookmarks, quizzes, and review navigation.")
