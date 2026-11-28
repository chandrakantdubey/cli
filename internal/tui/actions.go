package tui

import (
	"fmt"
	"strings"
	"time"

	"cli-silly/course"

	tea "github.com/charmbracelet/bubbletea"
)

type action struct{ label, key string }

type inspectionReady struct {
	body string
	err  error
}

func (m model) actions() []action {
	items := []action{{"Read this lesson", "read"}, {"Do this lesson", "do"}, {"Browse all lessons", "lessons"}, {"Resume last reading", "resume"}, {"Progress and skills", "progress"}, {"Inspect workspace files and Git", "inspect"}, {"Find archived attempts", "archives"}, {"Hints for this task", "hint"}, {"Check workspace", "check"}, {"Related practice exercises", "practice"}, {"Review queue", "review"}, {"Reference commands", "references"}, {"Welcome and getting started", "home"}, {"Keyboard help", "help"}}
	items = append(items, action{"Connected developer projects", "projects"})
	var found []action
	for _, item := range items {
		if strings.Contains(strings.ToLower(item.label), strings.ToLower(m.actionQuery)) {
			found = append(found, item)
		}
	}
	return found
}

func (m model) actionKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	items := m.actions()
	switch msg.String() {
	case "esc", "ctrl+k":
		m.restoreView(m.actionFrom)
	case "up":
		m.actionIndex = max(0, m.actionIndex-1)
	case "down":
		m.actionIndex = max(0, min(m.actionIndex+1, len(items)-1))
	case "backspace":
		r := []rune(m.actionQuery)
		if len(r) > 0 {
			m.actionQuery = string(r[:len(r)-1])
		}
		m.actionIndex = 0
	case "enter":
		if len(items) == 0 {
			return m, nil
		}
		key := items[m.actionIndex].key
		m.restoreView(m.actionFrom)
		m.stopWatch()
		m.status = ""
		m.scroll = 0
		switch key {
		case "read":
			m.resumeReading()
		case "do":
			m.screen = taskScreen
		case "lessons":
			m.screen = mapScreen
		case "resume":
			if m.reading.Last.LessonID != "" {
				m.openPosition(m.reading.Last)
			} else {
				m.resumeReading()
			}
		case "home":
			m.screen = homeScreen
		case "projects":
			m.refreshScores()
			m.detailFrom = m.viewPosition()
			m.screen = projectsScreen
			m.projectIndex = 0
		case "progress":
			m.refreshScores()
			m.detailFrom = m.viewPosition()
			m.screen = progressScreen
		case "hint":
			m.screen = taskScreen
			return m.onScreenKey("h")
		case "check":
			m.screen = taskScreen
			return m.onScreenKey("c")
		case "practice":
			m.openRelatedPractice()
		case "review":
			return m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'V'}})
		case "references":
			return m.openReferences()
		case "help":
			m.helpFrom = m.viewPosition()
			m.screen = commandsScreen
		case "inspect", "archives":
			m.detailFrom = m.viewPosition()
			m.screen = inspectionScreen
			if key == "inspect" {
				m.busy = true
				m.status = "Inspecting workspace…"
				m.inspection = ""
				return m, func() tea.Msg { body, err := m.store.Inspect(m.lesson); return inspectionReady{body, err} }
			} else {
				paths, err := m.store.Archives(m.lesson)
				m.inspection = "ARCHIVED ATTEMPTS\n" + m.lesson.Title + "\n\n"
				if err != nil {
					m.inspection += err.Error()
				} else if len(paths) == 0 {
					m.inspection += "No archived attempts. Reset archives your work and keeps earned scores."
				} else {
					m.inspection += strings.Join(paths, "\n") + "\n\nUse cli-silly restore <lesson-id> <.saved-name> from your normal shell. The current attempt will also be archived."
				}
			}
		}
	default:
		if msg.Type == tea.KeyRunes {
			m.actionQuery = appendQuery(m.actionQuery, string(msg.Runes))
			m.actionIndex = 0
		}
	}
	return m, nil
}

func (m model) projectsView(width int) string {
	projects, err := course.Projects()
	if err != nil {
		return err.Error()
	}
	lines := []string{m.paint(violet, "CONNECTED PROJECTS", true), "Apply familiar skills to a development workflow.", ""}
	for i, p := range projects {
		label := "  " + p.Title
		if i == m.projectIndex {
			label = m.paint(teal, "▸ "+p.Title, true)
		}
		lines = append(lines, label)
	}
	p := projects[m.projectIndex]
	aside := m.paint(violet, p.Title, true) + "\n\n" + p.Outcome + "\n"
	for _, id := range p.Lessons {
		for i, e := range m.entries {
			if e.Lesson.ID == id {
				aside += fmt.Sprintf("\n%d/100  %s", m.scores[i].Best, e.Lesson.Title)
			}
		}
	}
	aside += "\n\nEnter opens the first unfinished step. Each step has its own workspace and score. Return through Ctrl+k to choose the next step."
	return m.lessonPanels(width, strings.Join(lines, "\n"), aside)
}

func (m *model) startProject() {
	projects, err := course.Projects()
	if err != nil {
		m.status = err.Error()
		return
	}
	selected := -1
	for _, id := range projects[m.projectIndex].Lessons {
		for i, e := range m.entries {
			if e.Lesson.ID == id {
				if selected < 0 {
					selected = i
				}
				if m.scores[i].Best < 100 {
					m.selectLesson(i)
					m.resumeReading()
					return
				}
			}
		}
	}
	if selected >= 0 {
		m.selectLesson(selected)
		m.resumeReading()
	}
}

func (m model) actionsView() string {
	lines := []string{m.paint(violet, "ACTIONS", true), "Type to find an action: " + m.actionQuery, ""}
	for i, a := range m.actions() {
		label := "  " + a.label
		if i == m.actionIndex {
			label = m.paint(teal, "▸ "+a.label, true)
		}
		lines = append(lines, label)
	}
	if len(m.actions()) == 0 {
		lines = append(lines, "No matching actions. Backspace to edit; Esc to return.")
	}
	return strings.Join(lines, "\n")
}

func (m model) homeView(width int) string {
	title := "WELCOME TO CLI-SILLY"
	start := "Enter  Start " + m.lesson.Title
	if m.reading.Last.LessonID != "" {
		title = "WELCOME BACK"
		start = "Enter  Continue your last reading"
	}
	complete, total := m.requiredProgress()
	main := m.paint(violet, title, true) + "\n\nLearn Linux developer tools by reading, using them, and checking your work.\n\n" + m.paint(teal, start, true) + "\nm      Browse all lessons\nCtrl+k Find an action\n\n" + fmt.Sprintf("%d/%d core labs complete", complete, total)
	aside := m.paint(violet, "YOUR FIRST SESSION", true) + "\n\n1. Read the explanation and example.\n2. Tab opens Do: review the task.\n3. Enter opens a real Bash workspace.\n4. Run cli-silly check for feedback.\n5. Type exit to return here.\n\nHints are welcome. Your work is saved.\n\nBash has your normal Linux user permissions. Practice in lab/; it is not an isolated system."
	return m.lessonPanels(width, main, aside)
}

func (m model) progressView() string {
	lines := []string{m.paint(violet, "YOUR PRACTICE", true), "Scores record exercise results. Reading and reviews are tracked separately.", ""}
	for i, e := range m.entries {
		p := m.scores[i]
		label := fmt.Sprintf("%s  %d/100", e.Lesson.Title, p.Best)
		if p.Checks == 0 {
			label = e.Lesson.Title + "  Not attempted"
		}
		if p.PreviousBest > 0 {
			label += fmt.Sprintf(" (earlier rubric %d/100 retained)", p.PreviousBest)
		}
		if e.Lesson.ReviewOf != "" {
			label = "Practice: " + label
			if p.ReviewAt != nil {
				due := "Due now"
				if p.ReviewAt.After(time.Now()) {
					due = "Due " + p.ReviewAt.Local().Format("Jan 2")
				}
				label += " · " + due
			}
		}
		lines = append(lines, label)
	}
	if m.reading.Last.LessonID != "" {
		lines = append(lines, "", "Reading bookmark: "+m.reading.Last.LessonID+" / "+m.reading.Last.Anchor)
	}
	lines = append(lines, "", "Complete linked practice to revisit a skill. Spaced passes grow the review interval from 1 to 3, 7, 14, and 30 days. Early practice is always available.")
	return strings.Join(lines, "\n")
}

// Keep score records fresh when opening the dashboard after an external check.
func (m *model) refreshScores() {
	for i, e := range m.entries {
		if p, err := m.store.Progress(e.Lesson); err == nil {
			m.scores[i] = p
		}
	}
	if m.menu < len(m.scores) {
		m.progress = m.scores[m.menu]
	}
}
