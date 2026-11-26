package tui

import (
	"fmt"
	"slices"
	"strings"
	"unicode"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

func (m model) matches() []int {
	terms := strings.Fields(strings.ToLower(m.query))
	var matches []int
	for i, entry := range m.entries {
		lesson := entry.Lesson
		searchable := []string{entry.Workshop, lesson.Title, lesson.ID, strings.Join(lesson.Tools, " ")}
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
	pos := slices.Index(matches, m.menu)
	m.selectLesson(matches[(pos+direction+len(matches))%len(matches)])
}

func (m model) searchLabel() string {
	query := m.query
	if m.searching {
		query += "▏"
	}
	return m.paint("87", "/ "+query, true)
}

func (m model) searchKey(key tea.KeyMsg) model {
	switch key.Type {
	case tea.KeyEsc:
		m.query, m.searching = "", false
		return m
	case tea.KeyEnter:
		if len(m.matches()) > 0 {
			m.searching, m.screen, m.scroll = false, demoScreen, 0
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
	if len(matches) > 0 && !slices.Contains(matches, m.menu) {
		m.selectLesson(matches[0])
	}
	return m
}

func (m *model) appendSearch(text string) {
	text = strings.Map(func(r rune) rune {
		if unicode.IsSpace(r) {
			return ' '
		}
		if unicode.IsPrint(r) {
			return r
		}
		return -1
	}, text)
	runes := []rune(m.query + text)
	m.query = string(runes[:min(len(runes), 80)])
}

func (m model) splitCatalog() bool {
	return m.width >= 80 && m.height >= 15
}

func (m model) catalogView(width int) string {
	leftWidth := width * 38 / 100
	rightWidth := width - leftWidth - 2
	inside := leftWidth - 4

	complete := 0
	for _, score := range m.scores {
		if score.Best == 100 {
			complete++
		}
	}
	matches := m.matches()
	progress := fmt.Sprintf("COURSE PROGRESS  %d/%d complete", complete, len(m.entries))
	if lipgloss.Width(progress) > inside {
		progress = fmt.Sprintf("Progress  %d/%d complete", complete, len(m.entries))
	}
	heading := []string{m.paint("213", "LEARNING PATHS", true), progress}
	if m.searching || m.query != "" {
		heading = append(heading, m.searchLabel())
	}
	crop := lipgloss.NewStyle().MaxWidth(inside).Render
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
			rows = append(rows, crop(m.paint("251", m.workshopLabel(workshop), true)))
		}
		label := entry.Lesson.Title
		if i == m.menu {
			selected = len(rows)
			rows = append(rows, crop(m.paint("87", "▸ "+label, true)))
		} else {
			rows = append(rows, crop("  "+label))
		}
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

func (m model) preview(found bool, width int) string {
	if !found || len(m.entries) == 0 {
		return "Search by lesson, workshop, tool, or command."
	}
	entry := m.entries[m.menu]
	parts := []string{
		m.paint("213", strings.ToUpper(entry.Workshop), true),
		m.lesson.Title,
		fmt.Sprintf("%d demo steps · %d scoring criteria", len(m.lesson.Demo), len(m.lesson.Checks)),
	}
	if m.progress.Checks > 0 {
		parts = append(parts, fmt.Sprintf("Latest %d/100 · Best %d/100", m.progress.Latest.Score, m.progress.Best))
	} else {
		parts = append(parts, "Ready to start")
	}
	if m.height >= 21 {
		parts = append(parts, "", m.lesson.Task)
	}
	if m.height >= 25 && len(m.lesson.Demo) > 0 {
		frame := m.lesson.Demo[0]
		parts = append(parts, "", m.paint("221", "$ "+frame.Command, true))
		// Leave enough room for the panel border at typical laptop heights.
		if m.height >= 29 || lipgloss.Width(m.lesson.Task) <= width-4 {
			parts = append(parts, m.paint("87", frame.Visual, false))
		}
	}
	return strings.Join(parts, "\n")
}
