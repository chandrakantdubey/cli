package tui

import (
	"fmt"
	"slices"
	"strings"
	"unicode"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/x/ansi"
)

type catalogSearch struct {
	from   viewPosition
	query  string
	status string
	index  int
}

func (m model) matches() []int {
	terms := strings.Fields(strings.ToLower(m.query))
	var matches []int
	for i, entry := range m.entries {
		lesson := entry.Lesson
		searchable := []string{entry.Workshop, lesson.Title, lesson.ID, lesson.ReferenceID, strings.Join(lesson.Tools, " ")}
		for _, frame := range lesson.Demo {
			searchable = append(searchable, frame.Command)
		}
		text := strings.ToLower(strings.Join(searchable, " "))
		if !slices.ContainsFunc(terms, func(term string) bool { return !strings.Contains(text, term) }) {
			matches = append(matches, i)
		}
	}
	return matches
}

func (m *model) moveSelection(direction int) {
	matches := m.matches()
	if len(matches) == 0 {
		return
	}
	pos := slices.Index(matches, m.catalogIndex())
	index := matches[(pos+direction+len(matches))%len(matches)]
	if m.search != nil {
		m.search.index = index
	} else {
		m.selectLesson(index)
	}
}

func (m model) catalogIndex() int {
	if m.search != nil {
		return m.search.index
	}
	return m.menu
}

func (m *model) moveWorkshop(direction int) {
	matches := m.matches()
	position := slices.Index(matches, m.menu)
	for next := position + direction; next >= 0 && next < len(matches); next += direction {
		entry := m.entries[matches[next]]
		if entry.Workshop == m.entries[m.menu].Workshop && entry.PackID == m.entries[m.menu].PackID {
			continue
		}
		for next > 0 {
			previous := m.entries[matches[next-1]]
			if previous.Workshop != entry.Workshop || previous.PackID != entry.PackID {
				break
			}
			next--
		}
		m.selectLesson(matches[next])
		return
	}
}

func (m model) searchLabel() string {
	query := m.query
	if m.search != nil {
		query += "▏"
	}
	return m.paint(teal, "/ "+query, true)
}

func (m model) searchKey(key tea.KeyMsg) model {
	search := *m.search
	m.search = &search
	switch key.Type {
	case tea.KeyEsc:
		m.query, m.search = search.query, nil
		m.status = search.status
		m.restoreView(search.from)
		return m
	case tea.KeyEnter:
		if len(m.matches()) > 0 {
			m.search = nil
			m.selectLesson(search.index)
			m.resumeReading()
		}
		return m
	case tea.KeyUp:
		m.moveSelection(-1)
		return m
	case tea.KeyDown:
		m.moveSelection(1)
		return m
	case tea.KeyBackspace:
		runes := []rune(m.query)
		if len(runes) > 0 {
			m.query = string(runes[:len(runes)-1])
		}
	case tea.KeySpace:
		m.appendSearch(" ")
	case tea.KeyRunes:
		m.appendSearch(string(key.Runes))
	default:
		return m
	}
	matches := m.matches()
	if len(matches) > 0 && !slices.Contains(matches, search.index) {
		m.search.index = matches[0]
	}
	return m
}

func (m *model) appendSearch(text string) {
	m.query = appendQuery(m.query, text)
}

func appendQuery(query, text string) string {
	text = strings.Map(func(r rune) rune {
		if unicode.IsSpace(r) {
			return ' '
		}
		if unicode.IsPrint(r) {
			return r
		}
		return -1
	}, text)
	runes := []rune(query + text)
	return string(runes[:min(len(runes), 80)])
}

func (m model) splitCatalog() bool {
	return m.width >= 80 && m.height >= 15
}

func (m model) catalogView(width int) string {
	leftWidth := width * 42 / 100
	rightWidth := width - leftWidth - 2
	inside := leftWidth - 4

	complete, total := m.requiredProgress()
	matches := m.matches()
	progress := fmt.Sprintf("COURSE PROGRESS  %d/%d complete", complete, total)
	if lipgloss.Width(progress) > inside {
		progress = fmt.Sprintf("Progress  %d/%d complete", complete, total)
	}
	heading := []string{m.paint(violet, "LEARNING PATHS", true), progress}
	if m.search != nil || m.query != "" {
		heading = append(heading, m.searchLabel())
	}
	crop := func(label string) string { return ansi.Truncate(label, inside, "…") }
	for i, line := range heading {
		heading[i] = crop(line)
	}

	rows := make([]string, 0, len(matches)*2)
	workshop := ""
	selected := 0
	for _, i := range matches {
		entry := m.entries[i]
		if entry.Workshop != workshop {
			workshop = entry.Workshop
			rows = append(rows, crop(m.paint(muted, m.workshopLabel(workshop), true)))
		}
		if i == m.catalogIndex() {
			selected = len(rows)
		}
		rows = append(rows, m.lessonRow(i, inside))
	}
	if len(rows) == 0 {
		rows = append(rows, "No matching lessons.")
	}
	n := max(1, m.bodyHeight()-len(heading)-3)
	start := max(0, min(selected-n+1, len(rows)-n))
	left := m.panel(leftWidth, strings.Join(heading, "\n")+"\n\n"+strings.Join(rows[start:min(start+n, len(rows))], "\n"))
	preview := lipgloss.NewStyle().Width(rightWidth - 4).Render(m.preview(len(matches) != 0, rightWidth))
	previewLines := strings.Split(preview, "\n")
	right := m.panel(rightWidth, strings.Join(previewLines[:min(len(previewLines), m.bodyHeight()-2)], "\n"))
	return lipgloss.JoinHorizontal(lipgloss.Top, left, "  ", right)
}

func (m model) lessonRow(index, width int) string {
	prefix := "  "
	score := m.scores[index]
	status, color := "○", muted
	switch {
	case score.Best == 100:
		status, color = "✓", green
	case score.Checks > 0:
		status, color = "◐", amber
	}
	title := ansi.Truncate(m.entries[index].Lesson.Title, max(1, width-4), "…")
	if index == m.catalogIndex() {
		return m.selectedRow("▸ "+status+" "+title, width)
	}
	return prefix + m.paint(color, status, false) + " " + title
}

func (m model) preview(found bool, width int) string {
	if !found || len(m.entries) == 0 {
		return "Search by lesson, workshop, tool, or command."
	}
	index := m.catalogIndex()
	entry, progress := m.entries[index], m.scores[index]
	lesson := entry.Lesson
	parts := []string{
		m.paint(violet, strings.ToUpper(entry.Workshop), true),
		m.paint(teal, lesson.Title, true),
		fmt.Sprintf("%d example steps  ·  %d scoring criteria", len(lesson.Demo), len(lesson.Checks)),
	}
	if progress.Checks > 0 {
		parts = append(parts, m.score(progress.Best)+" best")
	} else {
		parts = append(parts, "Ready to start")
	}
	if m.height >= 21 {
		parts = append(parts, "", wrapTask(lesson.Task, width-4))
	}
	if m.height >= 25 && len(lesson.Demo) > 0 {
		frame := lesson.Demo[0]
		parts = append(parts, "", m.paint(amber, "$ "+frame.Command, true))
		// Leave enough room for the panel border at typical laptop heights.
		if m.height >= 29 || lipgloss.Width(lesson.Task) <= width-4 {
			parts = append(parts, m.paint(teal, frame.Visual, false))
		}
	}
	return strings.Join(parts, "\n")
}
