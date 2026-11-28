package lab

import (
	"encoding/json"
	"errors"
	"io/fs"
	"os"
	"path/filepath"
	"time"
)

type QuizAnswer struct {
	Correct    bool      `json:"correct"`
	AnsweredAt time.Time `json:"answeredAt"`
}

type ReviewState struct {
	Version int                              `json:"version"`
	Answers map[string]map[string]QuizAnswer `json:"answers"`
}

func (s Store) Review() (ReviewState, error) {
	state := ReviewState{Version: 1, Answers: make(map[string]map[string]QuizAnswer)}
	data, err := os.ReadFile(filepath.Join(s.StateDir, "review.json"))
	if errors.Is(err, fs.ErrNotExist) {
		return state, nil
	}
	if err != nil {
		return ReviewState{}, err
	}
	if err := json.Unmarshal(data, &state); err != nil {
		return ReviewState{}, err
	}
	if state.Version != 1 || state.Answers == nil {
		return ReviewState{}, errors.New("invalid review state")
	}
	return state, nil
}

func (s Store) RecordAnswer(lessonID, questionID string, correct bool, now time.Time) (ReviewState, error) {
	var state ReviewState
	err := withFileLock(filepath.Join(s.StateDir, "locks", "review.lock"), func() error {
		var err error
		state, err = s.Review()
		if err != nil {
			return err
		}
		if state.Answers[lessonID] == nil {
			state.Answers[lessonID] = make(map[string]QuizAnswer)
		}
		state.Answers[lessonID][questionID] = QuizAnswer{correct, now.UTC()}
		return writeJSON(filepath.Join(s.StateDir, "review.json"), state)
	})
	return state, err
}

// ReviewDue preserves older records while new attempts use spaced intervals.
func ReviewDue(base, variant Progress) time.Time {
	if variant.ReviewAt != nil {
		return *variant.ReviewAt
	}
	if variant.LastPassed != nil {
		return variant.LastPassed.Add(24 * time.Hour)
	}
	if base.LastPassed != nil {
		return base.LastPassed.Add(24 * time.Hour)
	}
	return time.Time{}
}

// Only a passing attempt at or after its due time advances the schedule. Early
// practice and repeated checks are welcome, but do not postpone the next review.
func advanceReview(p *Progress, passed bool, now time.Time) {
	intervals := [...]time.Duration{24 * time.Hour, 3 * 24 * time.Hour, 7 * 24 * time.Hour, 14 * 24 * time.Hour, 30 * 24 * time.Hour}
	if !passed {
		p.ReviewStage = 0
		due := now.Add(intervals[0])
		if p.ReviewAt == nil || due.Before(*p.ReviewAt) {
			p.ReviewAt = &due
		}
		return
	}
	if p.ReviewedAt != nil && p.ReviewAt != nil && now.Before(*p.ReviewAt) {
		return
	}
	if p.ReviewedAt != nil {
		p.ReviewStage = min(max(0, p.ReviewStage)+1, len(intervals)-1)
	}
	p.ReviewedAt = &now
	due := now.Add(intervals[max(0, min(p.ReviewStage, len(intervals)-1))])
	p.ReviewAt = &due
}
