package tui

import (
	"strings"
	"testing"

	"cli-silly/course"
	"cli-silly/internal/lab"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/x/ansi"
)

func TestSectionsSkipHeadingsInsideCodeAndHandleDuplicateTitles(t *testing.T) {
	body := "Opening paragraph.\n\n## One\n\n```bash\n## Shell comment\n```\n\n## One\n\nOverview\n========\n"
	sections := lessonSections(body)
	want := []string{"overview", "one", "one-2", "overview-2"}
	if len(sections) != len(want) {
		t.Fatalf("unexpected sections: %+v", sections)
	}
	for i, name := range want {
		if sections[i].anchor != name {
			t.Fatalf("section %d: %q != %q", i, sections[i].anchor, name)
		}
	}
	if !strings.Contains(sections[1].markdown, "## Shell comment") {
		t.Fatal("missing fenced code")
	}
}

func TestReaderRendersTableAndFenceAtDifferentWidths(t *testing.T) {
	lesson, err := course.Load(course.First)
	if err != nil {
		t.Fatal(err)
	}
	for _, width := range []int{24, 40, 80, 110} {
		for _, color := range []bool{false, true} {
			m := model{lesson: lesson, width: width, height: 27, useColor: color, store: lab.Store{StateDir: t.TempDir()}}
			m.openReading("read-the-permission-bits")
			view := m.View()
			content := ansi.Strip(strings.Join(m.readerLines, "\n"))
			if !strings.Contains(content, "Owner") || !strings.Contains(content, "show-version") || strings.Contains(content, "```") {
				t.Errorf("width %d: missing formatted table or code:\n%s", width, content)
			}
			if lipgloss.Width(view) > width || lipgloss.Height(view) > m.height {
				t.Errorf("reader overflow at width %d", width)
			}
			if !color && strings.Contains(view, "\x1b[") {
				t.Errorf("reader emitted color escape with color disabled")
			}
		}
	}
}

func TestReaderScrollsAndJumpsBetweenSections(t *testing.T) {
	lesson, err := course.Load("foundations/cli-intro")
	if err != nil {
		t.Fatal(err)
	}
	m := model{lesson: lesson, width: 105, height: 27, store: lab.Store{StateDir: t.TempDir()}}
	m.openReading(overview)
	if view := m.View(); !strings.Contains(view, "Contents") || !strings.Contains(view, "Run pwd first") {
		t.Fatalf("the opening should show the contents and reading beyond the overview:\n%s", view)
	}
	m.scrollReading(m.readerStarts[1])
	if m.position().Anchor != "commands-and-operations" {
		t.Fatal("scrolling should save the heading at the top of the reader")
	}
	for _, key := range []tea.KeyMsg{
		{Type: tea.KeyRunes, Runes: []rune{'i'}},
		{Type: tea.KeyDown},
		{Type: tea.KeyEnter},
	} {
		updated, _ := m.Update(key)
		m = updated.(model)
	}
	if m.outline || m.position().Anchor != "why-learn-the-cli" || !strings.Contains(m.View(), "Automation") {
		t.Fatal("choosing a heading should jump to the selected part of the document")
	}
	updated, _ := m.Update(tea.WindowSizeMsg{Width: 40, Height: 20})
	m = updated.(model)
	if m.position().Anchor != "why-learn-the-cli" || lipgloss.Width(m.View()) > 40 {
		t.Fatal("resize should preserve position and fit the narrow view")
	}
}

func TestReaderPersistsSectionAndBookmark(t *testing.T) {
	entries, err := course.Entries()
	if err != nil {
		t.Fatal(err)
	}
	store := lab.Store{StateDir: t.TempDir()}
	m := model{
		store: store, width: 80, height: 27, entries: entries,
		scores: make([]lab.Progress, len(entries)),
	}
	m.selectLesson(0)
	m.openReading("deleted-heading")
	if m.position().Anchor != overview {
		t.Fatal("missing anchor did not fall back to overview")
	}
	updated, _ := m.Update(tea.KeyMsg{Type: tea.KeyRight})
	updated, _ = updated.(model).Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'b'}})
	m = updated.(model)
	saved, err := store.Reading()
	if err != nil || len(saved.Bookmarks) != 1 || saved.Last != m.position() || saved.Bookmarks[0] != saved.Last {
		t.Fatalf("saved reading state: %+v (%v)", saved, err)
	}
	m.selectLesson(1)
	m.openReading(overview)
	m.query = "paths and quoting"
	m.openPosition(saved.Bookmarks[0])
	if m.position() != saved.Last || m.query != "" {
		t.Fatal("opening the bookmark did not return to its section")
	}
	m.updateReading(true)
	if len(m.reading.Bookmarks) != 0 {
		t.Fatal("second bookmark press should remove the saved section")
	}
}
