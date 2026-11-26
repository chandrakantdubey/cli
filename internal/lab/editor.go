package lab

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"slices"
	"strings"

	"cli-silly/course"
)

// StarterFiles returns authored regular files, including paths with spaces.
func StarterFiles(lesson course.Lesson) ([]string, error) {
	var files []string
	if lesson.Starter == nil {
		return files, nil
	}
	err := fs.WalkDir(lesson.Starter, ".", func(path string, entry fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if entry.Type().IsRegular() && !entry.IsDir() {
			files = append(files, path)
		}
		return nil
	})
	return files, err
}

func (s Store) Editor(lesson course.Lesson, name string) (*exec.Cmd, error) {
	files, err := StarterFiles(lesson)
	if err != nil {
		return nil, err
	}
	if !slices.Contains(files, name) {
		return nil, errors.New("choose a file from the lab's starter list")
	}
	path, err := s.Open(lesson)
	if err != nil {
		return nil, err
	}
	root, err := os.OpenRoot(path)
	if err != nil {
		return nil, err
	}
	defer root.Close()
	// Check every path component before handing the filename to an external editor.
	parent := name
	for parent != "." {
		info, err := root.Lstat(parent)
		if err != nil {
			return nil, fmt.Errorf("open starter file: %w", err)
		}
		if (parent == name && !info.Mode().IsRegular()) || (parent != name && !info.IsDir()) {
			return nil, errors.New("the selected starter path must contain only directories and a regular file")
		}
		parent = filepath.Dir(parent)
	}

	editor := strings.TrimSpace(os.Getenv("VISUAL"))
	if editor == "" {
		editor = strings.TrimSpace(os.Getenv("EDITOR"))
	}
	if editor == "" {
		return nil, errors.New("set VISUAL or EDITOR to a terminal editor, for example EDITOR=nano")
	}
	if _, err := exec.LookPath(editor); err == nil {
		cmd := exec.Command(editor, "./"+name)
		cmd.Dir = path
		return cmd, nil
	}
	args := strings.Fields(editor)
	if _, err := exec.LookPath(args[0]); err != nil {
		return nil, fmt.Errorf("editor unavailable: %s", args[0])
	}
	cmd := exec.Command(args[0], append(args[1:], "./"+name)...)
	cmd.Dir = path
	return cmd, nil
}
