package tui

import (
	"strings"
	"testing"

	"cli-silly/course"
	"cli-silly/internal/lab"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

func TestWelcomeActionsAndWorkspaceReturn(t *testing.T) {
	entries, err := course.Entries()
	if err != nil {
		t.Fatal(err)
	}
	m := model{entries: entries, scores: make([]lab.Progress, len(entries)), store: lab.Store{DataDir: t.TempDir(), StateDir: t.TempDir()}, width: 96, height: 27, screen: homeScreen}
	m.selectLesson(0)
	press := func(msg tea.KeyMsg) tea.Cmd { updated, cmd := m.Update(msg); m = updated.(model); return cmd }
	if !strings.Contains(m.View(), "WELCOME") {
		t.Fatal("missing first session")
	}
	press(tea.KeyMsg{Type: tea.KeyEnter})
	if m.screen != readScreen {
		t.Fatal("start did not teach before practice")
	}
	press(tea.KeyMsg{Type: tea.KeyTab})
	if m.screen != taskScreen {
		t.Fatal("Do not reachable")
	}
	if _, err := m.store.Open(m.lesson); err != nil {
		t.Fatal(err)
	}
	press(tea.KeyMsg{Type: tea.KeyCtrlK})
	press(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("inspect"), Paste: true})
	cmd := press(tea.KeyMsg{Type: tea.KeyEnter})
	if cmd == nil || !m.busy {
		t.Fatal("inspection should run asynchronously")
	}
	updated, _ := m.Update(cmd())
	m = updated.(model)
	if m.busy || !strings.Contains(m.View(), "WORKSPACE INSPECTION") {
		t.Fatal("inspection did not finish")
	}
	press(tea.KeyMsg{Type: tea.KeyEsc})
	if m.screen != taskScreen {
		t.Fatal("inspection lost its return context")
	}
	press(tea.KeyMsg{Type: tea.KeyCtrlK})
	press(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("connected")})
	press(tea.KeyMsg{Type: tea.KeyEnter})
	if m.screen != projectsScreen {
		t.Fatal("projects unavailable")
	}
	press(tea.KeyMsg{Type: tea.KeyEnter})
	if m.lesson.ID != "streams/log-report" || m.screen != readScreen {
		t.Fatal("project did not open first step")
	}
}

func TestNewViewsFitSmallAndASCIIWindows(t *testing.T) {
	entries, err := course.Entries()
	if err != nil {
		t.Fatal(err)
	}
	m := model{entries: entries, scores: make([]lab.Progress, len(entries)), store: lab.Store{StateDir: t.TempDir()}}
	m.selectLesson(0)
	for _, screen := range []screen{homeScreen, actionsScreen, progressScreen, projectsScreen, inspectionScreen} {
		for _, size := range [][2]int{{24, 12}, {40, 18}, {96, 27}} {
			m.screen = screen
			m.width = size[0]
			m.height = size[1]
			m.ascii = true
			view := m.View()
			if lipgloss.Width(view) > m.width || lipgloss.Height(view) > m.height {
				t.Errorf("screen %d overflow %v", screen, size)
			}
			if strings.ContainsAny(view, "╭╮╯╰│─█░▸◈") {
				t.Fatal("ASCII view retained decorative Unicode")
			}
		}
	}
}
