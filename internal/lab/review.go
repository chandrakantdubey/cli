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

// ReviewDue uses a fixed initial interval. Older completed records are due immediately.
func ReviewDue(base, variant Progress) time.Time {
	if variant.LastPassed != nil {
		return variant.LastPassed.Add(24 * time.Hour)
	}
	if base.LastPassed != nil {
		return base.LastPassed.Add(24 * time.Hour)
	}
	return time.Time{}
}
