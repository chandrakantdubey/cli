package lab

import (
	"errors"
	"io/fs"
	"os"

	"cli-silly/course"
)

// WriteTemplate creates a new file in the current managed workspace.
func (s Store) WriteTemplate(lesson course.Lesson, attempt, path, contents string) error {
	if !fs.ValidPath(path) || path == "." || path == Marker {
		return errors.New("choose a relative file path inside the workspace")
	}
	return s.withLock(lesson, func() error {
		workspace := s.Path(lesson)
		if err := verifyAttempt(workspace, attempt); err != nil {
			return err
		}
		if err := verifyWorkspace(workspace, lesson); err != nil {
			return err
		}
		root, err := os.OpenRoot(workspace)
		if err != nil {
			return err
		}
		defer root.Close()
		file, err := root.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0644)
		if errors.Is(err, fs.ErrExist) {
			return errors.New("the destination already exists; choose another filename")
		}
		if err != nil {
			return err
		}
		_, writeErr := file.WriteString(contents)
		if err := errors.Join(writeErr, file.Close()); err != nil {
			return errors.Join(err, root.Remove(path))
		}
		return nil
	})
}
