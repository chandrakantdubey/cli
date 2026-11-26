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

func renderMarkdown(source string, width int, color bool) (string, error) {
	profile, style := termenv.Ascii, "notty"
	if color {
		profile, style = termenv.ANSI256, "dark"
	}
	renderer, err := glamour.NewTermRenderer(
		glamour.WithColorProfile(profile), glamour.WithStandardStyle(style),
		glamour.WithWordWrap(width), glamour.WithTableWrap(true),
	)
	if err != nil {
		return "", err
	}
	result, err := renderer.Render(source)
	if !color {
		result = ansi.Strip(result)
	}
	return strings.Trim(result, "\n"), err
}

func (m *model) openReading(anchor string) {
	m.sections = lessonSections(m.lesson.Body)
	m.sectionIndex = max(0, slices.IndexFunc(m.sections, func(s section) bool { return s.anchor == anchor }))
	m.screen, m.scroll, m.status = readScreen, 0, ""
	m.updateReading(false)
	m.renderSection()
}

func (m *model) renderSection() {
	if len(m.sections) == 0 {
		return
	}
	m.readerText = m.sections[m.sectionIndex].markdown
	if m.readerText == "" && len(m.sections) > 1 {
		m.readerText = "Sections\n\n"
		for _, next := range m.sections[1:] {
			m.readerText += "- " + next.title + "\n"
		}
	}
	rendered, err := renderMarkdown(m.readerText, max(12, min(m.width-2, 100)), m.useColor)
	if err != nil {
		m.status = err.Error()
		return
	}
	m.readerText = rendered
}

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

func (m *model) moveSection(delta int) {
	index := min(len(m.sections)-1, max(0, m.sectionIndex+delta))
	if index == m.sectionIndex {
		return
	}
	m.sectionIndex, m.scroll = index, 0
	m.updateReading(false)
	m.renderSection()
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

func (m model) readerView() string {
	if len(m.sections) == 0 {
		return ""
	}
	current := m.sections[m.sectionIndex]
	label := fmt.Sprintf("SECTION %d/%d · %s", m.sectionIndex+1, len(m.sections), current.title)
	if slices.Contains(m.reading.Bookmarks, m.position()) {
		label += " · Bookmarked"
	}
	return m.paint("213", m.lesson.Title, true) + "\n" + m.paint("87", label, false) + "\n\n" + m.readerText
}

func (m model) bookmarksView() string {
	lines := []string{m.paint("213", "BOOKMARKS", true), ""}
	if len(m.reading.Bookmarks) == 0 {
		return strings.Join(append(lines, "Press b while reading to save a section."), "\n")
	}
	for i, saved := range m.reading.Bookmarks {
		label := saved.LessonID + " · " + saved.Anchor
		for _, entry := range m.entries {
			if saved.LessonID == entry.Lesson.ID {
				label = entry.Lesson.Title
				for _, section := range lessonSections(entry.Lesson.Body) {
					if section.anchor == saved.Anchor {
						label += " · " + section.title
						break
					}
				}
				break
			}
		}
		if i == m.bookmarkIndex {
			label = m.paint("87", "▸ "+label, true)
		} else {
			label = "  " + label
		}
		lines = append(lines, label)
	}
	return strings.Join(lines, "\n")
}
