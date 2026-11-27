package tui

import "github.com/charmbracelet/lipgloss"

type viewPosition struct {
	screen screen
	scroll int
}

func (m model) viewPosition() viewPosition {
	return viewPosition{m.screen, m.scroll}
}

func (m *model) restoreView(position viewPosition) {
	m.screen, m.scroll = position.screen, position.scroll
	if m.screen == readScreen {
		m.renderReading()
	}
}

func (m model) inLesson() bool {
	switch m.screen {
	case readScreen, demoScreen, taskScreen, resultsScreen, filesScreen:
		return true
	default:
		return false
	}
}

func (m *model) switchMode() {
	switch m.screen {
	case readScreen, demoScreen:
		m.screen, m.scroll, m.status = taskScreen, 0, ""
	case taskScreen, resultsScreen, filesScreen:
		m.stopWatch()
		m.resumeReading()
	}
}

func (m *model) back() {
	switch m.screen {
	case mapScreen:
		m.query = ""
	case readScreen:
		if m.outline {
			m.outline = false
			return
		}
		m.screen = mapScreen
		m.scroll = 0
	case demoScreen:
		m.resumeReading()
	case taskScreen:
		if m.readerLesson == m.lesson.ID {
			m.resumeReading()
		} else {
			m.screen, m.scroll, m.status = mapScreen, 0, ""
		}
	case resultsScreen:
		m.screen, m.scroll, m.status = taskScreen, 0, ""
	case filesScreen:
		m.restoreView(m.filesFrom)
	case bookmarksScreen:
		m.restoreView(m.bookmarksFrom)
	case commandsScreen:
		m.restoreView(m.helpFrom)
	case referenceDetailScreen:
		m.screen, m.scroll = referenceScreen, 0
	default:
		m.screen, m.scroll, m.status = mapScreen, 0, ""
	}
}

func (m model) adjacentLesson(direction int) int {
	if len(m.entries) == 0 {
		return -1
	}
	for i := m.menu + direction; i >= 0 && i < len(m.entries); i += direction {
		entry := m.entries[i]
		if entry.PackID == m.entries[m.menu].PackID && entry.Lesson.ReviewOf == "" {
			return i
		}
	}
	return -1
}

func (m *model) browseLesson(direction int) {
	if index := m.adjacentLesson(direction); index >= 0 {
		m.stopWatch()
		m.query = ""
		m.selectLesson(index)
		m.resumeReading()
	}
}

type keyLabel struct {
	key   string
	label string
}

func (m model) footer(width int) string {
	keys := m.visibleKeys(width)
	footer := ""
	for _, item := range keys {
		label := item.key + " " + item.label
		if footer == "" {
			footer = m.paint(teal, label, true)
			continue
		}
		if lipgloss.Width(footer)+3+lipgloss.Width(label) > width {
			break
		}
		footer += m.paint(muted, " · ", false) + m.paint(violet, label, false)
	}
	return footer
}

func (m model) visibleKeys(width int) []keyLabel {
	if m.search != nil {
		return []keyLabel{{"Type", "to filter"}, {"Enter", "Read"}, {"Esc", "Back"}, {"↑↓", "Select"}}
	}
	if m.refSearching {
		return []keyLabel{{"Type", "to filter"}, {"Enter", "Results"}, {"Esc", "Clear"}}
	}
	switch m.screen {
	case mapScreen:
		return []keyLabel{{"Enter", "Read"}, {"d", "Do"}, {"/", "Find"}, {"↑↓", "Choose"}, {"←→", "Workshops"}, {"?", "Help"}}
	case readScreen:
		if m.outline {
			return []keyLabel{{"Enter", "Jump"}, {"↑↓", "Choose"}, {"Esc", "Back"}}
		}
		if width < 55 {
			return []keyLabel{{"Tab", "Do"}, {"i", "Index"}, {"e", "Demo"}, {"↑↓", "Scroll"}, {"[ ]", "Lessons"}}
		}
		return []keyLabel{{"Tab", "Do"}, {"e", "Example"}, {"i", "Contents"}, {"↑↓", "Scroll"}, {"←→", "Sections"}, {"[ ]", "Lessons"}, {"?", "Help"}}
	case demoScreen:
		action := keyLabel{"Enter", "Next step"}
		if width < 55 {
			action.label = "Next"
		}
		if m.frame == len(m.lesson.Demo)-1 {
			action.label = "Start Do"
		}
		keys := []keyLabel{action, {"Esc", "Read"}}
		if len(m.lesson.Demo[m.frame].Permissions) == 9 {
			keys = append(keys, keyLabel{"x", "Toggle mode"})
		}
		return append(keys, keyLabel{"←→", "Steps"}, keyLabel{"d", "Do"}, keyLabel{"[ ]", "Lessons"})
	case taskScreen:
		return []keyLabel{{"Enter", "Bash"}, {"Tab", "Read"}, {"h", "Hint"}, {"e", "Files"}, {"c", "Check"}, {"[ ]", "Lessons"}, {"Esc", "Back"}}
	case resultsScreen:
		if m.progress.Latest.Complete && m.adjacentLesson(1) >= 0 {
			return []keyLabel{{"]", "Next lesson"}, {"Tab", "Read"}, {"Esc", "Task"}, {"Enter", "Bash"}, {"c", "Check"}}
		}
		if m.progress.Checks == 0 || m.progress.Latest.Complete {
			return []keyLabel{{"Enter", "Bash"}, {"Tab", "Read"}, {"Esc", "Task"}, {"c", "Check"}, {"m", "Lessons"}}
		}
		return []keyLabel{{"Enter", "Retry"}, {"h", "Hint"}, {"Tab", "Read"}, {"Esc", "Task"}, {"c", "Check"}, {"[ ]", "Lessons"}}
	case filesScreen:
		return []keyLabel{{"Enter", "Edit"}, {"Esc", "Back"}, {"↑↓", "Choose"}}
	case bookmarksScreen:
		return []keyLabel{{"Enter", "Open"}, {"Esc", "Back"}, {"↑↓", "Choose"}}
	case referenceScreen:
		return []keyLabel{{"Enter", "Preview"}, {"/", "Search"}, {"Esc", "Back"}, {"↑↓", "Choose"}}
	case referenceDetailScreen:
		return []keyLabel{{"p", "Print"}, {"Esc", "Back"}, {"↑↓", "Scroll"}}
	case quizScreen:
		return []keyLabel{{"Enter", "Submit"}, {"↑↓", "Choose"}, {"r", "Retry"}, {"Esc", "Back"}}
	case reviewScreen:
		return []keyLabel{{"Enter", "Challenge"}, {"↑↓", "Choose"}, {"Esc", "Back"}}
	default:
		return []keyLabel{{"Esc", "Back"}, {"↑↓", "Scroll"}}
	}
}
