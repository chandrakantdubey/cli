package tui

import (
	"strings"
	"testing"

	"cli-silly/course"
	"cli-silly/internal/lab"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

func TestDemoToggleAndAdvance(t *testing.T) {
	lesson, err := course.Load(course.First)
	if err != nil {
		t.Fatal(err)
	}
	m := model{lesson: lesson, width: 80, height: 24, screen: demoScreen}
	if !strings.Contains(m.View(), "permission denied") {
		t.Fatal("first frame should show why the script cannot run directly")
	}
	updated, _ := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'x'}})
	m = updated.(model)
	if !strings.Contains(m.View(), "execute bit set (preview)") {
		t.Fatal("toggle did not change the preview")
	}
	updated, _ = m.Update(tea.KeyMsg{Type: tea.KeyRight})
	m = updated.(model)
	if m.frame != 1 || m.toggleExec || !strings.Contains(m.View(), "chmod u+x") {
		t.Fatal("next frame should display the real permission change")
	}
}

func TestLayoutFitsTerminal(t *testing.T) {
	lesson, err := course.Load(course.First)
	if err != nil {
		t.Fatal(err)
	}
	for _, width := range []int{24, 40, 80} {
		for _, screen := range []screen{mapScreen, demoScreen, taskScreen, resultsScreen} {
			m := model{lesson: lesson, width: width, height: 12, screen: screen}
			m.hints = 1
			view := m.View()
			if lipgloss.Width(view) > m.width || lipgloss.Height(view) > m.height {
				t.Errorf("screen %d overflowed %dx%d terminal (%dx%d)", screen, m.width, m.height, lipgloss.Width(view), lipgloss.Height(view))
			}
			if strings.Contains(view, "\x1b[") {
				t.Errorf("screen %d emitted color with useColor disabled", screen)
			}
		}
	}
}

func TestCatalogSearchAndPreview(t *testing.T) {
	entries, err := course.Entries()
	if err != nil {
		t.Fatal(err)
	}
	m := model{
		entries: entries, scores: make([]lab.Progress, len(entries)), width: 96, height: 27,
		store: lab.Store{StateDir: t.TempDir()},
	}
	m.selectLesson(0)
	view := m.View()
	if !strings.Contains(view, "FOUNDATIONS") || !strings.Contains(view, "$ pwd") || !strings.Contains(view, "welcome-copy.txt") {
		t.Fatal("wide map should display the path list and selected lesson preview")
	}
	updated, _ := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'/'}})
	updated, _ = updated.(model).Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("Git selective"), Paste: true})
	m = updated.(model)
	if m.entries[m.catalogIndex()].Lesson.ID != "git/selective-staging" || len(m.matches()) != 1 {
		t.Fatalf("search should highlight the one matching lesson")
	}
	for _, size := range [][2]int{{40, 18}, {80, 18}, {96, 27}, {120, 32}} {
		m.width, m.height = size[0], size[1]
		view := m.View()
		if lipgloss.Width(view) > m.width || lipgloss.Height(view) > m.height {
			t.Errorf("catalog overflowed a %dx%d terminal", m.width, m.height)
		}
	}
	updated, _ = m.Update(tea.KeyMsg{Type: tea.KeyEnter})
	m = updated.(model)
	if m.screen != readScreen || m.search != nil || m.lesson.ID != "git/selective-staging" {
		t.Fatal("Enter should open the selected search result")
	}
	updated, _ = m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'d'}})
	m = updated.(model)
	if m.screen != taskScreen || !strings.Contains(m.View(), "YOUR TASK") {
		t.Fatal("Do should show the selected lesson's task before starting Bash")
	}
}

func TestLessonNavigationKeepsYourPlace(t *testing.T) {
	entries, err := course.Entries()
	if err != nil {
		t.Fatal(err)
	}
	m := model{
		entries: entries, scores: make([]lab.Progress, len(entries)), width: 110, height: 30,
		store: lab.Store{StateDir: t.TempDir()},
	}
	m.selectLesson(0)
	m.openReading(overview)
	m.jumpSection(1)
	m.scrollReading(3)
	start := m.readerScroll
	position := m.position()
	press := func(key tea.KeyMsg) {
		t.Helper()
		updated, _ := m.Update(key)
		m = updated.(model)
	}
	letter := func(r string) tea.KeyMsg { return tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune(r)} }
	assertReading := func() {
		t.Helper()
		if m.screen != readScreen || m.position() != position || m.readerScroll != start {
			t.Fatalf("reading did not resume at %+v, line %d", position, start)
		}
	}

	press(tea.KeyMsg{Type: tea.KeyTab})
	if m.screen != taskScreen || !strings.Contains(m.View(), "What earns points") {
		t.Fatal("Tab should show the task before opening Bash")
	}
	press(letter("e"))
	press(tea.KeyMsg{Type: tea.KeyEsc})
	if m.screen != taskScreen {
		t.Fatal("closing files should return to the current task")
	}
	press(tea.KeyMsg{Type: tea.KeyShiftTab})
	assertReading()

	press(letter("e"))
	press(tea.KeyMsg{Type: tea.KeyEsc})
	assertReading()
	press(letter("/"))
	press(letter("git"))
	press(tea.KeyMsg{Type: tea.KeyEsc})
	assertReading()

	press(letter("]"))
	if m.lesson.ID != "foundations/mental-model" || m.screen != readScreen {
		t.Fatal("next lesson should open in Read")
	}
	press(letter("["))
	assertReading()

	press(letter("m"))
	press(tea.KeyMsg{Type: tea.KeyRight})
	if m.entries[m.menu].Workshop != "Core tools" {
		t.Fatal("right should select the first lesson in the next workshop")
	}
	press(tea.KeyMsg{Type: tea.KeyLeft})
	if m.menu != 0 {
		t.Fatal("left should select the first lesson in the previous workshop")
	}
}

func TestReadAndDoActionsFitNarrowTerminal(t *testing.T) {
	lesson, err := course.Load("foundations/cli-intro")
	if err != nil {
		t.Fatal(err)
	}
	m := model{lesson: lesson, width: 24, height: 12, store: lab.Store{StateDir: t.TempDir()}}
	m.openReading(overview)
	if view := m.View(); !strings.Contains(view, "Tab Do") || lipgloss.Width(view) > 24 {
		t.Fatalf("the Read action bar must fit: %s", view)
	}
	m.switchMode()
	if view := m.View(); !strings.Contains(view, "Tab Read") || lipgloss.Width(view) > 24 {
		t.Fatalf("the Do action bar must fit: %s", view)
	}
}
