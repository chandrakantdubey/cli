package tui

import (
	"context"
	"errors"
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
	filesScreen
	quizScreen
	reviewScreen
	referenceScreen
	referenceDetailScreen
)

type model struct {
	lesson          course.Lesson
	entries         []course.Entry
	scores          []lab.Progress
	store           lab.Store
	helpFrom        viewPosition
	filesFrom       viewPosition
	bookmarksFrom   viewPosition
	busy            bool
	progress        lab.Progress
	openShell       func(course.Lesson) (*exec.Cmd, error)
	screen          screen
	menu            int
	frame           int
	hints           int
	scroll          int
	width           int
	height          int
	flash           int
	frameToken      int
	toggleExec      bool
	useColor        bool
	status          string
	search          *catalogSearch
	query           string
	reading         lab.ReadingState
	sections        []section
	sectionIndex    int
	readerLines     []string
	readerStarts    []int
	readerScroll    int
	readerLesson    string
	readerWidth     int
	readingVisits   map[string]readingVisit
	outline         bool
	outlineIndex    int
	bookmarkIndex   int
	files           []string
	fileIndex       int
	watching        bool
	watchGeneration int
	watchContext    context.Context
	watchCancel     context.CancelFunc
	watchAttempt    os.FileInfo
	quizState       lab.ReviewState
	quizIndex       int
	quizChoice      int
	quizAnswered    bool
	quizFinished    bool
	reviewTasks     []reviewTask
	reviewIndex     int
	references      []course.Reference
	refQuery        string
	refIndex        int
	refSearching    bool
	activeReference course.Reference
	printReference  string
}

type shellExited struct {
	err     error
	attempt os.FileInfo
}
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
	case "do":
		m.screen = taskScreen
	case "resume":
		anchor := ""
		if m.reading.Last.LessonID == lessonID {
			anchor = m.reading.Last.Anchor
		}
		m.openReading(anchor)
	}
	final, err := tea.NewProgram(m, tea.WithAltScreen(), tea.WithInput(os.Stdin), tea.WithOutput(os.Stdout)).Run()
	if err == nil {
		if output := final.(model).printReference; output != "" {
			fmt.Print(output)
		}
	}
	return err
}

func (m *model) selectLesson(index int) {
	m.rememberReading()
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
		if m.screen == readScreen || (m.screen == commandsScreen && m.helpFrom.screen == readScreen) {
			m.renderReading()
		}
	case editorExited:
		m.restoreView(m.filesFrom)
		if message.err != nil {
			m.status = "Editor: " + message.err.Error()
		}
	case watchTick:
		if m.watching && message.generation == m.watchGeneration {
			return m, func() tea.Msg {
				_, err := m.store.WatchCheck(m.watchContext, m.lesson, m.watchAttempt)
				return watched{message.generation, err}
			}
		}
	case watched:
		if m.watching && message.generation == m.watchGeneration {
			if message.err != nil {
				m.stopWatch()
				m.status = message.err.Error()
				return m, nil
			}
			if progress, err := m.store.Progress(m.lesson); err == nil {
				m.progress, m.scores[m.menu] = progress, progress
			} else {
				m.stopWatch()
				m.status = err.Error()
				return m, nil
			}
			return m, nextWatch(m.watchGeneration)
		}
	case shellExited:
		var exit *exec.ExitError
		if message.attempt == nil || (message.err != nil && !errors.As(message.err, &exit)) {
			return m.showResults(message.err), nil
		}
		m.busy, m.status = true, "Checking lab files…"
		return m, func() tea.Msg {
			_, err := m.store.WatchCheck(context.Background(), m.lesson, message.attempt)
			return graded{err}
		}
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
		if m.refSearching && key != "ctrl+c" {
			return m.referenceSearch(message), nil
		}
		if m.search != nil && key != "ctrl+c" {
			return m.searchKey(message), nil
		}
		if m.busy && key != "ctrl+c" {
			return m, nil
		}
		switch key {
		case "m", "B", "/", "?", "esc", "q", "ctrl+c":
			m.stopWatch()
		}
		switch key {
		case "q", "ctrl+c":
			return m, tea.Quit
		case "n":
			return m.openReferences()
		case "z":
			return m.openQuiz()
		case "V":
			m.stopWatch()
			m.screen, m.reviewIndex, m.scroll, m.status = reviewScreen, 0, 0, ""
			if state, err := m.store.Review(); err == nil {
				m.quizState = state
			} else {
				m.status = err.Error()
			}
			m.reviewTasks = m.dueReviews()
		case "m":
			m.screen, m.scroll, m.status = mapScreen, 0, ""
		case "B":
			if m.screen == bookmarksScreen {
				m.restoreView(m.bookmarksFrom)
				return m, nil
			}
			m.bookmarksFrom = m.viewPosition()
			m.screen, m.scroll, m.bookmarkIndex, m.status = bookmarksScreen, 0, 0, ""
		case "/":
			if m.screen == referenceScreen {
				m.refSearching = true
				return m, nil
			}
			m.search = &catalogSearch{from: m.viewPosition(), index: m.menu, query: m.query, status: m.status}
			m.screen, m.scroll, m.status = mapScreen, 0, ""
		case "?":
			if m.screen != commandsScreen {
				m.helpFrom = m.viewPosition()
				m.screen, m.scroll = commandsScreen, 0
			} else {
				m.restoreView(m.helpFrom)
			}
		case "esc":
			m.back()
		case "tab", "shift+tab":
			m.switchMode()
		case "[", "]":
			if m.inLesson() {
				direction := 1
				if key == "[" {
					direction = -1
				}
				m.browseLesson(direction)
			}
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
	switch m.screen {
	case demoScreen, taskScreen, resultsScreen, commandsScreen, referenceDetailScreen:
		last := func() int { return max(0, len(m.contentLines())-m.bodyHeight()) }
		switch key {
		case "pgup":
			m.scroll = max(0, m.scroll-m.bodyHeight())
			return m, nil
		case "pgdown":
			m.scroll = min(last(), m.scroll+m.bodyHeight())
			return m, nil
		case "home":
			m.scroll = 0
			return m, nil
		case "end":
			m.scroll = last()
			return m, nil
		}
	}
	if key == "up" || key == "down" {
		direction := 1
		if key == "up" {
			direction = -1
		}
		switch m.screen {
		case mapScreen:
			m.moveSelection(direction)
		case readScreen:
			if m.outline {
				m.outlineIndex = max(0, min(m.outlineIndex+direction, len(m.sections)-1))
			} else {
				m.scrollReading(direction)
			}
		case referenceScreen:
			m.refIndex = max(0, min(m.refIndex+direction, len(m.referenceMatches())-1))
		case reviewScreen:
			m.reviewIndex = max(0, min(m.reviewIndex+direction, len(m.reviewTasks)-1))
		case quizScreen:
			if m.quizAnswered || m.quizFinished {
				m.scroll = max(0, min(m.scroll+direction, len(m.contentLines())-m.bodyHeight()))
			} else {
				m.quizChoice = max(0, min(m.quizChoice+direction, len(m.lesson.Quiz[m.quizIndex].Choices)-1))
			}
		case filesScreen:
			m.fileIndex = max(0, min(m.fileIndex+direction, len(m.files)-1))
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
		case "enter", " ", "r", "l":
			m.resumeReading()
			return m, nil
		case "d", "t":
			m.screen = taskScreen
		case "s":
			m.screen = resultsScreen
		case "g":
			if m.reading.Last.LessonID == "" {
				m.openReading(overview)
			} else {
				m.openPosition(m.reading.Last)
			}
			return m, nil
		case "left":
			m.moveWorkshop(-1)
		case "right":
			m.moveWorkshop(1)
		case "home":
			m.selectLesson(m.matches()[0])
		case "end":
			matches := m.matches()
			m.selectLesson(matches[len(matches)-1])
		}
		m.scroll, m.status = 0, ""
	case readScreen:
		switch key {
		case "i":
			m.outline, m.outlineIndex = !m.outline, m.sectionIndex
		case "enter", "e":
			if m.outline {
				m.outline = false
				m.jumpSection(m.outlineIndex)
				return m, nil
			}
			m.screen, m.frame, m.scroll, m.flash, m.toggleExec = demoScreen, 0, 0, 0, false
		case "d", "t":
			m.screen, m.scroll = taskScreen, 0
		case "left":
			m.jumpSection(max(0, m.sectionIndex-1))
		case "right":
			m.jumpSection(min(m.sectionIndex+1, len(m.sections)-1))
		case "b":
			m.updateReading(true)
		case "pgup":
			m.scrollReading(-m.readingHeight())
		case "pgdown":
			m.scrollReading(m.readingHeight())
		case "home":
			m.jumpSection(0)
		case "end":
			m.scrollReading(len(m.readerLines))
		}
	case bookmarksScreen:
		if key == "enter" && len(m.reading.Bookmarks) > 0 {
			m.openPosition(m.reading.Bookmarks[m.bookmarkIndex])
		}
	case demoScreen:
		switch key {
		case "r", "l":
			m.resumeReading()
		case "left", "b":
			m.frame, m.scroll, m.flash, m.toggleExec = max(0, m.frame-1), 0, 0, false
		case "x":
			if len(m.lesson.Demo[m.frame].Permissions) == 9 {
				m.toggleExec = !m.toggleExec
			}
		case "d", "t":
			m.screen, m.scroll = taskScreen, 0
		case "right", "enter", " ":
			m.scroll = 0
			if m.frame == len(m.lesson.Demo)-1 {
				if key != "right" {
					m.screen = taskScreen
				}
			} else {
				m.frame++
				m.toggleExec = false
				if m.useColor {
					m.flash, m.frameToken = 3, m.frameToken+1
					return m, animate(m.frameToken)
				}
			}
		}
	case referenceScreen:
		matches := m.referenceMatches()
		if key == "enter" && len(matches) > 0 {
			m.activeReference = matches[m.refIndex]
			m.screen, m.scroll = referenceDetailScreen, 0
		}
	case referenceDetailScreen:
		if key == "p" {
			body, err := m.activeReference.Format(nil)
			if err != nil {
				m.status = err.Error()
				return m, nil
			}
			m.printReference = body
			return m, tea.Quit
		}
	case quizScreen:
		return m.quizKey(key)
	case reviewScreen:
		if key == "enter" && len(m.reviewTasks) > 0 {
			m.selectLesson(m.reviewTasks[m.reviewIndex].entry)
			m.query, m.screen = "", taskScreen
		}
	case filesScreen:
		if key == "enter" && len(m.files) > 0 {
			cmd, err := m.store.Editor(m.lesson, m.files[m.fileIndex])
			if err != nil {
				m.status = err.Error()
				return m, nil
			}
			m.status = ""
			return m, tea.ExecProcess(cmd, func(err error) tea.Msg { return editorExited{err} })
		}
	case taskScreen, resultsScreen:
		switch key {
		case "e":
			return m.openFiles()
		case "w":
			return m.toggleWatch()
		case "h":
			m.stopWatch()
			m.screen, m.scroll, m.status = taskScreen, 0, ""
			m.hints = m.hints%len(m.lesson.Hints) + 1
			lines := m.contentLines()
			label := fmt.Sprintf("Hint %d/%d", m.hints, len(m.lesson.Hints))
			for index, line := range lines {
				if strings.Contains(line, label) {
					m.scroll = min(index, max(0, len(lines)-m.bodyHeight()))
					break
				}
			}
		case "c":
			m.stopWatch()
			m.busy = true
			m.status = "Checking workspace…"
			return m, func() tea.Msg {
				if _, err := m.store.Open(m.lesson); err != nil {
					return graded{err}
				}
				_, err := m.store.Check(m.lesson, "")
				return graded{err}
			}
		case "d", "t":
			m.stopWatch()
			m.screen, m.scroll, m.status = taskScreen, 0, ""
		case "r", "l":
			m.stopWatch()
			m.resumeReading()
		case "enter", "o":
			m.stopWatch()
			cmd, err := m.openShell(m.lesson)
			if err != nil {
				m.status = err.Error()
				return m, nil
			}
			m.status = ""
			var attempt os.FileInfo
			if lab.Watchable(m.lesson) {
				attempt, err = os.Stat(cmd.Dir)
				if err != nil {
					m.status = err.Error()
					return m, nil
				}
			}
			return m, tea.ExecProcess(cmd, func(err error) tea.Msg { return shellExited{err, attempt} })
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
	width := m.layoutWidth()
	header := m.paint(violet, "◈  CLI-SILLY", true)
	if width >= 32 {
		title := "LESSONS"
		if m.screen != mapScreen && len(m.entries) > 0 {
			title = m.entries[m.menu].Workshop
		}
		header += "  " + m.paint(muted, "›  "+title, false)
		if m.inLesson() && len(m.entries) > 0 {
			position, total := m.workshopPosition()
			header += m.paint(muted, fmt.Sprintf("  ·  Lesson %d/%d", position, total), false)
		}
	}
	section := "Read  ·  Do"
	switch m.screen {
	case readScreen, demoScreen:
		section = m.activeTab("Read") + "  ·  Do"
	case taskScreen, resultsScreen, filesScreen:
		section = "Read  ·  " + m.activeTab("Do")
	}
	if m.inLesson() && width >= 33 {
		section += "   " + m.paint(muted, "Tab switch", false)
	}
	if m.screen == mapScreen && !m.splitCatalog() && (m.search != nil || m.query != "") {
		section = lipgloss.NewStyle().MaxWidth(width).Render(m.searchLabel())
	}

	body := m.contentLines()
	n := m.bodyHeight()
	offset := m.scroll
	if (m.screen == mapScreen && !m.splitCatalog()) || m.screen == bookmarksScreen || m.screen == filesScreen || m.screen == reviewScreen || m.screen == referenceScreen {
		for index, line := range body {
			if strings.Contains(line, "▸") {
				offset = max(0, index-n+2)
				break
			}
		}
	}
	start := max(0, min(offset, len(body)-n))
	end := min(start+n, len(body))
	if width >= 48 && (m.screen == taskScreen || m.screen == resultsScreen || m.screen == demoScreen) && len(body) > n {
		section += "   " + m.paint(muted, fmt.Sprintf("↑↓ %d/%d lines", end, len(body)), false)
	}
	visible := append([]string(nil), body[start:end]...)
	for len(visible) < n {
		visible = append(visible, "")
	}

	footer := m.footer(width)
	return lipgloss.NewStyle().PaddingLeft(1).Render(strings.Join([]string{
		lipgloss.NewStyle().MaxWidth(width).Render(header), section, "", strings.Join(visible, "\n"), "",
		footer,
	}, "\n"))
}

func (m model) bodyHeight() int { return max(1, m.height-5) }

func (m model) layoutWidth() int { return max(12, min(m.width-2, 118)) }

func (m model) contentLines() []string {
	width := m.layoutWidth()
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
		body = m.readerView(width)
	case referenceScreen, referenceDetailScreen:
		body = m.referenceView()
	case quizScreen:
		body = m.quizView()
	case reviewScreen:
		body = m.reviewView()
	case filesScreen:
		body = m.filesView()
	case bookmarksScreen:
		body = m.bookmarksView()
	case commandsScreen:
		body = m.commandsView()
	}
	if m.status != "" {
		body = m.paint(amber, m.status, true) + "\n\n" + body
	}
	return strings.Split(lipgloss.NewStyle().Width(width).Render(body), "\n")
}

func (m model) mapView() string {
	complete, total := m.requiredProgress()
	body := m.paint(violet, fmt.Sprintf("COURSE PROGRESS  %d/%d complete", complete, total), true)
	workshop := ""
	matches := m.matches()
	if len(matches) == 0 {
		return body + "\n\nNo matching lessons."
	}
	for _, i := range matches {
		entry := m.entries[i]
		if entry.Workshop != workshop {
			workshop = entry.Workshop
			body += "\n\n" + m.paint(muted, m.workshopLabel(workshop), true)
		}
		body += "\n" + m.lessonRow(i, m.layoutWidth())
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

func (m model) workshopPosition() (int, int) {
	current := m.entries[m.menu]
	position, total := 0, 0
	for index, entry := range m.entries {
		if entry.PackID == current.PackID && entry.Workshop == current.Workshop {
			total++
			if index <= m.menu {
				position++
			}
		}
	}
	return position, total
}

func (m model) commandsView() string {
	id := m.lesson.ID
	return strings.Join([]string{
		m.paint(violet, "CONTROLS AND COMMANDS", true),
		"", "Lessons: ↑↓ Select · ←→ Workshops · Enter/r Read · d Do · s Score",
		"Read: ↑↓ Scroll · i Contents · ←→ Sections · e Example · b Bookmark",
		"Lesson views: Tab switches Read and Do · [ / ] browse lessons · Esc Back",
		"Example: ←→ Steps · x Toggle permission preview · r Read · d Do",
		"Map: g Resume reading · B Bookmarks · z Quiz · V Review",
		"Do: Enter Bash · e Edit · c Check · w Watch · h Hint · r Read",
		"Anywhere: n References · m Lessons · / Search · ? Commands · q Quit",
		"", "From any terminal:",
		"  cli-silly list",
		"  cli-silly resume",
		"  cli-silly read " + id,
		"  cli-silly do " + id,
		"  cli-silly check " + id,
		"  cli-silly progress",
		"", "Inside the practice shell:",
		"  cli-silly read", "  cli-silly task", "  cli-silly hint 2", "  cli-silly check", "  exit",
	}, "\n")
}

func (m model) demoView(width int) string {
	frame := m.lesson.Demo[m.frame]
	var steps []string
	for i := range m.lesson.Demo {
		switch {
		case i == m.frame:
			steps = append(steps, m.paint(violet, "◆", true))
		case i < m.frame:
			steps = append(steps, m.paint(teal, "✓", false))
		default:
			steps = append(steps, "◇")
		}
	}
	title := m.paint(violet, fmt.Sprintf("WORKED EXAMPLE  %d / %d", m.frame+1, len(steps)), true) + "   " + strings.Join(steps, " ")
	graphic := m.paint(teal, frame.Visual, false)
	if len(frame.Permissions) == 9 {
		owner := frame.Permissions[:3]
		if m.toggleExec {
			if owner[2] == 'x' {
				owner = owner[:2] + "-"
			} else {
				owner = owner[:2] + "x"
			}
		}
		accent := teal
		if m.flash%2 == 1 {
			accent = amber
		}
		graphic = m.paint(teal, strings.ReplaceAll(frame.Visual, "   -"+frame.Permissions, ""), false)
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
		accent := teal
		if m.flash%2 == 1 {
			accent = amber
		}
		graphic += "\n\nstdout " + m.paint(accent, "━━━━▶  "+frame.OutputFile, true)
	}
	main := m.paint(muted, "Command", true) + "\n\n" + m.paint(amber, "$ "+frame.Command, true)
	observation := m.paint(muted, "What you see", true) + "\n\n" + graphic
	panels := m.panel(width, main+"\n\n"+observation)
	if width >= 86 {
		panels = m.lessonPanels(width, main, observation)
	}
	return title + "\n\n" + m.lesson.Title + "\n\n" + panels +
		"\n\n" + m.paint(violet, "Why it works", true) + "\n" + frame.Explanation
}

func (m model) taskView(width int) string {
	mainWidth, _ := lessonPanelWidths(width)
	body := m.paint(violet, "YOUR TASK", true) + "\n" + m.lesson.Title + "\n\n" + wrapTask(m.lesson.Task, mainWidth-4)
	body += "\n\n" + m.paint(muted, "What earns points", true)
	for _, check := range m.lesson.Checks {
		body += "\n" + m.paint(teal, fmt.Sprintf("%3d pts  ", check.Points), true) + check.Label
	}
	if m.hints > 0 {
		body += "\n\n" + m.paint(amber, fmt.Sprintf("Hint %d/%d", m.hints, len(m.lesson.Hints)), true) + "\n" + m.lesson.Hints[m.hints-1]
	}
	aside := m.paint(violet, "PRACTICE WORKSPACE", true) + "\n" + m.paint(teal, "Enter Open Bash", true)
	if files, err := lab.StarterFiles(m.lesson); err == nil && len(files) > 0 {
		aside += "\n\n" + m.paint(muted, "Starter files", true) + "\n  " + strings.Join(files[:min(4, len(files))], "\n  ")
		if len(files) > 4 {
			aside += fmt.Sprintf("\n  and %d more", len(files)-4)
		}
		aside += "\n" + m.paint(teal, "e Browse and edit files", false)
	}
	if len(m.lesson.Tools) > 0 {
		aside += "\n\nTools: Bash, " + strings.Join(m.lesson.Tools, ", ")
	}
	return m.lessonPanels(width, body, aside+m.practiceCommands()) + "\n\n" +
		m.paint(muted, "Workspace  ", false) + m.store.Path(m.lesson)
}

func (m model) resultsView(width int) string {
	watch := ""
	if m.watching {
		watch = m.paint(teal, "Watching file checks · w Pause", true) + "\n"
	}
	if m.progress.Checks == 0 {
		return watch + m.panel(width, m.lesson.Title+"\n\nNo score yet\nPress Esc to review the task, or press Enter to open Bash.")
	}
	result := m.progress.Latest
	title := m.paint(amber, "KEEP GOING", true)
	if result.Complete {
		title = m.paint(green, "✦  LAB COMPLETE", true)
	}
	body := m.lesson.Title + "\n" + title + "\n" + m.score(result.Score) + fmt.Sprintf("  ·  Best %d/100\n", m.progress.Best)
	for _, item := range result.Items {
		badge := m.paint(green, "✓ PASS", true)
		if item.Earned != item.Possible {
			badge = m.paint(amber, "○  TRY", true)
		}
		body += fmt.Sprintf("\n%s  %d/%d  %s", badge, item.Earned, item.Possible, item.Label)
		if item.Feedback != "" {
			body += "\n       " + item.Feedback
		}
	}
	aside := m.paint(violet, "NEXT ACTION", true)
	if next := m.adjacentLesson(1); result.Complete && next >= 0 {
		aside += "\n" + m.paint(teal, "] Read next", true) + "\n" + m.entries[next].Lesson.Title
	} else if !result.Complete {
		aside += "\n" + m.paint(teal, "Enter Retry in Bash", true) + "\nh Show the next hint"
	} else {
		aside += "\n" + m.paint(teal, "m Browse lessons", true)
	}
	aside += "\n\nTab Return to reading\nEsc Review the task"
	return watch + m.lessonPanels(width, body, aside)
}

func (m model) lessonPanels(width int, main, aside string) string {
	left, right := lessonPanelWidths(width)
	if right > 0 {
		return lipgloss.JoinHorizontal(lipgloss.Top, m.panel(left, main), "  ", m.panel(right, aside))
	}
	return m.panel(width, main) + "\n\n" + m.panel(width, aside)
}

func lessonPanelWidths(width int) (int, int) {
	if width >= 86 {
		right := width * 37 / 100
		return width - right - 2, right
	}
	return width, 0
}

func (m model) panel(width int, content string) string {
	if width < 27 {
		return content
	}
	style := lipgloss.NewStyle().Border(lipgloss.RoundedBorder()).Width(width-2).Padding(0, 1)
	if m.useColor {
		style = style.BorderForeground(border)
	}
	return style.Render(content)
}

func (m model) score(score int) string {
	filled := score / 10
	return m.paint(teal, strings.Repeat("█", filled), true) + strings.Repeat("░", 10-filled) + fmt.Sprintf("  %d/100", score)
}
