package tui

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"

	"cli-silly/course"
	"cli-silly/internal/lab"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type screen int

const (
	mapScreen screen = iota
	demoScreen
	taskScreen
	resultsScreen
	readScreen
	commandsScreen
	bookmarksScreen
)

type model struct {
	lesson        course.Lesson
	entries       []course.Entry
	scores        []lab.Progress
	store         lab.Store
	previous      screen
	busy          bool
	progress      lab.Progress
	openShell     func(course.Lesson) (*exec.Cmd, error)
	screen        screen
	menu          int
	frame         int
	hints         int
	scroll        int
	width         int
	height        int
	flash         int
	frameToken    int
	toggleExec    bool
	useColor      bool
	status        string
	searching     bool
	query         string
	reading       lab.ReadingState
	sections      []section
	sectionIndex  int
	readerText    string
	bookmarkIndex int
}

type shellExited struct{ err error }
type graded struct{ err error }
type pulse struct{ token int }

func Run(entries []course.Entry, lessonID, start string, useColor bool, store lab.Store, openShell func(course.Lesson) (*exec.Cmd, error)) error {
	if len(entries) == 0 {
		return fmt.Errorf("the course contains no lessons")
	}
	m := model{
		entries: entries, scores: make([]lab.Progress, len(entries)), openShell: openShell, store: store,
		width: 80, height: 24, useColor: useColor,
	}
	for i, entry := range entries {
		progress, err := store.Progress(entry.Lesson)
		if err != nil {
			return fmt.Errorf("%s: %w", entry.Lesson.ID, err)
		}
		m.scores[i] = progress
		if entry.Lesson.ID == lessonID {
			m.menu = i
		}
	}
	m.selectLesson(m.menu)
	reading, err := store.Reading()
	if err != nil {
		m.status = err.Error()
	} else {
		m.reading = reading
	}
	switch start {
	case "learn", "continue":
		m.screen = demoScreen
	case "resume":
		anchor := ""
		if m.reading.Last.LessonID == lessonID {
			anchor = m.reading.Last.Anchor
		}
		m.openReading(anchor)
	}
	_, err = tea.NewProgram(m, tea.WithAltScreen(), tea.WithInput(os.Stdin), tea.WithOutput(os.Stdout)).Run()
	return err
}

func (m *model) selectLesson(index int) {
	m.menu = index
	m.lesson, m.progress = m.entries[index].Lesson, m.scores[index]
	m.frame, m.flash, m.hints, m.scroll = 0, 0, 0, 0
	m.toggleExec, m.status = false, ""
}

func (m model) Init() tea.Cmd { return nil }

func (m model) Update(message tea.Msg) (tea.Model, tea.Cmd) {
	switch message := message.(type) {
	case tea.WindowSizeMsg:
		m.width, m.height = message.Width, message.Height
		if m.screen == readScreen || (m.screen == commandsScreen && m.previous == readScreen) {
			m.renderSection()
			if m.screen == readScreen {
				m.scroll = min(m.scroll, max(0, len(m.contentLines())-m.bodyHeight()))
			}
		}
	case shellExited:
		return m.showResults(message.err), nil
	case graded:
		return m.showResults(message.err), nil
	case pulse:
		if message.token == m.frameToken && m.flash > 0 {
			m.flash--
			if m.flash > 0 {
				return m, animate(message.token)
			}
		}
	case tea.KeyMsg:
		key := message.String()
		if m.searching && key != "ctrl+c" {
			return m.searchKey(message), nil
		}
		if m.busy && key != "ctrl+c" {
			return m, nil
		}
		switch key {
		case "q", "ctrl+c":
			return m, tea.Quit
		case "m":
			m.screen, m.scroll, m.status = mapScreen, 0, ""
		case "B":
			m.screen, m.scroll, m.bookmarkIndex, m.status = bookmarksScreen, 0, 0, ""
		case "/":
			m.screen, m.scroll, m.searching, m.status = mapScreen, 0, true, ""
		case "?":
			if m.screen != commandsScreen {
				m.previous, m.screen, m.scroll = m.screen, commandsScreen, 0
			} else {
				m.screen, m.scroll = m.previous, 0
			}
		case "esc":
			switch m.screen {
			case mapScreen:
				if m.query == "" {
					return m, tea.Quit
				}
				m.query = ""
			case taskScreen:
				m.screen = demoScreen
			case commandsScreen:
				m.screen = m.previous
			default:
				m.screen = mapScreen
			}
			m.scroll = 0
		default:
			return m.onScreenKey(key)
		}
	}
	return m, nil
}

func (m model) showResults(actionErr error) model {
	m.screen, m.scroll = resultsScreen, 0
	m.busy = false
	progress, err := m.store.Progress(m.lesson)
	if err == nil {
		m.progress = progress
		m.scores[m.menu] = progress
	}
	if actionErr != nil {
		m.status = actionErr.Error()
	} else if err != nil {
		m.status = err.Error()
	} else {
		m.status = ""
	}
	return m
}

func (m model) onScreenKey(key string) (tea.Model, tea.Cmd) {
	if key == "up" || key == "down" {
		direction := 1
		if key == "up" {
			direction = -1
		}
		switch m.screen {
		case mapScreen:
			m.moveSelection(direction)
		case bookmarksScreen:
			m.bookmarkIndex = max(0, min(m.bookmarkIndex+direction, len(m.reading.Bookmarks)-1))
		default:
			m.scroll = max(0, min(m.scroll+direction, len(m.contentLines())-m.bodyHeight()))
		}
		return m, nil
	}
	switch m.screen {
	case mapScreen:
		if len(m.matches()) == 0 {
			return m, nil
		}
		switch key {
		case "enter", " ":
			m.screen = demoScreen
		case "t":
			m.screen = taskScreen
		case "s":
			m.screen = resultsScreen
		case "l":
			m.openReading(overview)
			return m, nil
		case "g":
			if m.reading.Last.LessonID == "" {
				m.openReading(overview)
			} else {
				m.openPosition(m.reading.Last)
			}
			return m, nil
		}
		m.scroll, m.status = 0, ""
	case readScreen:
		switch key {
		case "enter":
			m.screen, m.scroll = demoScreen, 0
		case "t":
			m.screen, m.scroll = taskScreen, 0
		case "left", "[":
			m.moveSection(-1)
		case "right", "]":
			m.moveSection(1)
		case "b":
			m.updateReading(true)
		case "pgup":
			m.scroll = max(0, m.scroll-m.bodyHeight())
		case "pgdown":
			m.scroll = min(max(0, len(m.contentLines())-m.bodyHeight()), m.scroll+m.bodyHeight())
		}
	case bookmarksScreen:
		if key == "enter" && len(m.reading.Bookmarks) > 0 {
			m.openPosition(m.reading.Bookmarks[m.bookmarkIndex])
		}
	case demoScreen:
		switch key {
		case "l":
			m.openReading(overview)
		case "left", "b":
			m.frame, m.scroll, m.flash, m.toggleExec = max(0, m.frame-1), 0, 0, false
		case "x":
			if len(m.lesson.Demo[m.frame].Permissions) == 9 {
				m.toggleExec = !m.toggleExec
			}
		case "t":
			m.screen, m.scroll = taskScreen, 0
		case "right", "enter", " ":
			m.scroll = 0
			if m.frame == len(m.lesson.Demo)-1 {
				m.screen = taskScreen
			} else {
				m.frame++
				m.toggleExec = false
				if m.useColor {
					m.flash, m.frameToken = 3, m.frameToken+1
					return m, animate(m.frameToken)
				}
			}
		}
	case taskScreen, resultsScreen:
		switch key {
		case "h":
			m.screen = taskScreen
			m.hints = m.hints%len(m.lesson.Hints) + 1
			m.scroll = max(0, len(m.contentLines())-m.bodyHeight())
		case "c":
			m.busy = true
			m.status = "Checking workspace…"
			return m, func() tea.Msg {
				if _, err := m.store.Open(m.lesson); err != nil {
					return graded{err}
				}
				_, err := m.store.Check(m.lesson, "")
				return graded{err}
			}
		case "r":
			m.toggleExec = false
			m.screen, m.frame, m.scroll, m.flash = demoScreen, 0, 0, 0
		case "enter", "o":
			cmd, err := m.openShell(m.lesson)
			if err != nil {
				m.status = err.Error()
				return m, nil
			}
			m.status = ""
			return m, tea.ExecProcess(cmd, func(err error) tea.Msg { return shellExited{err} })
		}
	}
	return m, nil
}

func animate(token int) tea.Cmd {
	return tea.Tick(150*time.Millisecond, func(time.Time) tea.Msg { return pulse{token} })
}

func (m model) View() string {
	if m.width < 24 || m.height < 8 {
		return "Resize terminal · q quit"
	}
	width := min(m.width-2, 100)
	header := m.paint("213", "◈  CLI-SILLY", true)
	if width >= 32 {
		title := "WORKSHOPS"
		if m.screen != mapScreen && len(m.entries) > 0 {
			title = m.entries[m.menu].Workshop
		}
		header += "  " + m.paint("251", title, false)
	}
	tabs := []string{"Map", "01 Learn", "02 Do", "03 Score", "Read", "Commands", "Bookmarks"}
	tabs[m.screen] = m.paint("87", "["+tabs[m.screen]+"]", true)
	section := strings.Join(tabs, "  ")
	if m.screen == mapScreen && !m.splitCatalog() && (m.searching || m.query != "") {
		section = lipgloss.NewStyle().MaxWidth(width).Render(m.searchLabel())
	} else if lipgloss.Width(section) > width {
		section = tabs[m.screen]
	}

	body := m.contentLines()
	n := m.bodyHeight()
	offset := m.scroll
	if (m.screen == mapScreen && !m.splitCatalog()) || m.screen == bookmarksScreen {
		for index, line := range body {
			if strings.Contains(line, "▸") {
				offset = max(0, index-n+2)
				break
			}
		}
	}
	start := max(0, min(offset, len(body)-n))
	end := min(start+n, len(body))
	visible := append([]string(nil), body[start:end]...)
	for len(visible) < n {
		visible = append(visible, "")
	}

	var footer string
	switch m.screen {
	case mapScreen:
		footer = "↑↓ Select   Enter Demo   l Read   g Resume   B Bookmarks   / Search   ? Help"
	case demoScreen:
		footer = "←→ Steps   l Read   t Task   ? Commands   m Map"
		if len(m.lesson.Demo[m.frame].Permissions) == 9 {
			footer = "x Toggle   " + footer
		}
	case taskScreen, resultsScreen:
		footer = "Enter Bash   c Check   h Hint   ? Commands   m Map"
	case readScreen:
		footer = "←→ Sections   ↑↓ Scroll   b Bookmark   B Saved   Enter Demo   ? Help"
	case bookmarksScreen:
		footer = "↑↓ Select   Enter Open   m Map   ? Help"
	case commandsScreen:
		footer = "? / Esc Back   ↑↓ Scroll   m Map"
	}
	if width < 55 {
		switch m.screen {
		case mapScreen:
			footer = "↑↓ Enter l g B / ?"
		case demoScreen:
			footer = "←→  l t x  ? m"
		case taskScreen, resultsScreen:
			footer = "Enter  c h r  ? m"
		case readScreen:
			footer = "←→ ↑↓ b B Enter ? m"
		case commandsScreen:
			footer = "↑↓  Esc  m"
		}
	} else if len(body) > n && m.screen != mapScreen {
		footer += fmt.Sprintf("   ↑↓ %d/%d", end, len(body))
	}
	if m.searching {
		footer = "Type to filter   ↑↓ Select   Enter Open   Esc Clear"
	}
	footer = lipgloss.NewStyle().MaxWidth(width).Render(footer)
	return lipgloss.NewStyle().PaddingLeft(1).Render(strings.Join([]string{
		lipgloss.NewStyle().MaxWidth(width).Render(header), section, "", strings.Join(visible, "\n"), "",
		m.paint("87", footer, true),
	}, "\n"))
}

func (m model) bodyHeight() int { return max(1, m.height-5) }

func (m model) contentLines() []string {
	width := max(12, min(m.width-2, 100))
	var body string
	switch m.screen {
	case mapScreen:
		if m.splitCatalog() {
			body = m.catalogView(width)
		} else {
			body = m.mapView()
		}
	case demoScreen:
		body = m.demoView(width)
	case taskScreen:
		body = m.taskView(width)
	case resultsScreen:
		body = m.resultsView(width)
	case readScreen:
		body = m.readerView()
	case bookmarksScreen:
		body = m.bookmarksView()
	case commandsScreen:
		body = m.commandsView()
	}
	if m.status != "" {
		body = m.paint("221", m.status, true) + "\n\n" + body
	}
	return strings.Split(lipgloss.NewStyle().Width(width).Render(body), "\n")
}

func (m model) mapView() string {
	complete := 0
	for _, score := range m.scores {
		if score.Best == 100 {
			complete++
		}
	}
	body := m.paint("213", fmt.Sprintf("COURSE PROGRESS  %d/%d complete", complete, len(m.entries)), true)
	workshop := ""
	matches := m.matches()
	if len(matches) == 0 {
		return body + "\n\nNo matching lessons."
	}
	for _, i := range matches {
		entry := m.entries[i]
		if entry.Workshop != workshop {
			workshop = entry.Workshop
			body += "\n\n" + m.paint("251", m.workshopLabel(workshop), true)
		}
		label := "   " + entry.Lesson.Title
		if m.scores[i].Checks > 0 {
			label += fmt.Sprintf("  %d/100", m.scores[i].Best)
		}
		if i == m.menu {
			label = m.paint("87", "▸  "+label[3:], true)
		}
		body += "\n" + label
	}
	return body
}

func (m model) workshopLabel(workshop string) string {
	total, complete := 0, 0
	for i, entry := range m.entries {
		if entry.Workshop == workshop {
			total++
			if m.scores[i].Best == 100 {
				complete++
			}
		}
	}
	return fmt.Sprintf("%s  %d/%d", strings.ToUpper(workshop), complete, total)
}

func (m model) commandsView() string {
	id := m.lesson.ID
	return strings.Join([]string{
		m.paint("213", "CONTROLS AND COMMANDS", true),
		"", "Map: ↑↓ Select · Enter Demo · l Read · t Task · s Score",
		"Demo: ←→ Steps · x Toggle permission preview",
		"Read: ←→ Sections · ↑↓ Scroll · b Bookmark · Enter Demo",
		"Map: g Resume reading · B Bookmarks",
		"Task: Enter Bash · c Check · h Hint · r Replay",
		"Anywhere: m Map · / Search · ? Commands · q Quit",
		"", "From any terminal:",
		"  cli-silly list",
		"  cli-silly resume",
		"  cli-silly read " + id,
		"  cli-silly open " + id,
		"  cli-silly check " + id,
		"  cli-silly progress",
		"", "Inside the practice shell:",
		"  cli-silly task", "  cli-silly hint 2", "  cli-silly check", "  exit",
	}, "\n")
}

func (m model) demoView(width int) string {
	frame := m.lesson.Demo[m.frame]
	var steps []string
	for i := range m.lesson.Demo {
		if i == m.frame {
			steps = append(steps, m.paint("213", "◆", true))
		} else {
			steps = append(steps, "◇")
		}
	}
	title := m.paint("213", fmt.Sprintf("WORKED EXAMPLE  %d / %d", m.frame+1, len(steps)), true) + "   " + strings.Join(steps, " ")
	graphic := m.paint("87", frame.Visual, false)
	if len(frame.Permissions) == 9 {
		owner := frame.Permissions[:3]
		if m.toggleExec {
			if owner[2] == 'x' {
				owner = owner[:2] + "-"
			} else {
				owner = owner[:2] + "x"
			}
		}
		accent := "87"
		if m.flash%2 == 1 {
			accent = "221"
		}
		graphic = m.paint("87", strings.ReplaceAll(frame.Visual, "   -"+frame.Permissions, ""), false)
		if width >= 42 {
			graphic += "\n\nOWNER    GROUP    OTHER\n" + m.paint(accent, owner, true) + "      " + frame.Permissions[3:6] + "      " + frame.Permissions[6:9]
		} else {
			graphic += "\n\nOwner " + m.paint(accent, owner, true) + " · Group " + frame.Permissions[3:6] + " · Other " + frame.Permissions[6:9]
		}
		state := "Direct run as owner: permission denied"
		if owner[2] == 'x' {
			state = "Direct run as owner: execute bit set"
		}
		graphic += "\n" + state
		if m.toggleExec {
			graphic += " (preview)"
		}
	} else if frame.OutputFile != "" {
		accent := "87"
		if m.flash%2 == 1 {
			accent = "221"
		}
		graphic += "\n\nstdout " + m.paint(accent, "━━━━▶  "+frame.OutputFile, true)
	}
	return title + "\n" + m.panel(width, m.paint("221", "$ "+frame.Command, true)+"\n\n"+graphic) +
		"\n" + frame.Explanation
}

func (m model) taskView(width int) string {
	body := m.panel(width, m.paint("213", "YOUR TASK", true)+"\n"+m.lesson.Task)
	for _, check := range m.lesson.Checks {
		body += "\n" + m.paint("87", fmt.Sprintf("%3d pts  ", check.Points), true) + check.Label
	}
	if m.hints > 0 {
		body += "\n\n" + m.paint("221", fmt.Sprintf("Hint %d/%d", m.hints, len(m.lesson.Hints)), true) + "\n" + m.lesson.Hints[m.hints-1]
	}
	if len(m.lesson.Tools) > 0 {
		body += "\n\nTools: " + strings.Join(m.lesson.Tools, ", ")
	}
	return body + "\n\nEnter opens Bash. Run cli-silly task, cli-silly hint, and cli-silly check. Type exit to return."
}

func (m model) resultsView(width int) string {
	if m.progress.Checks == 0 {
		return m.panel(width, "No score yet\n\nOpen Bash and run cli-silly check after attempting the task.")
	}
	result := m.progress.Latest
	title := m.paint("221", "KEEP GOING", true)
	if result.Complete {
		title = m.paint("120", "✦  LAB COMPLETE", true)
	}
	body := m.panel(width, title+"\n"+m.score(result.Score)+fmt.Sprintf("  ·  Best %d/100", m.progress.Best))
	for _, item := range result.Items {
		badge := m.paint("120", "✓ PASS", true)
		if item.Earned != item.Possible {
			badge = m.paint("221", "○  TRY", true)
		}
		body += fmt.Sprintf("\n%s  %d/%d  %s", badge, item.Earned, item.Possible, item.Label)
		if item.Feedback != "" {
			body += "\n       " + item.Feedback
		}
	}
	return body
}

func (m model) panel(width int, content string) string {
	if width < 27 {
		return content
	}
	style := lipgloss.NewStyle().Border(lipgloss.RoundedBorder()).Width(width-2).Padding(0, 1)
	if m.useColor {
		style = style.BorderForeground(lipgloss.Color("99"))
	}
	return style.Render(content)
}

func (m model) score(score int) string {
	filled := score / 10
	return m.paint("87", strings.Repeat("█", filled), true) + strings.Repeat("░", 10-filled) + fmt.Sprintf("  %d/100", score)
}

func (m model) paint(color, text string, bold bool) string {
	if !m.useColor {
		return text
	}
	return lipgloss.NewStyle().Foreground(lipgloss.Color(color)).Bold(bold).Render(text)
}
