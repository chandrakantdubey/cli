package course

import (
	"archive/tar"
	"compress/gzip"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"unicode"
	"unicode/utf8"
)

const maxPackFile = 1 << 20
const maxPackBytes = 8 << 20
const maxPackEntries = 1024

type packWriter struct {
	root    *os.Root
	bytes   int64
	entries int
}

func stagePack(source, destination string) error {
	root, err := os.OpenRoot(destination)
	if err != nil {
		return err
	}
	defer root.Close()
	writer := packWriter{root: root}
	source, err = filepath.Abs(source)
	if err != nil {
		return err
	}
	source, err = filepath.EvalSymlinks(source)
	if err != nil {
		return err
	}
	info, err := os.Stat(source)
	if err != nil {
		return err
	}
	if info.IsDir() {
		final, err := filepath.EvalSymlinks(destination)
		if err != nil {
			return err
		}
		if inside, err := filepath.Rel(source, final); err == nil && filepath.IsLocal(inside) {
			return errors.New("use a course directory outside the installation staging directory")
		}
		root, err := os.OpenRoot(source)
		if err != nil {
			return err
		}
		defer root.Close()
		return fs.WalkDir(root.FS(), ".", func(name string, entry fs.DirEntry, err error) error {
			if err != nil || name == "." {
				return err
			}
			if entry.IsDir() {
				return writer.directory(name)
			}
			if !entry.Type().IsRegular() {
				return fmt.Errorf("unsupported pack entry %q", name)
			}
			info, err := entry.Info()
			if err != nil {
				return err
			}
			file, err := root.Open(name)
			if err != nil {
				return err
			}
			defer file.Close()
			return writer.file(name, info.Size(), file)
		})
	}
	if !info.Mode().IsRegular() {
		return errors.New("provide a course directory or a .tar.gz archive")
	}
	archive, err := os.Open(source)
	if err != nil {
		return err
	}
	defer archive.Close()
	gz, err := gzip.NewReader(archive)
	if err != nil {
		return fmt.Errorf("read course archive: %w", err)
	}
	defer gz.Close()
	bounded := &io.LimitedReader{R: gz, N: maxPackBytes + maxPackEntries*4096}
	tarReader := tar.NewReader(bounded)
	for {
		entry, err := tarReader.Next()
		if errors.Is(err, io.EOF) {
			if _, err := io.Copy(io.Discard, bounded); err != nil {
				return err
			}
			if bounded.N == 0 {
				return errors.New("course archive exceeds expanded size limit")
			}
			return nil
		}
		if err != nil {
			return err
		}
		name := strings.TrimPrefix(entry.Name, "./")
		switch entry.Typeflag {
		case tar.TypeDir:
			name = strings.TrimSuffix(name, "/")
			if name != "" && name != "." {
				if err := writer.directory(name); err != nil {
					return err
				}
			}
		case tar.TypeReg:
			if err := writer.file(name, entry.Size, tarReader); err != nil {
				return err
			}
		default:
			return fmt.Errorf("unsupported archive entry %q", name)
		}
	}
}

func (w *packWriter) validate(name string) error {
	w.entries++
	if !fs.ValidPath(name) || name == "." || strings.Contains(name, "\\") || strings.ContainsFunc(name, unicode.IsControl) || w.entries > maxPackEntries {
		return fmt.Errorf("invalid or excessive course entry %q", name)
	}
	return nil
}

func (w *packWriter) directory(name string) error {
	if err := w.validate(name); err != nil {
		return err
	}
	return w.root.MkdirAll(name, 0700)
}

func (w *packWriter) file(name string, size int64, source io.Reader) error {
	if err := w.validate(name); err != nil {
		return err
	}
	if size < 0 || size > maxPackFile || w.bytes+size > maxPackBytes {
		return fmt.Errorf("course file exceeds size limit: %s", name)
	}
	if err := w.root.MkdirAll(filepath.Dir(name), 0700); err != nil {
		return err
	}
	file, err := w.root.OpenFile(name, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0644)
	if err != nil {
		return err
	}
	_, err = io.CopyN(file, source, size)
	w.bytes += size
	return errors.Join(err, file.Close())
}

func validatePackText(data []byte, manifest bool) error {
	if !utf8.Valid(data) || unsafePackText(string(data)) {
		return errors.New("lesson content contains unsupported control characters")
	}
	if manifest {
		var fields any
		if err := json.Unmarshal(data, &fields); err != nil {
			return err
		}
		if unsafeManifest(fields) {
			return errors.New("manifest contains unsupported control characters")
		}
	}
	return nil
}

func unsafePackText(value string) bool {
	return strings.ContainsFunc(value, func(r rune) bool {
		return unicode.IsControl(r) && r != '\t' && r != '\n' && r != '\r'
	})
}

func unsafeManifest(value any) bool {
	switch value := value.(type) {
	case string:
		return unsafePackText(value)
	case map[string]any:
		for key, field := range value {
			if unsafePackText(key) || unsafeManifest(field) {
				return true
			}
		}
	case []any:
		for _, field := range value {
			if unsafeManifest(field) {
				return true
			}
		}
	}
	return false
}
