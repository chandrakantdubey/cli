package tui

import (
	"context"
	"os"
	"strings"
	"time"

	"cli-silly/internal/lab"

	tea "github.com/charmbracelet/bubbletea"
)

type editorExited struct{ err error }
type watchTick struct{ generation int }
type watched struct {
	generation int
	err        error
}

func (m *model) stopWatch() {
	if m.watching {
		m.watchCancel()
		m.watching = false
		m.watchGeneration++
	}
}

func nextWatch(generation int) tea.Cmd {
	return tea.Tick(400*time.Millisecond, func(time.Time) tea.Msg { return watchTick{generation} })
}

func (m model) toggleWatch() (tea.Model, tea.Cmd) {
	if m.watching {
		m.stopWatch()
		m.status = "Automatic checks paused."
		return m, nil
	}
	if !lab.Watchable(m.lesson) {
		m.status = "This lab uses a manual check. Press c when ready."
		return m, nil
	}
	path, err := m.store.Open(m.lesson)
	if err != nil {
		m.status = err.Error()
		return m, nil
	}
	attempt, err := os.Stat(path)
	if err != nil {
		m.status = err.Error()
		return m, nil
	}
	m.watchAttempt = attempt
	m.watchContext, m.watchCancel = context.WithCancel(context.Background())
	m.watching = true
	m.watchGeneration++
	m.screen, m.scroll, m.status = resultsScreen, 0, ""
	return m, nextWatch(m.watchGeneration)
}

func (m model) openFiles() (tea.Model, tea.Cmd) {
	files, err := lab.StarterFiles(m.lesson)
	if err != nil {
		m.status = err.Error()
		return m, nil
	}
	m.stopWatch()
	m.files, m.fileIndex = files, 0
	m.screen, m.scroll, m.status = filesScreen, 0, ""
	return m, nil
}

func (m model) filesView() string {
	lines := []string{m.paint(violet, "STARTER FILES", true), m.store.Path(m.lesson), ""}
	if len(m.files) == 0 {
		return strings.Join(append(lines, "This lab has no starter files to edit."), "\n")
	}
	for i, name := range m.files {
		label := "  " + name
		if i == m.fileIndex {
			label = m.paint(teal, "▸ "+name, true)
		}
		lines = append(lines, label)
	}
	return strings.Join(lines, "\n")
}

func (m model) practiceCommands() string {
	body := "\n\n" + m.paint(violet, "Inside Bash", true) + "\n  cli-silly read   Reopen the full reading as text\n  cli-silly task   See the goal again\n  cli-silly check  See your score now\n  exit             Return to this lesson\n\nWorkspace: " + m.store.Path(m.lesson)
	if lab.Watchable(m.lesson) {
		body += "\nResults refresh when you return from Bash."
	}
	switch {
	case m.lesson.Setup == "compose":
		body += "\nInspect or stop this project's containers inside the lab with docker compose ps or docker compose down."
	default:
		for _, check := range m.lesson.Checks {
			if check.Kind == "http-body" {
				body += "\nFor an API started as a background job in Bash, use jobs -l and stop that job before leaving."
				break
			}
		}
	}
	return body
}

func (m model) nextLesson() int {
	if len(m.entries) == 0 {
		return -1
	}
	for i := m.menu + 1; i < len(m.entries); i++ {
		entry := m.entries[i]
		if entry.PackID == m.entries[m.menu].PackID && entry.Lesson.ReviewOf == "" {
			return i
		}
	}
	return -1
}
