package tui

import (
	"strings"
	"time"

	"cli-silly/course"

	tea "github.com/charmbracelet/bubbletea"
)

type motionKind int

const (
	still motionKind = iota
	demoMotion
	checkMotion
	scoreMotion
)

const motionFrames = 7

type motion struct {
	kind       motionKind
	generation int
	step       int
	fromScore  int
}

type motionTick int

func motionDelay(generation int) tea.Cmd {
	return tea.Tick(70*time.Millisecond, func(time.Time) tea.Msg { return motionTick(generation) })
}

func (m *model) beginMotion(kind motionKind, fromScore int) tea.Cmd {
	m.motion.kind = still
	if m.reduceMotion || !m.useColor {
		return nil
	}
	m.motion = motion{kind: kind, generation: m.motion.generation + 1, fromScore: fromScore}
	return motionDelay(m.motion.generation)
}

func (m model) advanceMotion(tick motionTick) (tea.Model, tea.Cmd) {
	if int(tick) != m.motion.generation || m.motion.kind == still {
		return m, nil
	}
	switch m.motion.kind {
	case demoMotion:
		if m.screen != demoScreen {
			m.motion.kind = still
		}
	case checkMotion:
		if !m.busy {
			m.motion.kind = still
		}
	case scoreMotion:
		if m.screen != resultsScreen {
			m.motion.kind = still
		}
	}
	m.motion.step++
	if m.motion.kind == still || (m.motion.kind != checkMotion && m.motion.step >= motionFrames) {
		m.motion.kind = still
		return m, nil
	}
	return m, motionDelay(m.motion.generation)
}

func hasDemoMotion(frame course.Frame) bool {
	return len(frame.Flow) > 0 || frame.OutputFile != "" || len(frame.Permissions) == 9
}

func (m *model) playDemo() tea.Cmd {
	if hasDemoMotion(m.lesson.Demo[m.frame]) {
		return m.beginMotion(demoMotion, 0)
	}
	return nil
}

func (m model) flowDiagram(stages []string) string {
	nodes := make([]string, 0, len(stages))
	focus := -1
	if m.motion.kind == demoMotion {
		focus = m.motion.step * (len(stages) - 1) / (motionFrames - 1)
	}
	for i, stage := range stages {
		if i == focus {
			nodes = append(nodes, m.paint(amber, "◆ "+stage, true))
		} else {
			nodes = append(nodes, m.paint(teal, "◇ "+stage, false))
		}
	}
	return strings.Join(nodes, m.paint(muted, " → ", false))
}

func (m model) stdoutArrow() string {
	if m.motion.kind != demoMotion {
		return "━━━━━━━▶"
	}
	step := m.motion.step
	return strings.Repeat("━", step) + m.paint(amber, "●", true) + strings.Repeat("━", motionFrames-1-step) + "▶"
}

func (m model) scoreFill(score int) int {
	if m.motion.kind != scoreMotion {
		return score
	}
	from := m.motion.fromScore
	return from + (score-from)*m.motion.step/(motionFrames-1)
}

func (m model) checkingStatus() string {
	if m.motion.kind != checkMotion {
		return m.status
	}
	spinner := []string{"⠋", "⠙", "⠹", "⠸", "⠼", "⠴", "⠦", "⠧", "⠇", "⠏"}
	return spinner[m.motion.step%len(spinner)] + " " + m.status
}
