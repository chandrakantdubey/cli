package lab

import (
	"testing"
	"time"
)

func TestReviewKeepsQuizAnswersAndSchedulesAfterPractice(t *testing.T) {
	store := Store{StateDir: t.TempDir()}
	now := time.Date(2026, 11, 26, 11, 0, 0, 0, time.UTC)
	if _, err := store.RecordAnswer("streams/log-report", "fields", false, now); err != nil {
		t.Fatal(err)
	}
	if _, err := store.RecordAnswer("streams/log-report", "groups", true, now.Add(time.Minute)); err != nil {
		t.Fatal(err)
	}
	saved, err := store.Review()
	if err != nil {
		t.Fatal(err)
	}
	answers := saved.Answers["streams/log-report"]
	if len(answers) != 2 || answers["fields"].Correct || !answers["groups"].Correct || !answers["fields"].AnsweredAt.Equal(now) {
		t.Fatalf("incorrect persisted quiz answers: %+v", answers)
	}
	if !ReviewDue(Progress{Best: 100}, Progress{}).IsZero() {
		t.Fatal("older complete grades should be due")
	}
	base := Progress{Best: 100, LastPassed: &now}
	later := now.Add(5 * time.Hour)
	if ReviewDue(base, Progress{}).Sub(now) != 24*time.Hour ||
		ReviewDue(base, Progress{LastPassed: &later}).Sub(later) != 24*time.Hour {
		t.Fatal("due date should follow latest relevant passing check")
	}
}
