package lab

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"cli-silly/course"
)

func TestFeedbackDescribesActualReport(t *testing.T) {
	for _, tc := range []struct{ actual, expected, want string }{
		{"", "a\n", "empty"}, {"a", "a\n", "newline"}, {"a\r\n", "a\n", "CRLF"},
		{"b\na\n", "a\nb\n", "order"}, {"a\n", "a\nb\n", "Expected 2 lines; found 1"},
		{"\x1b[2J\n", "a\n", "Line 1 differs"},
	} {
		got := textDifference([]byte(tc.actual), tc.expected)
		if !strings.Contains(got, tc.want) || strings.Contains(got, "\x1b") {
			t.Fatalf("feedback %q, want %q", got, tc.want)
		}
	}
}

func TestIntroRequiresPreservedInputsAndRetainsOldGrade(t *testing.T) {
	l, err := course.Load("foundations/cli-intro")
	if err != nil {
		t.Fatal(err)
	}
	s := Store{DataDir: t.TempDir(), StateDir: t.TempDir()}
	dir, err := s.Open(l)
	if err != nil {
		t.Fatal(err)
	}
	old := l
	old.GradingRevision = 1
	if err := writeJSON(s.progressPath(old), Progress{Revision: 1, Best: 100, Checks: 1, Latest: Result{Score: 100, Complete: true}}); err != nil {
		t.Fatal(err)
	}
	for _, c := range l.Checks {
		if err := os.WriteFile(filepath.Join(dir, c.Path), []byte(c.Expected), 0600); err != nil {
			t.Fatal(err)
		}
	}
	r, err := s.Check(l, dir)
	if err != nil || r.Score != 100 {
		t.Fatalf("%+v %v", r, err)
	}
	if err := os.Remove(filepath.Join(dir, "docs/checklist.txt")); err != nil {
		t.Fatal(err)
	}
	r, err = s.Check(l, dir)
	if err != nil || r.Score != 50 || !strings.Contains(r.Items[1].Feedback, "unchanged") {
		t.Fatalf("%+v %v", r, err)
	}
	p, err := s.Progress(l)
	if err != nil || p.PreviousBest != 100 || p.Best != 100 {
		t.Fatalf("%+v %v", p, err)
	}
	p, err = s.Progress(old)
	if err != nil || p.Checks != 1 || p.Best != 100 {
		t.Fatalf("old grade changed: %+v %v", p, err)
	}
}

func TestSpacedReviewDoesNotRewardRepeatedChecks(t *testing.T) {
	now := time.Date(2026, 11, 27, 12, 0, 0, 0, time.UTC)
	p := Progress{}
	advanceReview(&p, true, now)
	first := *p.ReviewAt
	advanceReview(&p, true, now.Add(time.Hour))
	if !p.ReviewAt.Equal(first) || p.ReviewStage != 0 {
		t.Fatal("early practice changed schedule")
	}
	advanceReview(&p, true, first)
	if p.ReviewStage != 1 || p.ReviewAt.Sub(first) != 72*time.Hour {
		t.Fatal("spaced pass did not advance")
	}
	advanceReview(&p, false, first.Add(time.Hour))
	if p.ReviewStage != 0 || p.ReviewAt.Sub(first.Add(time.Hour)) != 24*time.Hour {
		t.Fatal("failure did not shorten interval")
	}
}

func TestArchiveReviewRejectsDamagedStarterAfterPassing(t *testing.T) {
	l, err := course.Load("review/audit-archive")
	if err != nil {
		t.Fatal(err)
	}
	s := Store{DataDir: t.TempDir(), StateDir: t.TempDir()}
	dir, err := s.Open(l)
	if err != nil {
		t.Fatal(err)
	}
	runTask(t, dir, scenarios[l.ID].first)
	runTask(t, dir, scenarios[l.ID].finish)
	checkScore(t, s, l, dir, "finished", true)
	if err := os.WriteFile(filepath.Join(dir, "reports/summary.txt"), []byte("changed\n"), 0600); err != nil {
		t.Fatal(err)
	}
	result, err := s.Check(l, dir)
	if err != nil || result.Score != 50 || !strings.Contains(result.Items[0].Feedback, "unchanged") {
		t.Fatalf("accepted damaged starter: %+v %v", result, err)
	}
}
