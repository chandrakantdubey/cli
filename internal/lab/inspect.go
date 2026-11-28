package lab

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"cli-silly/course"
)

// Inspect reports existing state only. It never prepares a lab or runs its code.
func (s Store) Inspect(l course.Lesson) (string, error) {
	dir := s.Path(l)
	if err := verifyWorkspace(dir, l); err != nil {
		return "", err
	}
	root, err := os.OpenRoot(dir)
	if err != nil {
		return "", err
	}
	defer root.Close()
	lines := []string{"WORKSPACE INSPECTION", l.Title, "Root: " + strconv.Quote(dir), "", "Existing files (maximum 100 entries, four levels):"}
	count := 0
	err = fs.WalkDir(root.FS(), ".", func(name string, entry fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if name == "." {
			return nil
		}
		if entry.IsDir() && (entry.Name() == ".git" || strings.Count(name, "/") >= 4) {
			return fs.SkipDir
		}
		if count >= 100 {
			lines = append(lines, "... more entries omitted")
			return fs.SkipAll
		}
		count++
		info, err := entry.Info()
		if err != nil {
			return err
		}
		lines = append(lines, fmt.Sprintf("%s  %7d B  %s", info.Mode(), info.Size(), strconv.Quote(name)))
		return nil
	})
	if err != nil {
		return "", err
	}
	if strings.HasPrefix(l.Setup, "git-") {
		info, err := root.Lstat(".git")
		if err != nil || !info.IsDir() {
			return "", fmt.Errorf("Git inspection requires the lab's own .git directory; restore a saved attempt or reset this lab")
		}
		status, err := git(dir, "status", "--short", "--branch")
		if err != nil {
			return "", err
		}
		lines = append(lines, "", "Git working tree and index:", safeTerminal(status))
		graph, err := git(dir, "log", "--all", "--graph", "--oneline", "--decorate=short", "-8")
		if err != nil {
			return "", err
		}
		lines = append(lines, "Git history (last eight commits):", safeTerminal(graph))
	}
	lines = append(lines, "", "This is your current filesystem state. Bash runs with your Linux user permissions.", "Use cli-silly archives to find saved attempts. Inspection does not change your score.")
	return strings.Join(lines, "\n"), nil
}

func safeTerminal(value string) string {
	var out strings.Builder
	for _, r := range value {
		if r == '\n' || r == '\t' || (r >= 32 && r < 127) {
			out.WriteRune(r)
		} else {
			out.WriteString(strconv.QuoteRuneToASCII(r))
		}
	}
	return out.String()
}

func (s Store) Archives(l course.Lesson) ([]string, error) {
	parent := filepath.Dir(s.Path(l))
	entries, err := os.ReadDir(parent)
	if os.IsNotExist(err) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	var result []string
	for _, e := range entries {
		if !e.IsDir() || !strings.HasPrefix(e.Name(), ".saved-") {
			continue
		}
		path := filepath.Join(parent, e.Name(), filepath.Base(l.ID))
		if verifyWorkspace(path, l) == nil {
			result = append(result, path)
		}
	}
	return result, nil
}

// Restore swaps a validated saved attempt with the active workspace. Both attempts
// remain recoverable; earned scores are untouched.
func (s Store) Restore(l course.Lesson, archive string) (string, error) {
	if !strings.HasPrefix(archive, ".saved-") || filepath.Base(archive) != archive {
		return "", fmt.Errorf("use the .saved-... name shown by cli-silly archives")
	}
	current := s.Path(l)
	parent := filepath.Dir(current)
	saved := filepath.Join(parent, archive, filepath.Base(current))
	var backup string
	err := s.withLock(l, func() error {
		info, err := os.Lstat(filepath.Dir(saved))
		if err != nil {
			return err
		}
		if !info.IsDir() {
			return fmt.Errorf("archive must be a managed directory")
		}
		if err := verifyWorkspace(saved, l); err != nil {
			return err
		}
		if err := verifyWorkspace(current, l); err != nil {
			return err
		}
		holder, err := os.MkdirTemp(parent, ".saved-")
		if err != nil {
			return err
		}
		backup = filepath.Join(holder, filepath.Base(current))
		if err := os.Rename(current, backup); err != nil {
			return err
		}
		if err := os.Rename(saved, current); err != nil {
			if restoreErr := os.Rename(backup, current); restoreErr != nil {
				return fmt.Errorf("restore failed: %v; rollback failed: %w; previous workspace is at %s", err, restoreErr, backup)
			}
			return err
		}
		return nil
	})
	return backup, err
}
