package main

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"syscall"
	"unsafe"

	"cli-silly/course"
	"cli-silly/internal/lab"
	"cli-silly/internal/tui"
)

var version = "dev"

const help = `cli-silly · Learn developer tools in your terminal.

Usage: cli-silly [command]

  list                       List available lessons
  learn [lesson-id]          Walk through a lesson and open Bash
  continue                   Open the first unfinished lesson
  resume                     Open the last-read section
  read [lesson-id]           Print the offline lesson
  open [lesson-id]           Open or resume the lab in Bash
  path [lesson-id]           Create or resume the lab and print its path
  task [lesson-id]           Show the task and rubric
  hint [number] [lesson-id]   Show a hint (1 through 3)
  check [--json] [lesson-id]  Grade the workspace (0 pass, 1 retry, 2 error)
  progress [lesson-id]       Show scores for the course or one lesson
  reset [lesson-id]          Archive lab files and start fresh
  doctor                     Check workshop dependencies
  --version                  Print the CLI version
`

type app struct {
	store          lab.Store
	lesson         course.Lesson
	out            io.Writer
	in             *bufio.Reader
	color          bool
	hintNumber     int
	singleProgress bool
	workspace      string
}

func main() {
	code, err := run(os.Args[1:])
	if err != nil {
		fmt.Fprintln(os.Stderr, "cli-silly:", err)
	}
	os.Exit(code)
}

func run(args []string) (int, error) {
	command := "menu"
	if len(args) != 0 {
		command, args = args[0], args[1:]
	}
	if command == "--version" {
		fmt.Println("cli-silly", version)
		return 0, nil
	}
	if command == "--help" || command == "-h" || command == "help" {
		fmt.Print(help)
		return 0, nil
	}

	encoded := false
	if command == "check" && len(args) > 0 && args[0] == "--json" {
		encoded, args = true, args[1:]
	}
	id := course.First
	explicit := false
	hint := 1
	if command == "hint" && len(args) > 0 && len(args[0]) == 1 && strings.Contains("123", args[0]) {
		hint, args = int(args[0][0]-'0'), args[1:]
	}
	switch command {
	case "learn", "read", "open", "path", "check", "task", "hint", "progress", "reset":
		if len(args) == 1 {
			id, explicit, args = args[0], true, nil
		}
	}
	if len(args) != 0 {
		return 2, errors.New("unexpected argument (see cli-silly --help)")
	}

	var cwd string
	if !explicit {
		switch command {
		case "check", "task", "hint", "progress", "reset":
			var err error
			cwd, id, err = lab.FindWorkspace()
			if err != nil {
				if command == "check" || !errors.Is(err, lab.ErrOutsideWorkspace) {
					return 2, err
				}
				id = course.First
			}
		}
	}

	lesson, err := course.Load(id)
	if err != nil {
		return 2, err
	}
	home, err := os.UserHomeDir()
	if err != nil {
		return 2, err
	}
	a := app{
		store: lab.Store{
			DataDir:  userDir("XDG_DATA_HOME", filepath.Join(home, ".local", "share")),
			StateDir: userDir("XDG_STATE_HOME", filepath.Join(home, ".local", "state")),
		},
		lesson: lesson, hintNumber: hint, singleProgress: explicit || cwd != "", workspace: cwd,
		out:   os.Stdout,
		in:    bufio.NewReader(os.Stdin),
		color: terminal(os.Stdout) && os.Getenv("NO_COLOR") == "" && os.Getenv("TERM") != "dumb",
	}
	if command == "resume" {
		reading, err := a.store.Reading()
		if err != nil {
			return 2, err
		}
		if saved, err := course.Load(reading.Last.LessonID); err == nil {
			a.lesson = saved
		}
	}
	if command == "continue" {
		entries, err := course.Entries()
		if err != nil {
			return 2, err
		}
		for _, entry := range entries {
			if entry.Lesson.ReviewOf != "" {
				continue
			}
			progress, err := a.store.Progress(entry.Lesson)
			if err != nil {
				return 2, err
			}
			if progress.Best < 100 {
				a.lesson = entry.Lesson
				break
			}
		}
	}
	if command == "check" {
		result, err := a.store.Check(lesson, cwd)
		if err != nil {
			return 2, err
		}
		if encoded {
			if err := json.NewEncoder(a.out).Encode(result); err != nil {
				return 2, err
			}
		} else {
			a.result(result)
		}
		if result.Complete {
			return 0, nil
		}
		return 1, nil
	}
	if err := a.dispatch(command); err != nil {
		return 2, err
	}
	return 0, nil
}

func (a app) dispatch(command string) error {
	interactive := terminal(os.Stdin) && terminal(os.Stdout) && os.Getenv("TERM") != "dumb"
	if interactive && (command == "menu" || command == "learn" || command == "continue" || command == "resume") {
		entries, err := course.Entries()
		if err != nil {
			return err
		}
		return tui.Run(entries, a.lesson.ID, command, a.color, a.store, func(lesson course.Lesson) (*exec.Cmd, error) {
			a.lesson = lesson
			return a.shellCommand()
		})
	}
	switch command {
	case "menu", "list":
		entries, err := course.Entries()
		if err != nil {
			return err
		}
		for _, entry := range entries {
			fmt.Fprintln(a.out, entry.Lesson.ID+"  "+entry.Lesson.Title)
		}
	case "learn", "continue", "resume":
		a.printLesson()
	case "read":
		fmt.Fprintf(a.out, "# %s\n\n%s\n", a.lesson.Title, a.lesson.Body)
		a.task()
	case "open":
		return a.shell()
	case "path":
		path, err := a.store.Open(a.lesson)
		if err != nil {
			return err
		}
		fmt.Fprintln(a.out, path)
	case "task":
		a.task()
	case "hint":
		fmt.Fprintln(a.out, a.lesson.Hints[a.hintNumber-1])
		fmt.Fprintln(a.out, "Use cli-silly hint 2 or cli-silly hint 3 for more detail.")
	case "progress":
		if a.singleProgress {
			return a.progress()
		}
		entries, err := course.Entries()
		if err != nil {
			return err
		}
		complete, required := 0, 0
		for _, entry := range entries {
			a.lesson = entry.Lesson
			if err := a.progress(); err != nil {
				return err
			}
			p, err := a.store.Progress(a.lesson)
			if err != nil {
				return err
			}
			if a.lesson.ReviewOf == "" {
				required++
				if p.Best == 100 {
					complete++
				}
			}
		}
		fmt.Fprintf(a.out, "\nCourse: %d/%d labs complete\n", complete, required)
		return nil
	case "reset":
		if !terminal(os.Stdin) {
			return errors.New("reset requires an interactive terminal")
		}
		fmt.Fprintf(a.out, "Archive %s and recreate its starter files?\n", a.store.Path(a.lesson))
		if a.lesson.Setup == "compose" {
			fmt.Fprintln(a.out, "Run docker compose down in this workspace first if the old project is still running.")
		}
		if a.ask("Type reset to confirm. ") == "reset" {
			archive, err := a.store.ResetFrom(a.lesson, a.workspace)
			if err != nil {
				return err
			}
			fmt.Fprintln(a.out, "Previous work saved in", archive)
		}
	case "doctor":
		entries, err := course.Entries()
		if err != nil {
			return err
		}
		for _, entry := range entries {
			status := "ready"
			if err := lab.Ready(entry.Lesson); err != nil {
				status = err.Error()
			}
			fmt.Fprintf(a.out, "%s  %s\n", entry.Lesson.ID, status)
		}
	default:
		return fmt.Errorf("unknown command %q (see cli-silly --help)", command)
	}
	return nil
}

func (a app) printLesson() {
	fmt.Fprintf(a.out, "# %s\n\n%s\n", a.lesson.Title, a.lesson.Body)
	for i, step := range a.lesson.Demo {
		fmt.Fprintf(a.out, "\nWorked example %d/%d\n\n$ %s\n\n%s\n\n%s\n",
			i+1, len(a.lesson.Demo), step.Command, step.Visual, step.Explanation)
	}
	a.task()
}

func (a app) shell() error {
	if !terminal(os.Stdin) || !terminal(os.Stdout) {
		return errors.New("opening Bash requires a terminal; use cli-silly path to find the workspace")
	}
	cmd, err := a.shellCommand()
	if err != nil {
		return err
	}
	cmd.Stdin, cmd.Stdout, cmd.Stderr = os.Stdin, os.Stdout, os.Stderr
	if err := cmd.Run(); err != nil {
		var exit *exec.ExitError
		if !errors.As(err, &exit) {
			return err
		}
		fmt.Fprintf(a.out, "\nBash exited with status %d.\n", exit.ExitCode())
	}
	return a.progress()
}

func (a app) shellCommand() (*exec.Cmd, error) {
	if err := lab.Ready(a.lesson); err != nil {
		return nil, err
	}
	bash, err := exec.LookPath("bash")
	if err != nil {
		return nil, err
	}
	executable, err := os.Executable()
	if err != nil {
		return nil, err
	}
	path, err := a.store.Open(a.lesson)
	if err != nil {
		return nil, err
	}
	const intro = `printf '\nWorkspace %s\nRun cli-silly task, cli-silly hint, or cli-silly check. Type exit to return.\n\n' "$PWD"
exec "$0" --noprofile --norc -i`
	cmd := exec.Command(bash, "--noprofile", "--norc", "-c", intro, bash)
	cmd.Dir = path
	cmd.Env = append(cmd.Environ(), "PATH="+filepath.Dir(executable)+":"+os.Getenv("PATH"), "PROMPT_COMMAND=PS1='cli-silly> '; unset PROMPT_COMMAND")
	if a.lesson.Setup == "compose" {
		project, err := lab.ComposeProject(path)
		if err != nil {
			return nil, err
		}
		cmd.Env = append(cmd.Env, "COMPOSE_PROJECT_NAME="+project, "COMPOSE_FILE="+filepath.Join(path, "compose.yaml"))
	}
	return cmd, nil
}

func (a app) task() {
	fmt.Fprintf(a.out, "\n%s\n%s\n\n", a.style("1;95", "YOUR TASK"), a.lesson.Task)
	for _, check := range a.lesson.Checks {
		fmt.Fprintf(a.out, "  %2d points  %s\n", check.Points, check.Label)
	}
}

func (a app) result(result lab.Result) {
	for _, item := range result.Items {
		label := a.style("92", "PASS")
		if item.Earned != item.Possible {
			label = a.style("93", "TRY ")
		}
		fmt.Fprintf(a.out, "%s  %d/%d  %s\n", label, item.Earned, item.Possible, item.Label)
		if item.Feedback != "" {
			fmt.Fprintln(a.out, "     "+item.Feedback)
		}
	}
	fmt.Fprintf(a.out, "\n%s\n", a.score(result.Score))
	if result.Complete {
		fmt.Fprintln(a.out, a.style("1;92", "✦ LAB COMPLETE · "+a.lesson.Title+""))
	}
}

func (a app) progress() error {
	p, err := a.store.Progress(a.lesson)
	if err != nil {
		return err
	}
	if p.Checks == 0 {
		fmt.Fprintf(a.out, "%s  ·  Ready to start\n", a.lesson.ID)
		return nil
	}
	fmt.Fprintf(a.out, "%s  ·  Latest %s · Best %d/100\n", a.lesson.ID, a.score(p.Latest.Score), p.Best)
	return nil
}

func (a app) score(points int) string {
	return fmt.Sprintf("%d/100  %s", points, a.style("96", strings.Repeat("█", points/10)+strings.Repeat("░", 10-points/10)))
}

func (a app) style(code, text string) string {
	if !a.color {
		return text
	}
	return "\x1b[" + code + "m" + text + "\x1b[0m"
}

func (a app) ask(prompt string) string {
	fmt.Fprint(a.out, prompt)
	line, err := a.in.ReadString('\n')
	if err != nil && line == "" {
		return "EOF"
	}
	return strings.TrimSpace(line)
}

func terminal(file *os.File) bool {
	var settings syscall.Termios
	_, _, errno := syscall.Syscall(syscall.SYS_IOCTL, file.Fd(), syscall.TCGETS, uintptr(unsafe.Pointer(&settings)))
	return errno == 0
}

func userDir(variable, fallback string) string {
	dir := os.Getenv(variable)
	if !filepath.IsAbs(dir) {
		dir = fallback
	}
	return filepath.Join(dir, "cli-silly")
}
