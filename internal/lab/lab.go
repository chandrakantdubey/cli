package lab

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"syscall"
	"time"

	"cli-silly/course"
)

const Marker = ".cli-silly-lab.json"

type Store struct {
	DataDir  string
	StateDir string
}

type workspace struct {
	LessonID string `json:"lessonId"`
	Revision int    `json:"starterRevision"`
	Project  string `json:"composeProject,omitempty"`
}

type Item struct {
	Label    string `json:"label"`
	Feedback string `json:"feedback,omitempty"`
	Earned   int    `json:"earned"`
	Possible int    `json:"possible"`
}

type Result struct {
	Items    []Item `json:"items"`
	Score    int    `json:"score"`
	Complete bool   `json:"complete"`
}

type Progress struct {
	PreviousBest int        `json:"previousBest,omitempty"`
	ReviewStage  int        `json:"reviewStage,omitempty"`
	ReviewAt     *time.Time `json:"reviewAt,omitempty"`
	ReviewedAt   *time.Time `json:"reviewedAt,omitempty"`
	Revision     int        `json:"gradingRevision"`
	Best         int        `json:"best"`
	Checks       int        `json:"checks"`
	Latest       Result     `json:"latest"`
	LastPassed   *time.Time `json:"lastPassed,omitempty"`
}

func (s Store) Path(lesson course.Lesson) string {
	return filepath.Join(s.DataDir, "workspaces", lesson.ID)
}

func (s Store) progressPath(lesson course.Lesson) string {
	return filepath.Join(s.StateDir, "progress", lesson.ID, fmt.Sprintf("%d.json", lesson.GradingRevision))
}

func (s Store) Open(lesson course.Lesson) (string, error) {
	path := s.Path(lesson)
	err := s.withLock(lesson, func() error {
		if _, err := os.Lstat(path); !errors.Is(err, fs.ErrNotExist) {
			if err != nil {
				return err
			}
			return verifyWorkspace(path, lesson)
		}
		return create(path, lesson)
	})
	return path, err
}

func (s Store) Reset(lesson course.Lesson) (string, error) {
	return s.ResetFrom(lesson, "")
}

func (s Store) ResetFrom(lesson course.Lesson, expectedDirectory string) (string, error) {
	path := s.Path(lesson)
	var backup string
	err := s.withLock(lesson, func() error {
		if err := verifyAttempt(path, expectedDirectory); err != nil {
			return err
		}
		if err := verifyWorkspace(path, lesson); err != nil {
			return err
		}
		archive, err := os.MkdirTemp(filepath.Dir(path), ".saved-")
		if err != nil {
			return err
		}
		backup = filepath.Join(archive, filepath.Base(path))
		if err := os.Rename(path, backup); err != nil {
			return err
		}
		if err := create(path, lesson); err != nil {
			return errors.Join(err, os.Rename(backup, path))
		}
		return nil
	})
	return backup, err
}

func create(path string, lesson course.Lesson) error {
	parent := filepath.Dir(path)
	if err := os.MkdirAll(parent, 0700); err != nil {
		return err
	}
	tmp, err := os.MkdirTemp(parent, ".new-")
	if err != nil {
		return err
	}
	defer os.RemoveAll(tmp)

	err = fs.WalkDir(lesson.Starter, ".", func(name string, entry fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		target := filepath.Join(tmp, name)
		if entry.IsDir() {
			return os.MkdirAll(target, 0700)
		}
		data, err := fs.ReadFile(lesson.Starter, name)
		if err != nil {
			return err
		}
		return os.WriteFile(target, data, 0644)
	})
	if err != nil {
		return err
	}
	marker := workspace{LessonID: lesson.ID, Revision: lesson.StarterRevision}
	if err := initialize(tmp, lesson.Setup, &marker); err != nil {
		return err
	}
	if err := writeJSON(filepath.Join(tmp, Marker), marker); err != nil {
		return err
	}
	return os.Rename(tmp, path)
}

func verifyWorkspace(path string, lesson course.Lesson) error {
	info, err := os.Lstat(path)
	if err != nil {
		return fmt.Errorf("open the lab first: %w", err)
	}
	if !info.IsDir() {
		return fmt.Errorf("expected a managed directory at %s", path)
	}
	root, err := os.OpenRoot(path)
	if err != nil {
		return err
	}
	defer root.Close()
	data, err := readRegular(root, Marker)
	if err != nil {
		return fmt.Errorf("read workspace marker: %w", err)
	}
	var marker workspace
	if err := json.Unmarshal(data, &marker); err != nil || marker.LessonID != lesson.ID || marker.Revision < 1 {
		return fmt.Errorf("workspace marker does not match %s", lesson.ID)
	}
	return nil
}

func verifyAttempt(path, expectedDirectory string) error {
	if expectedDirectory == "" {
		return nil
	}
	expected, err := os.Stat(expectedDirectory)
	if err != nil {
		return err
	}
	actual, err := os.Stat(path)
	if err != nil || !os.SameFile(expected, actual) {
		return errors.New("this shell belongs to an archived attempt; reopen the lab")
	}
	return nil
}

func (s Store) Check(lesson course.Lesson, expectedDirectory string) (Result, error) {
	return s.check(context.Background(), lesson, expectedDirectory, nil, false)
}

// WatchCheck grades only local file criteria and saves a result when feedback changes.
func (s Store) WatchCheck(ctx context.Context, lesson course.Lesson, attempt fs.FileInfo) (Result, error) {
	if !Watchable(lesson) {
		return Result{}, errors.New("use a manual check for this lab")
	}
	return s.check(ctx, lesson, "", attempt, true)
}

func Watchable(lesson course.Lesson) bool {
	return len(lesson.Checks) > 0 && !slices.ContainsFunc(lesson.Checks, func(item course.Criterion) bool {
		return item.Kind != "text" && item.Kind != "owner-executable" && item.Kind != "absent"
	})
}

func (s Store) check(ctx context.Context, lesson course.Lesson, expectedDirectory string, attempt fs.FileInfo, watch bool) (Result, error) {
	var result Result
	err := s.withLock(lesson, func() error {
		path := s.Path(lesson)
		if err := ctx.Err(); err != nil {
			return err
		}
		if attempt != nil {
			current, err := os.Stat(path)
			if err != nil || !os.SameFile(attempt, current) {
				return errors.New("the watched attempt was reset; reopen it to restart automatic checks")
			}
		}
		if err := verifyAttempt(path, expectedDirectory); err != nil {
			return err
		}
		if err := Ready(lesson); err != nil {
			return err
		}
		if err := verifyWorkspace(path, lesson); err != nil {
			return err
		}
		root, err := os.OpenRoot(path)
		if err != nil {
			return err
		}
		defer root.Close()

		for _, criterion := range lesson.Checks {
			item, err := evaluate(root, path, criterion)
			if err != nil {
				return fmt.Errorf("check %s: %w", criterion.ID, err)
			}
			for _, name := range criterion.Preserve {
				expected, err := fs.ReadFile(lesson.Starter, name)
				if err != nil {
					return err
				}
				actual, err := readRegular(root, name)
				if err != nil || string(actual) != string(expected) {
					item.Earned = 0
					item.Feedback = fmt.Sprintf("Keep starter file lab/%s unchanged. Inspect it or use cli-silly reset to archive this attempt and restore the starter files.", name)
					break
				}
			}
			result.Items = append(result.Items, item)
			result.Score += item.Earned
		}
		result.Complete = result.Score == 100

		progress, err := s.Progress(lesson)
		if err != nil {
			return err
		}
		if err := ctx.Err(); err != nil {
			return err
		}
		if watch && progress.Checks > 0 && slices.Equal(progress.Latest.Items, result.Items) {
			return nil
		}
		progress.Latest = result
		progress.Best = max(progress.Best, result.Score)
		progress.Checks++
		if lesson.ReviewOf != "" {
			advanceReview(&progress, result.Complete, time.Now().UTC())
		}
		if result.Complete {
			now := time.Now().UTC()
			progress.LastPassed = &now
		}
		return writeJSON(s.progressPath(lesson), progress)
	})
	return result, err
}

var errNotRegular = errors.New("expected a regular file")
var errTooLarge = errors.New("file exceeds the 64 KiB report limit")

func evaluate(root *os.Root, directory string, check course.Criterion) (Item, error) {
	item, err := evaluateCriterion(root, directory, check)
	if err != nil || item.Earned == item.Possible {
		return item, err
	}
	if check.Feedback != "" {
		item.Feedback = check.Feedback
	}
	if check.Kind == "text" {
		if _, err := root.Lstat(check.Path); errors.Is(err, fs.ErrNotExist) {
			if misplaced := misplacedReport(root, check.Path); misplaced != "" {
				item.Feedback = fmt.Sprintf("Found lab/%s. Save the report at lab/%s. %s", misplaced, check.Path, item.Feedback)
			} else {
				item.Feedback = fmt.Sprintf("Expected lab/%s. %s", check.Path, item.Feedback)
			}
		}
		if data, err := readRegular(root, check.Path); err == nil {
			item.Feedback = textDifference(data, check.Expected) + " " + item.Feedback
		}
	}
	return item, nil
}

// Look only one directory below the expected location, with a bounded directory read.
func misplacedReport(root *os.Root, path string) string {
	parent := filepath.Dir(path)
	dir, err := root.Open(parent)
	if err != nil {
		return ""
	}
	defer dir.Close()
	entries, _ := dir.ReadDir(32)
	for _, entry := range entries {
		if entry.IsDir() && !strings.HasPrefix(entry.Name(), ".") {
			location := filepath.Join(parent, entry.Name(), filepath.Base(path))
			if info, err := root.Lstat(location); err == nil && info.Mode().IsRegular() {
				return location
			}
		}
	}
	return ""
}

func evaluateCriterion(root *os.Root, directory string, check course.Criterion) (Item, error) {
	item := Item{Label: check.Label, Possible: check.Points}
	switch check.Kind {
	case "git-count", "git-tree", "git-index", "git-branch", "git-ancestor", "git-file":
		return match(item, check, gitCriterion(directory, root, check))
	case "docker-running", "docker-healthy", "docker-recreated":
		value, err := dockerCriterion(root, check)
		if err != nil {
			return item, err
		}
		return match(item, check, value)
	case "http-body":
		return match(item, check, httpBody(root, check.Path))
	case "archive-member":
		return match(item, check, archiveMember(root, check.Path, check.Member))
	case "script-config", "script-report", "node-cli":
		value, err := exerciseScript(root, check)
		if err != nil {
			return item, err
		}
		return match(item, check, value)
	case "absent":
		_, err := root.Lstat(check.Path)
		if errors.Is(err, fs.ErrNotExist) {
			item.Earned = check.Points
		} else if err != nil {
			return item, err
		} else {
			item.Feedback = "The starter file is still present."
		}
		return item, nil
	}
	if check.Kind == "owner-executable" {
		info, err := root.Lstat(check.Path)
		switch {
		case errors.Is(err, fs.ErrNotExist):
			item.Feedback = "Restore " + check.Path + " in the lab."
		case err != nil:
			return item, err
		case !info.Mode().IsRegular() || info.Mode().Perm()&0100 == 0:
			item.Feedback = "Inspect " + check.Path + " with ls -l and check the owner's execute bit."
		default:
			item.Earned = item.Possible
		}
		return item, nil
	}

	data, err := readRegular(root, check.Path)
	switch {
	case errors.Is(err, fs.ErrNotExist):
		item.Feedback = "Create " + check.Path + " by redirecting the script's output."
	case errors.Is(err, errNotRegular), errors.Is(err, errTooLarge):
		item.Feedback = "Replace " + check.Path + " with the plain text report from the script."
	case err != nil:
		return item, err
	case string(data) != check.Expected:
		item.Feedback = "Compare " + check.Path + " with the script's standard output, including its final newline."
	default:
		item.Earned = item.Possible
	}
	return item, nil
}

func readRegular(root *os.Root, path string) ([]byte, error) {
	info, err := root.Lstat(path)
	if err != nil {
		return nil, err
	}
	if !info.Mode().IsRegular() {
		return nil, errNotRegular
	}
	file, err := root.OpenFile(path, os.O_RDONLY|syscall.O_NONBLOCK|syscall.O_NOFOLLOW, 0)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	info, err = file.Stat()
	if err != nil {
		return nil, err
	}
	if !info.Mode().IsRegular() {
		return nil, errNotRegular
	}
	const limit = 64 << 10
	data, err := io.ReadAll(io.LimitReader(file, limit+1))
	if err == nil && len(data) > limit {
		err = errTooLarge
	}
	return data, err
}

func (s Store) Progress(lesson course.Lesson) (Progress, error) {
	progress := Progress{Revision: lesson.GradingRevision}
	data, err := os.ReadFile(s.progressPath(lesson))
	if errors.Is(err, fs.ErrNotExist) {
		files, readErr := os.ReadDir(filepath.Dir(s.progressPath(lesson)))
		if readErr != nil && !errors.Is(readErr, fs.ErrNotExist) {
			return progress, readErr
		}
		for _, file := range files {
			if file.IsDir() || !strings.HasSuffix(file.Name(), ".json") {
				continue
			}
			var old Progress
			data, readErr := os.ReadFile(filepath.Join(filepath.Dir(s.progressPath(lesson)), file.Name()))
			if readErr != nil {
				return progress, readErr
			}
			if json.Unmarshal(data, &old) != nil || old.Revision >= lesson.GradingRevision || old.Revision < 1 || old.Best < 0 || old.Best > 100 {
				continue
			}
			progress.PreviousBest = max(progress.PreviousBest, old.Best)
		}
		return progress, nil
	}
	if err != nil {
		return progress, err
	}
	err = json.Unmarshal(data, &progress)
	if err == nil && (progress.Revision != lesson.GradingRevision || progress.Best < 0 || progress.Best > 100 || progress.Latest.Score < 0 || progress.Latest.Score > progress.Best || progress.Checks < 0) {
		err = errors.New("invalid saved progress")
	}
	return progress, err
}

func (s Store) withLock(lesson course.Lesson, run func() error) error {
	return withFileLock(filepath.Join(s.StateDir, "locks", lesson.ID+".lock"), run)
}

func withFileLock(path string, run func() error) error {
	if err := os.MkdirAll(filepath.Dir(path), 0700); err != nil {
		return err
	}
	file, err := os.OpenFile(path, os.O_CREATE|os.O_RDWR, 0600)
	if err != nil {
		return err
	}
	defer file.Close()
	for {
		err = syscall.Flock(int(file.Fd()), syscall.LOCK_EX)
		if !errors.Is(err, syscall.EINTR) {
			break
		}
	}
	if err != nil {
		return err
	}
	return run()
}

func writeJSON(path string, value any) error {
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0700); err != nil {
		return err
	}
	file, err := os.CreateTemp(dir, ".write-")
	if err != nil {
		return err
	}
	defer os.Remove(file.Name())
	writeErr := json.NewEncoder(file).Encode(value)
	if err := errors.Join(writeErr, file.Close()); err != nil {
		return err
	}
	return os.Rename(file.Name(), path)
}
