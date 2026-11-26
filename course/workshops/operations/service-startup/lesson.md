A systemd service starts a process with the user, environment, and working directory from its unit. Use a journal to distinguish failures before the process starts from failures reported by the application itself.

This exercise uses a journal and unit fixture. On a live system, `systemctl cat UNIT` shows the unit, and `journalctl -u UNIT -n 30` shows recent messages. Keep your edits inside the lab.
