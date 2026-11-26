package tui

import (
	"fmt"
	"strings"

	"cli-silly/course"

	tea "github.com/charmbracelet/bubbletea"
)

func (m model) referenceMatches() []course.Reference {
	var matches []course.Reference
	for _, entry := range m.references {
		if entry.Matches(m.refQuery) {
			matches = append(matches, entry)
		}
	}
	return matches
}

func (m model) openReferences() (tea.Model, tea.Cmd) {
	m.stopWatch()
	references, err := course.References()
	if err != nil {
		m.status = err.Error()
		return m, nil
	}
	m.references, m.refQuery = references, ""
	m.screen, m.refIndex, m.scroll, m.status = referenceScreen, 0, 0, ""
	m.refSearching = false
	return m, nil
}

func (m model) referenceSearch(key tea.KeyMsg) model {
	switch key.Type {
	case tea.KeyEnter:
		m.refSearching = false
	case tea.KeyEsc:
		m.refSearching = false
		m.refQuery = ""
	case tea.KeyBackspace:
		runes := []rune(m.refQuery)
		if len(runes) > 0 {
			m.refQuery = string(runes[:len(runes)-1])
		}
	case tea.KeySpace:
		m.refQuery = appendQuery(m.refQuery, " ")
	case tea.KeyRunes:
		m.refQuery = appendQuery(m.refQuery, string(key.Runes))
	}
	m.refIndex, m.scroll = 0, 0
	return m
}

func (m model) referenceView() string {
	head := []string{m.paint("213", "REFERENCE LIBRARY", true), ""}
	if m.refSearching || m.refQuery != "" {
		cursor := ""
		if m.refSearching {
			cursor = "▏"
		}
		head = append(head, "/ "+m.refQuery+cursor, "")
	}
	if m.screen == referenceScreen {
		matches := m.referenceMatches()
		if len(matches) == 0 {
			return strings.Join(append(head, "No matching references."), "\n")
		}
		for i, entry := range matches {
			label := fmt.Sprintf("%s · %s", entry.Title, strings.Join(entry.Tools, ", "))
			if i == m.refIndex {
				label = m.paint("87", "▸ "+label, true)
			} else {
				label = "  " + label
			}
			head = append(head, label)
		}
		return strings.Join(head, "\n")
	}

	entry := m.activeReference
	code, _ := entry.Format(nil)
	head = append(head, entry.Title, "", m.paint("87", strings.TrimSuffix(code, "\n"), false), "",
		"Example: "+entry.Example, "Expected: "+entry.Expected, "Learn: cli-silly learn "+entry.LessonID,
		"Print with parameters: cli-silly snippet "+entry.ID)
	for _, parameter := range entry.Parameters {
		head = append(head, "  "+parameter.Name+" = "+parameter.Default)
	}
	if entry.Kind == "template" {
		head = append(head, "", "Generate inside a lab: cli-silly template "+entry.ID+" starter.sh")
	}
	return strings.Join(head, "\n")
}
