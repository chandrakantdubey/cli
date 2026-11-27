package tui

import (
	"strings"
	"testing"

	"cli-silly/course"

	tea "github.com/charmbracelet/bubbletea"
)

func TestDemoMotionReplaysAndStopsWhenLeaving(t *testing.T) {
	lesson, err := course.Load("foundations/cli-intro")
	if err != nil {
		t.Fatal(err)
	}
	m := model{lesson: lesson, screen: demoScreen, frame: 1, useColor: true}
	replay := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'a'}}
	updated, cmd := m.Update(replay)
	m = updated.(model)
	if cmd == nil || m.motion.kind != demoMotion || !strings.Contains(m.stdoutArrow(), "●") {
		t.Fatal("replay should start the redirect trace")
	}
	old := m.motion.generation
	start := m.stdoutArrow()
	updated, _ = m.Update(motionTick(old))
	m = updated.(model)
	if m.stdoutArrow() == start {
		t.Fatal("the mark should move along the arrow")
	}
	updated, cmd = m.Update(replay)
	m = updated.(model)
	if cmd == nil || m.stdoutArrow() != start {
		t.Fatal("replay should restart at the beginning")
	}
	updated, cmd = m.Update(motionTick(old))
	m = updated.(model)
	if cmd != nil || m.stdoutArrow() != start {
		t.Fatal("a tick from the last replay should be ignored")
	}
	updated, _ = m.Update(tea.KeyMsg{Type: tea.KeyTab})
	m = updated.(model)
	if m.screen != taskScreen || m.motion.kind != still {
		t.Fatal("switching modes should stop the animation")
	}
}

func TestReducedMotionAndAnimatedScoreDisplay(t *testing.T) {
	m := model{screen: resultsScreen, useColor: true, reduceMotion: true}
	if m.beginMotion(scoreMotion, 40) != nil || m.scoreFill(100) != 100 {
		t.Fatal("reduced motion should show the final score immediately")
	}
	m.reduceMotion = false
	m.beginMotion(scoreMotion, 40)
	if m.scoreFill(100) != 40 || !strings.Contains(m.scoreMeter(m.scoreFill(100), 100), "100/100") {
		t.Fatal("only the meter should interpolate while the number stays measured")
	}
	for i := 0; i < motionFrames; i++ {
		updated, _ := m.Update(motionTick(m.motion.generation))
		m = updated.(model)
	}
	if m.motion.kind != still || m.scoreFill(100) != 100 {
		t.Fatal("the animation should stop at the final score")
	}
}
