package tui

import (
	"fmt"
	"slices"
	"strconv"
	"strings"
	"time"

	"cli-silly/course"
	"cli-silly/internal/lab"

	tea "github.com/charmbracelet/bubbletea"
)

type reviewTask struct {
	entry  int
	title  string
	due    time.Time
	missed int
}

func (m model) requiredProgress() (int, int) {
	complete, total := 0, 0
	for i, entry := range m.entries {
		if entry.Lesson.ReviewOf == "" && entry.PackID == "" {
			total++
			if m.scores[i].Best == 100 {
				complete++
			}
		}
	}
	return complete, total
}

func (m model) dueReviews() []reviewTask {
	var tasks []reviewTask
	for i, variant := range m.entries {
		if variant.Lesson.ReviewOf == "" {
			continue
		}
		base := slices.IndexFunc(m.entries, func(e course.Entry) bool { return e.Lesson.ID == variant.Lesson.ReviewOf })
		if base >= 0 && m.scores[base].Best == 100 {
			missed := 0
			for _, answer := range m.quizState.Answers[variant.Lesson.ReviewOf] {
				if !answer.Correct {
					missed++
				}
			}
			tasks = append(tasks, reviewTask{i, variant.Lesson.Title, lab.ReviewDue(m.scores[base], m.scores[i]), missed})
		}
	}
	slices.SortStableFunc(tasks, func(a, b reviewTask) int { return a.due.Compare(b.due) })
	return tasks
}

func (m model) reviewView() string {
	lines := []string{m.paint(violet, "REVIEW QUEUE", true), ""}
	if len(m.reviewTasks) == 0 {
		return strings.Join(append(lines, "Complete a lesson with a linked challenge to schedule a review.", "You can also explore optional challenges from the map."), "\n")
	}
	now := time.Now()
	for i, task := range m.reviewTasks {
		due := "Due now"
		if task.due.After(now) {
			due = "Due " + task.due.Local().Format("Jan 2, 15:04")
		}
		label := task.title + " · " + due
		if task.missed > 0 {
			label += fmt.Sprintf(" · %d missed quiz items", task.missed)
		}
		if i == m.reviewIndex {
			label = m.paint(teal, "▸ "+label, true)
		} else {
			label = "  " + label
		}
		lines = append(lines, label)
	}
	return strings.Join(lines, "\n")
}

func (m model) practiceLinks() []reviewTask {
	var tasks []reviewTask
	for index, entry := range m.entries {
		if entry.Lesson.ReviewOf == m.lesson.ID && m.lesson.ID != "" ||
			m.lesson.ReviewOf != "" && entry.Lesson.ID == m.lesson.ReviewOf {
			tasks = append(tasks, reviewTask{entry: index, title: entry.Lesson.Title})
		}
	}
	return tasks
}

func (m model) practiceLabel() string {
	if len(m.practiceLinks()) == 0 {
		return ""
	}
	if m.lesson.ReviewOf != "" {
		return "Original lesson"
	}
	return "More practice"
}

func (m *model) openRelatedPractice() {
	tasks := m.practiceLinks()
	if len(tasks) == 0 {
		m.status = "Browse Review challenges in Lessons for more practice."
		return
	}
	m.stopWatch()
	m.reviewFrom = m.viewPosition()
	m.reviewTasks, m.reviewIndex, m.relatedReview = tasks, 0, true
	m.screen, m.scroll, m.status = reviewScreen, 0, ""
}

func (m model) relatedPracticeView(width int) string {
	title := "MORE PRACTICE"
	if m.lesson.ReviewOf != "" {
		title = "ORIGINAL LESSON"
	}
	lines := []string{m.paint(violet, title, true), m.lesson.Title, ""}
	for i, task := range m.reviewTasks {
		label := "  " + task.title
		if i == m.reviewIndex {
			label = m.paint(teal, "▸ "+task.title, true)
		}
		lines = append(lines, label)
	}
	task := m.reviewTasks[m.reviewIndex]
	lesson := m.entries[task.entry].Lesson
	score := fmt.Sprintf("Best %d/100", m.scores[task.entry].Best)
	if m.scores[task.entry].Checks == 0 {
		score = "Ready to try"
	}
	return m.lessonPanels(width, strings.Join(lines, "\n"), lesson.Title+"\n"+score+"\n\n"+lesson.Task)
}

func (m model) openQuiz() (tea.Model, tea.Cmd) {
	if len(m.lesson.Quiz) == 0 {
		m.status = "No quiz for this lesson. Try the permissions, Git staging, or log-report lesson."
		return m, nil
	}
	m.stopWatch()
	saved, err := m.store.Review()
	if err != nil {
		m.status = err.Error()
		return m, nil
	}
	m.quizState = saved
	m.quizIndex, m.quizChoice, m.scroll = 0, 0, 0
	m.quizAnswered, m.quizFinished = false, false
	m.screen, m.status = quizScreen, ""
	return m, nil
}

func (m model) quizKey(key string) (tea.Model, tea.Cmd) {
	question := m.lesson.Quiz[m.quizIndex]
	if key == "r" {
		return m.openQuiz()
	}
	if m.quizFinished {
		return m, nil
	}
	if digit, err := strconv.Atoi(key); err == nil && !m.quizAnswered && digit > 0 && digit <= len(question.Choices) {
		m.quizChoice = digit - 1
		return m, nil
	}
	if key != "enter" {
		return m, nil
	}
	if m.quizAnswered {
		if m.quizIndex+1 == len(m.lesson.Quiz) {
			m.quizFinished, m.scroll = true, 0
		} else {
			m.quizAnswered, m.quizIndex, m.quizChoice, m.scroll = false, m.quizIndex+1, 0, 0
		}
		return m, nil
	}
	correct := m.quizChoice+1 == question.Answer
	saved, err := m.store.RecordAnswer(m.lesson.ID, question.ID, correct, time.Now())
	if err != nil {
		m.status = "Couldn't save answer: " + err.Error()
		return m, nil
	}
	m.status, m.quizAnswered, m.quizState = "", true, saved
	if len(m.contentLines()) > m.bodyHeight() {
		m.scroll = len(m.contentLines()) - m.bodyHeight()
	}
	return m, nil
}

func (m model) quizView() string {
	lines := []string{m.paint(violet, "SKILL CHECK", true), m.lesson.Title, ""}
	if m.quizFinished {
		correct := 0
		for _, q := range m.lesson.Quiz {
			if m.quizState.Answers[m.lesson.ID][q.ID].Correct {
				correct++
			}
		}
		lines = append(lines, fmt.Sprintf("%d/%d correct", correct, len(m.lesson.Quiz)), "Press r to try again or m to return to the course map.")
		return strings.Join(lines, "\n")
	}
	question := m.lesson.Quiz[m.quizIndex]
	lines = append(lines, fmt.Sprintf("Question %d/%d", m.quizIndex+1, len(m.lesson.Quiz)), question.Prompt, "")
	for i, choice := range question.Choices {
		label := fmt.Sprintf("%d. %s", i+1, choice)
		if i == m.quizChoice {
			label = m.paint(teal, "▸ "+label, true)
		} else {
			label = "  " + label
		}
		lines = append(lines, label)
	}
	if m.quizAnswered {
		feedback := "Correct."
		if m.quizChoice+1 != question.Answer {
			feedback = "Try this idea again."
		}
		lines = append(lines, "", m.paint(amber, feedback, true), question.Explanation, "Enter to continue.")
	}
	return strings.Join(lines, "\n")
}
