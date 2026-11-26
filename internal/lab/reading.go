package lab

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"slices"
)

// ReadingPosition identifies a section independently of terminal width.
type ReadingPosition struct {
	LessonID string `json:"lessonId"`
	Anchor   string `json:"anchor"`
}

type ReadingState struct {
	Version   int               `json:"version"`
	Last      ReadingPosition   `json:"last"`
	Bookmarks []ReadingPosition `json:"bookmarks"`
}

func (s Store) Reading() (ReadingState, error) {
	data, err := os.ReadFile(filepath.Join(s.StateDir, "reading.json"))
	if errors.Is(err, fs.ErrNotExist) {
		return ReadingState{Version: 1}, nil
	}
	if err != nil {
		return ReadingState{}, err
	}
	var state ReadingState
	if err := json.Unmarshal(data, &state); err != nil {
		return ReadingState{}, fmt.Errorf("read saved reading state: %w", err)
	}
	if state.Version != 1 {
		return ReadingState{}, fmt.Errorf("unsupported reading state version %d", state.Version)
	}
	return state, nil
}

func (s Store) SaveReading(position ReadingPosition, toggleBookmark bool) (ReadingState, error) {
	var state ReadingState
	err := withFileLock(filepath.Join(s.StateDir, "locks", "reading.lock"), func() error {
		var err error
		state, err = s.Reading()
		if err != nil {
			return err
		}
		state.Last = position
		if toggleBookmark {
			if i := slices.Index(state.Bookmarks, position); i >= 0 {
				state.Bookmarks = slices.Delete(state.Bookmarks, i, i+1)
			} else {
				state.Bookmarks = append(state.Bookmarks, position)
			}
		}
		return writeJSON(filepath.Join(s.StateDir, "reading.json"), state)
	})
	return state, err
}
