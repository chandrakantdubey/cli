package lab

import (
	"context"
	"crypto/rand"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"net"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"slices"
	"strconv"
	"strings"
	"time"

	"cli-silly/course"
)

var ErrOutsideWorkspace = errors.New("run from inside a lab or specify its lesson ID")

func FindWorkspace() (string, string, error) {
	dir, err := os.Getwd()
	if err != nil {
		return "", "", err
	}
	for {
		root, err := os.OpenRoot(dir)
		if err != nil {
			return "", "", err
		}
		data, err := readRegular(root, Marker)
		root.Close()
		if err == nil {
			var marker workspace
			if err := json.Unmarshal(data, &marker); err != nil || marker.LessonID == "" {
				return "", "", fmt.Errorf("invalid workspace marker in %s", dir)
			}
			return dir, marker.LessonID, nil
		}
		if !errors.Is(err, fs.ErrNotExist) {
			return "", "", err
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			return "", "", ErrOutsideWorkspace
		}
		dir = parent
	}
}

func Ready(lesson course.Lesson) error {
	for _, tool := range append([]string{"bash"}, lesson.Tools...) {
		if _, err := exec.LookPath(tool); err != nil {
			return fmt.Errorf("install %s to practice this lesson", tool)
		}
	}
	if slices.Contains(lesson.Tools, "docker") {
		if _, err := output("", "docker", "compose", "version"); err != nil {
			return errors.New("Docker Compose is required for this lab")
		}
		if _, err := output("", "docker", "info", "--format", "{{.ServerVersion}}"); err != nil {
			return errors.New("Docker is installed but the daemon is unavailable to this user")
		}
	}
	return nil
}

func initialize(dir, setup string, marker *workspace) error {
	switch setup {
	case "":
		return nil
	case "compose":
		marker.Project = "silly-" + strings.ToLower(rand.Text())
		return os.WriteFile(filepath.Join(dir, ".env"), []byte("COMPOSE_PROJECT_NAME="+marker.Project+"\n"), 0600)
	case "git-staging", "git-merge", "git-remote", "git-conflict":
		steps := [][]string{
			{"init", "-q", "-b", "main"}, {"config", "user.name", "CLI Silly"},
			{"config", "user.email", "learner@example.invalid"},
			{"add", "."}, {"commit", "-qm", "Starter files"},
			{"tag", "cli-silly-start"},
		}
		if setup == "git-merge" || setup == "git-conflict" {
			steps = append(steps, []string{"switch", "-qc", "feature"})
		}
		for _, args := range steps {
			if _, err := git(dir, args...); err != nil {
				return err
			}
		}
		switch setup {
		case "git-staging":
			if err := os.WriteFile(filepath.Join(dir, "app.cfg"), []byte("mode=ready\n"), 0644); err != nil {
				return err
			}
			return os.WriteFile(filepath.Join(dir, "notes.md"), []byte("Draft: verify the next release.\n"), 0644)
		case "git-merge":
			if err := os.WriteFile(filepath.Join(dir, "feature.txt"), []byte("feature enabled\n"), 0644); err != nil {
				return err
			}
			for _, args := range [][]string{{"add", "feature.txt"}, {"commit", "-qm", "Enable feature"}, {"switch", "-q", "main"}} {
				if _, err := git(dir, args...); err != nil {
					return err
				}
			}
		case "git-conflict":
			if err := os.WriteFile(filepath.Join(dir, "settings.txt"), []byte("theme=dark\n"), 0644); err != nil {
				return err
			}
			for _, args := range [][]string{{"commit", "-qam", "Set dark theme"}, {"switch", "-q", "main"}} {
				if _, err := git(dir, args...); err != nil {
					return err
				}
			}
			if err := os.WriteFile(filepath.Join(dir, "settings.txt"), []byte("theme=light\n"), 0644); err != nil {
				return err
			}
			_, err := git(dir, "commit", "-qam", "Set light theme")
			return err
		case "git-remote":
			// The nested bare remote is local to the lab and ignored by the working repo.
			if _, err := git(dir, "init", "-q", "--bare", "-b", "main", "remote.git"); err != nil {
				return err
			}
			_, err := git(dir, "remote", "add", "origin", "remote.git")
			return err
		}
		return nil
	default:
		return fmt.Errorf("unknown setup %q", setup)
	}
}

func output(directory, name string, args ...string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	cmd := exec.CommandContext(ctx, name, args...)
	cmd.Dir = directory
	if name == "git" {
		for _, entry := range cmd.Environ() {
			if !strings.HasPrefix(entry, "GIT_") {
				cmd.Env = append(cmd.Env, entry)
			}
		}
		cmd.Env = append(cmd.Env, "GIT_CONFIG_NOSYSTEM=1", "GIT_CONFIG_GLOBAL=/dev/null", "GIT_TERMINAL_PROMPT=0")
	}
	var stdout boundedOutput
	cmd.Stdout = &stdout
	cmd.WaitDelay = time.Second
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("%s %s: %w", name, strings.Join(args, " "), err)
	}
	return stdout.String(), nil
}

func git(directory string, args ...string) (string, error) {
	return output(directory, "git", append([]string{"-c", "core.hooksPath=/dev/null", "-c", "core.fsmonitor=false", "-c", "commit.gpgsign=false", "-c", "init.templateDir="}, args...)...)
}

func gitCriterion(directory string, root *os.Root, check course.Criterion) string {
	info, err := root.Lstat(".git")
	if err != nil || !info.IsDir() {
		return ""
	}
	var args []string
	switch check.Kind {
	case "git-ancestor":
		if _, err := git(directory, "merge-base", "--is-ancestor", check.Path, "HEAD"); err == nil {
			return "merged"
		}
		return ""
	case "git-file":
		value, _ := git(directory, "show", "HEAD:"+check.Path)
		return value
	case "git-count":
		// Accept commits descended from the starter only.
		if _, err := git(directory, "merge-base", "--is-ancestor", "cli-silly-start", "HEAD"); err != nil {
			return ""
		}
		args = []string{"rev-list", "--count", "cli-silly-start..HEAD"}
	case "git-tree":
		if check.Path == "remote.git" {
			info, err := root.Lstat("remote.git")
			if err != nil || !info.IsDir() {
				return ""
			}
			local, err := git(directory, "rev-parse", "HEAD")
			remote, remoteErr := git(directory, "--git-dir=remote.git", "rev-parse", "refs/heads/main")
			if err == nil && remoteErr == nil && local == remote {
				return "matches"
			}
			return ""
		}
		value, err := git(directory, "diff", "--no-ext-diff", "--no-textconv", "--name-only", "cli-silly-start", "HEAD", "--")
		if err != nil || strings.TrimSpace(value) != check.Path {
			return ""
		}
		value, _ = git(directory, "show", "HEAD:"+check.Path)
		return value
	case "git-index":
		_, err := git(directory, "diff", "--no-ext-diff", "--no-textconv", "--cached", "--quiet")
		if err == nil {
			return "clean"
		}
		return ""
	case "git-branch":
		args = []string{"branch", "--show-current"}
	}
	value, _ := git(directory, args...)
	return strings.TrimSpace(value)
}

func ComposeProject(directory string) (string, error) {
	root, err := os.OpenRoot(directory)
	if err != nil {
		return "", err
	}
	defer root.Close()
	return readProject(root)
}

func readProject(root *os.Root) (string, error) {
	data, err := readRegular(root, Marker)
	if err != nil {
		return "", err
	}
	var marker workspace
	if err := json.Unmarshal(data, &marker); err != nil || !strings.HasPrefix(marker.Project, "silly-") {
		return "", errors.New("missing lab Compose project")
	}
	return marker.Project, nil
}

func dockerCriterion(root *os.Root, check course.Criterion) (string, error) {
	project, err := readProject(root)
	if err != nil {
		return "", err
	}
	value, err := output("", "docker", "ps", "-q",
		"--filter", "label=com.docker.compose.project="+project,
		"--filter", "label=com.docker.compose.service="+check.Path)
	if err != nil || strings.TrimSpace(value) == "" {
		return "", err
	}
	switch check.Kind {
	case "docker-running":
		return "running", nil
	case "docker-recreated":
		previous, err := readRegular(root, check.Member)
		if err != nil {
			return "", nil
		}
		ids := strings.Fields(value)
		if len(ids) != 1 {
			return "", nil
		}
		old := strings.TrimSpace(string(previous))
		if len(old) >= 12 && len(old) <= 64 && strings.Trim(old, "abcdef0123456789") == "" && !strings.HasPrefix(old, ids[0]) && !strings.HasPrefix(ids[0], old) {
			return "recreated", nil
		}
		return "", nil
	}
	value, err = output("", "docker", "inspect", "--format", "{{if .State.Health}}{{.State.Health.Status}}{{end}}", strings.Fields(value)[0])
	return strings.TrimSpace(value), err
}

func httpBody(root *os.Root, path string) string {
	data, err := readRegular(root, path)
	if err != nil {
		return ""
	}
	port := strings.TrimSpace(string(data))
	number, err := strconv.Atoi(port)
	if err != nil || number < 1024 || number > 65535 {
		return ""
	}
	transport := &http.Transport{}
	defer transport.CloseIdleConnections()
	client := http.Client{
		Timeout: 2 * time.Second, Transport: transport,
		CheckRedirect: func(*http.Request, []*http.Request) error { return http.ErrUseLastResponse },
	}
	res, err := client.Get("http://" + net.JoinHostPort("127.0.0.1", port) + "/health")
	if err != nil {
		return ""
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		return ""
	}
	const limit = 64 << 10
	body, err := io.ReadAll(io.LimitReader(res.Body, limit+1))
	if err != nil || len(body) > limit {
		return ""
	}
	return string(body)
}

func match(item Item, check course.Criterion, actual string) (Item, error) {
	if actual == check.Expected {
		item.Earned = item.Possible
	}
	return item, nil
}
