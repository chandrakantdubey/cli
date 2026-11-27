package tui

import (
	"bytes"
	"fmt"
	"slices"
	"strings"
	"unicode"

	"cli-silly/course"
	"cli-silly/internal/lab"

	"github.com/charmbracelet/glamour"
	glamouransi "github.com/charmbracelet/glamour/ansi"
	"github.com/charmbracelet/glamour/styles"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/x/ansi"
	"github.com/muesli/termenv"
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/text"
)

const overview = "overview"

type section struct {
	anchor   string
	title    string
	markdown string
}

type readingVisit struct {
	anchor string
	line   int
}

func lessonSections(body string) []section {
	source := []byte(body)
	doc := goldmark.DefaultParser().Parse(text.NewReader(source))
	sections := []section{{anchor: overview, title: "Overview"}}
	used := map[string]bool{overview: true}
	start := 0
	for node := doc.FirstChild(); node != nil; node = node.NextSibling() {
		heading, ok := node.(*ast.Heading)
		if !ok || heading.Lines().Len() == 0 {
			continue
		}
		pos := heading.Lines().At(0).Start
		line := bytes.LastIndexByte(source[:pos], '\n') + 1
		sections[len(sections)-1].markdown = strings.TrimSpace(body[start:line])
		title := strings.TrimSpace(string(heading.Text(source)))
		slug := strings.Join(strings.FieldsFunc(strings.ToLower(title), func(r rune) bool {
			return !unicode.IsLetter(r) && !unicode.IsDigit(r)
		}), "-")
		if slug == "" {
			slug = "section"
		}
		anchor := slug
		for i := 2; used[anchor]; i++ {
			anchor = fmt.Sprintf("%s-%d", slug, i)
		}
		used[anchor] = true
		sections = append(sections, section{anchor: anchor, title: title})
		start = line
	}
	sections[len(sections)-1].markdown = strings.TrimSpace(body[start:])
	return sections
}

func (m *model) openReading(anchor string) {
	m.sections = lessonSections(m.lesson.Reading())
	m.sectionIndex = max(0, slices.IndexFunc(m.sections, func(s section) bool { return s.anchor == anchor }))
	m.readerStarts, m.readerScroll, m.outline = nil, 0, false
	m.readerLesson, m.readerWidth = m.lesson.ID, 0
	m.screen, m.scroll, m.status = readScreen, 0, ""
	m.updateReading(false)
	m.renderReading()
}

func (m *model) rememberReading() {
	if m.readerLesson != m.lesson.ID || len(m.readerStarts) == 0 {
		return
	}
	if m.readingVisits == nil {
		m.readingVisits = make(map[string]readingVisit)
	}
	m.readingVisits[m.lesson.ID] = readingVisit{
		anchor: m.position().Anchor,
		line:   m.readerScroll - m.readerStarts[m.sectionIndex],
	}
}

func (m *model) resumeReading() {
	if m.readerLesson == m.lesson.ID && len(m.sections) > 0 {
		m.screen, m.scroll, m.status, m.outline = readScreen, 0, "", false
		m.renderReading()
		return
	}
	if saved, ok := m.readingVisits[m.lesson.ID]; ok {
		m.openReading(saved.anchor)
		end := len(m.readerLines)
		if m.sectionIndex+1 < len(m.readerStarts) {
			end = m.readerStarts[m.sectionIndex+1]
		}
		m.readerScroll += min(saved.line, max(0, end-m.readerScroll-1))
		return
	}
	anchor := overview
	if m.reading.Last.LessonID == m.lesson.ID {
		anchor = m.reading.Last.Anchor
	}
	m.openReading(anchor)
}

func (m model) contentsWidth() int {
	if m.width >= 86 && m.height >= 17 {
		return min(34, m.layoutWidth()/3)
	}
	return 0
}

func (m *model) renderReading() {
	if len(m.sections) == 0 {
		return
	}
	width := m.layoutWidth() - m.contentsWidth()
	if m.contentsWidth() > 0 {
		width -= 2
	}
	if width == m.readerWidth {
		return
	}
	lineInSection := 0
	if len(m.readerStarts) > m.sectionIndex {
		lineInSection = m.readerScroll - m.readerStarts[m.sectionIndex]
	}
	profile, style := termenv.Ascii, styles.NoTTYStyleConfig
	if m.useColor {
		profile, style = termenv.ANSI256, styles.LightStyleConfig
		if lipgloss.HasDarkBackground() {
			style = styles.TokyoNightStyleConfig
		}
	}
	// The reader supplies the visual hierarchy without displaying Markdown punctuation.
	for _, heading := range []*glamouransi.StyleBlock{&style.H1, &style.H2, &style.H3, &style.H4, &style.H5, &style.H6} {
		heading.Prefix, heading.Suffix = "", ""
	}
	style.H2.Prefix = "◆ "
	style.Strong.BlockPrefix, style.Strong.BlockSuffix = "", ""
	style.Emph.BlockPrefix, style.Emph.BlockSuffix = "", ""
	margin := uint(0)
	style.Document.Margin = &margin
	renderer, err := glamour.NewTermRenderer(
		glamour.WithColorProfile(profile), glamour.WithStyles(style),
		glamour.WithWordWrap(max(12, width-4)), glamour.WithTableWrap(true),
	)
	m.readerLines = nil
	m.readerStarts = make([]int, len(m.sections))
	if err != nil {
		m.status = err.Error()
		return
	}
	m.readerWidth = width
	for i, section := range m.sections {
		m.readerStarts[i] = len(m.readerLines)
		if section.markdown == "" {
			continue
		}
		rendered, err := renderer.Render(section.markdown)
		if err != nil {
			rendered = section.markdown
		}
		if !m.useColor {
			rendered = ansi.Strip(rendered)
		}
		for _, line := range strings.Split(strings.Trim(rendered, "\n"), "\n") {
			m.readerLines = append(m.readerLines, strings.Split(ansi.Hardwrap(line, width, false), "\n")...)
		}
		m.readerLines = append(m.readerLines, "")
	}
	end := len(m.readerLines)
	if m.sectionIndex+1 < len(m.readerStarts) {
		end = m.readerStarts[m.sectionIndex+1]
	}
	start := m.readerStarts[m.sectionIndex]
	m.readerScroll = start + min(max(0, end-start-1), lineInSection)
}

func (m model) readingHeight() int { return max(1, m.bodyHeight()-3) }

func (m model) position() lab.ReadingPosition {
	return lab.ReadingPosition{LessonID: m.lesson.ID, Anchor: m.sections[m.sectionIndex].anchor}
}

func (m *model) updateReading(toggleBookmark bool) {
	state, err := m.store.SaveReading(m.position(), toggleBookmark)
	if err != nil {
		m.status = "Could not save reading position: " + err.Error()
		return
	}
	m.reading, m.status = state, ""
}

func (m *model) jumpSection(index int) {
	if index < 0 || index >= len(m.readerStarts) {
		return
	}
	m.sectionIndex, m.readerScroll = index, m.readerStarts[index]
	m.updateReading(false)
}

func (m *model) scrollReading(delta int) {
	if len(m.readerStarts) == 0 {
		return
	}
	lastSection := m.readerStarts[len(m.readerStarts)-1]
	end := max(lastSection, len(m.readerLines)-m.readingHeight())
	m.readerScroll = max(0, min(m.readerScroll+delta, end))
	index := m.sectionIndex
	for index > 0 && m.readerStarts[index] > m.readerScroll {
		index--
	}
	for index+1 < len(m.readerStarts) && m.readerStarts[index+1] <= m.readerScroll {
		index++
	}
	if index != m.sectionIndex {
		m.sectionIndex = index
		m.updateReading(false)
	}
}

func (m *model) openPosition(position lab.ReadingPosition) {
	index := slices.IndexFunc(m.entries, func(e course.Entry) bool { return e.Lesson.ID == position.LessonID })
	if index < 0 {
		m.status = "Saved lesson unavailable. Choose a lesson from the map."
		return
	}
	m.query = ""
	m.selectLesson(index)
	m.openReading(position.Anchor)
}

func (m model) readerView(width int) string {
	if len(m.sections) == 0 {
		return ""
	}
	current := m.sections[m.sectionIndex]
	label := fmt.Sprintf("SECTION %d/%d · %s", m.sectionIndex+1, len(m.sections), current.title)
	if slices.Contains(m.reading.Bookmarks, m.position()) {
		label += " · Bookmarked"
	}
	crop := func(line string) string { return ansi.Truncate(line, width, "…") }
	header := crop(m.paint(violet, m.lesson.Title, true)) + "\n" + crop(m.paint(teal, label, false)) + "\n\n"
	start := min(m.readerScroll, len(m.readerLines))
	end := min(start+m.readingHeight(), len(m.readerLines))
	page := strings.Join(m.readerLines[start:end], "\n")
	side := m.contentsWidth()
	if side > 0 {
		index := m.contents(side-4, m.readingHeight()-2)
		page = lipgloss.JoinHorizontal(lipgloss.Top, m.panel(side, index), "  ", page)
	} else if m.outline {
		page = m.contents(width, m.readingHeight())
	}
	return header + page
}

func (m model) contents(width, height int) string {
	selected := m.sectionIndex
	prompt := "i Jump to a section"
	if m.outline {
		selected = m.outlineIndex
		prompt = "↑↓ Choose · Enter Jump"
	}
	crop := func(line string) string { return ansi.Truncate(line, width, "…") }
	lines := []string{crop(m.paint(violet, "Contents", true)), crop(m.paint(muted, prompt, false))}
	visible := max(1, height-len(lines))
	start := max(0, min(selected-visible+1, len(m.sections)-visible))
	for i := start; i < min(start+visible, len(m.sections)); i++ {
		label := fmt.Sprintf("  %d %s", i+1, m.sections[i].title)
		if i == selected {
			label = m.paint(teal, "▸"+label[1:], true)
		}
		lines = append(lines, crop(label))
	}
	return strings.Join(lines, "\n")
}

func (m model) bookmarksView() string {
	lines := []string{m.paint(violet, "BOOKMARKS", true), ""}
	if len(m.reading.Bookmarks) == 0 {
		return strings.Join(append(lines, "Press b while reading to save a section."), "\n")
	}
	for i, saved := range m.reading.Bookmarks {
		label := saved.LessonID + " · " + saved.Anchor
		for _, entry := range m.entries {
			if saved.LessonID == entry.Lesson.ID {
				label = entry.Lesson.Title
				for _, section := range lessonSections(entry.Lesson.Reading()) {
					if section.anchor == saved.Anchor {
						label += " · " + section.title
						break
					}
				}
				break
			}
		}
		if i == m.bookmarkIndex {
			label = m.paint(teal, "▸ "+label, true)
		} else {
			label = "  " + label
		}
		lines = append(lines, label)
	}
	return strings.Join(lines, "\n")
}
