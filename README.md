# Prompt

> I want you to generate **complete CLI/Linux course** from beginner to advanced. The project should be modular, interactive, and production-ready.
>
> **Project Requirements:**
>
> * Use **React functional components and hooks**.
> * Use **React Router** for module navigation.
> * Include a **sidebar** with all modules listed for navigation.
> * Each module should be a **separate page/component**, with:
>
>   * Key points summarized
>   * Examples & CLI code snippets
>   * Interactive CLI simulation for commands relevant to the module
>   * Quizzes (3–5 questions per module)
> * Include reusable components: `Sidebar.jsx`, `ModuleCard.jsx`, `CLIEmulator.jsx`, `Quiz.jsx`.
> * Include **state management** for CLI simulation and quiz scores.
> * Include comments explaining the code and functionality.
> * Make project **scalable**, so new modules can be added easily.
> * Ensure **responsive design** for desktop and mobile.
>
> **Modules & Core Content:**
>
> 1. **CLI & Linux Intro** – CLI vs GUI, Linux essentials (kernel, distro, filesystem, root), Example: `echo "Hello CLI!"`.
> 2. **Mental Model** – Flow: You → Terminal → Shell → Kernel → Output → You.
> 3. **CLI/System Core Terms** – CLI, Terminal, Shell, Prompt, Command, Argument, Option, Kernel, with examples and tips.
> 4. **Filesystem & Navigation** – Commands: `pwd`, `cd`, `ls`, `tree`. Paths: absolute, relative, home, root, current, parent.
> 5. **Users, Groups & Permissions** – `whoami`, `id`, `groups`, `su`, `sudo`, `passwd`, `chmod`, `chown`.
> 6. **File Operations** – Create, copy, move, delete, view/edit files: `touch`, `mkdir`, `cp`, `mv`, `rm`, `cat`, `less`, `nano`, `vim`.
> 7. **Redirection, Pipes & Wildcards** – `>`, `>>`, `<`, `2>`, `|`, `*`, `?`, `[]`.
> 8. **Text Processing & Search** – `grep`, `sort`, `uniq`, `cut`, `awk`, `sed`, `find`, `locate`.
> 9. **Processes, Jobs & Signals** – `ps`, `top`, `htop`, `&`, `bg`, `fg`, Ctrl+Z, `kill`, `pkill`, `systemctl`.
> 10. **Variables, Aliases & History** – `export`, `echo $VAR`, `alias`, `.bashrc`, history (`history`, `!!`, Ctrl+R), custom functions.
> 11. **Package Management** – Ubuntu `apt`, Fedora `dnf`, Arch `pacman`, Alpine `apk`, macOS `brew`. Install, remove, update, search.
> 12. **Networking Basics** – `ip a`, `ip link`, `ping`, `traceroute`, `mtr`, `ss`, `netstat`, `dig`, `nslookup`, `curl`, `wget`. Interfaces, routing, connectivity, DNS, ports, firewall, troubleshooting.
> 13. **SSH & Remote Access** – `ssh`, `scp`, `rsync`, key generation (`ssh-keygen`), ssh-copy-id, config file shortcuts (`~/.ssh/config`), tunneling (`-L`, `-R`).
> 14. **Advanced Process Management** – Priorities (`nice`, `renice`), background/persistence (`nohup`, `tmux`, `screen`), signals (`SIGTERM`, `SIGKILL`, `SIGHUP`), monitoring (`htop`, `glances`).
> 15. **Shell Scripting Deep Dive** – Bash scripts: variables, functions, loops, conditionals, stdin/stdout/stderr, debug (`set -x`, `set -e`), sample interactive script with user input.
> 16. **System Monitoring & Logs** – Disk (`iostat`, `iotop`), memory (`vmstat`, `free`), network (`iftop`, `nethogs`), logs (`/var/log`, `dmesg`, `journalctl`).
> 17. **Archive & Compression** – Tar, gzip, zip: create, extract, view. Commands: `tar -cvf`, `tar -xzvf`, `gzip`, `gunzip`, `zip`, `unzip`.
>
> **Output Structure:**
>
> ```
> /src
>   /components
>     Sidebar.jsx
>     ModuleCard.jsx
>     CLIEmulator.jsx
>     Quiz.jsx
>   /pages
>     Module1.jsx
>     Module2.jsx
>     ...
>   /utils
>     quizData.js
>   App.jsx
>   routes.jsx
> ```
>
> * Each file should contain **full JSX + Tailwind classes**.
> * CLI emulator should handle **basic relevant commands** per module.
> * Quizzes must be **functional and track scores**.
> * Include **comments** explaining components and hooks.
> * Ensure the code is **clean, modular, and maintainable**.