package lab

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"

	"cli-silly/course"
)

func Doctor(entries []course.Entry) string {
	lines := []string{fmt.Sprintf("cli-silly readiness · %s/%s", runtime.GOOS, runtime.GOARCH), "Supported application targets: native Linux amd64 and arm64.", "Reading works offline even when a practice tool is unavailable.", ""}
	cache := make(map[string]string)
	ready := 0
	for _, entry := range entries {
		key := strings.Join(entry.Lesson.Tools, "\x00")
		status, found := cache[key]
		if !found {
			status = "ready"
			if err := Ready(entry.Lesson); err != nil {
				status = err.Error()
			}
			cache[key] = status
		}
		if status == "ready" {
			ready++
		}
		lines = append(lines, entry.Lesson.ID+"  "+status)
	}
	lines = append(lines, "", fmt.Sprintf("%d/%d exercises ready", ready, len(entries)))
	if _, err := exec.LookPath("cli-silly"); err != nil {
		lines = append(lines, `PATH: add your install directory to PATH. For the default location: export PATH="$HOME/.local/bin:$PATH"`)
	}
	missing := []string{}
	for _, tool := range []string{"git", "jq", "python3", "curl", "tar", "ssh", "ps", "node"} {
		if _, err := exec.LookPath(tool); err != nil {
			missing = append(missing, tool)
		}
	}
	if len(missing) > 0 {
		lines = append(lines, "Missing optional commands: "+strings.Join(missing, ", "))
		release, _ := os.ReadFile("/etc/os-release")
		switch {
		case strings.Contains(string(release), "fedora") || strings.Contains(string(release), "rhel"):
			lines = append(lines, "Fedora package names: git jq python3 curl tar openssh-clients procps-ng nodejs. Use dnf info PACKAGE to inspect before installing.")
		case strings.Contains(string(release), "debian") || strings.Contains(string(release), "ubuntu"):
			lines = append(lines, "Debian/Ubuntu package names: git jq python3 curl tar openssh-client procps nodejs. Use apt show PACKAGE to inspect before installing.")
		default:
			lines = append(lines, "Use your Linux distribution's package manager to install the missing commands.")
		}
	}
	lines = append(lines, "Docker labs need Engine, Compose, daemon access, and the python:3.12-alpine image.", "Before going offline: docker pull python:3.12-alpine. Check Docker access with docker info.", "Upgrade using the installation guide; upgrades retain your workspaces and progress. No packages are installed by doctor.")
	return strings.Join(lines, "\n")
}
