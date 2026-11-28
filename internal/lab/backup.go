package lab

import (
	"archive/tar"
	"compress/gzip"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path"
	"path/filepath"
	"strings"
	"syscall"
)

const backupLimit = 256 << 20
const backupManifest = "cli-silly-backup-v1\n"

// ExportBackup copies directories and regular files. Symlinks and special files are refused,
// never followed. Close other sessions first for a consistent filesystem snapshot.
func (s Store) ExportBackup(destination string) (err error) {
	abs, err := filepath.Abs(destination)
	if err != nil {
		return err
	}
	parent, err := filepath.EvalSymlinks(filepath.Dir(abs))
	if err != nil {
		return err
	}
	abs = filepath.Join(parent, filepath.Base(abs))
	for _, dir := range []string{s.DataDir, s.StateDir} {
		base, e := filepath.Abs(dir)
		if e != nil {
			return e
		}
		if resolved, resolveErr := filepath.EvalSymlinks(base); resolveErr == nil {
			base = resolved
		}
		rel, e := filepath.Rel(base, abs)
		if e == nil && rel != ".." && !strings.HasPrefix(rel, ".."+string(filepath.Separator)) {
			return errors.New("save the backup outside cli-silly data and state directories")
		}
	}
	file, err := os.OpenFile(abs, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0600)
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			os.Remove(abs)
		}
	}()
	z := gzip.NewWriter(file)
	archive := tar.NewWriter(z)
	defer func() { err = errors.Join(err, archive.Close(), z.Close(), file.Close()) }()
	if err = archive.WriteHeader(&tar.Header{Name: "backup-version", Mode: 0600, Size: int64(len(backupManifest))}); err != nil {
		return err
	}
	if _, err = io.WriteString(archive, backupManifest); err != nil {
		return err
	}
	var size int64
	count := 0
	for _, part := range []struct{ name, dir string }{{"data", s.DataDir}, {"state", s.StateDir}} {
		root, e := os.OpenRoot(part.dir)
		if errors.Is(e, fs.ErrNotExist) {
			continue
		}
		if e != nil {
			return e
		}
		e = fs.WalkDir(root.FS(), ".", func(name string, entry fs.DirEntry, walkErr error) error {
			if walkErr != nil {
				return walkErr
			}
			if name == "." {
				return nil
			}
			if name == "locks" {
				return fs.SkipDir
			}
			if strings.HasPrefix(entry.Name(), ".write-") || strings.HasPrefix(entry.Name(), ".new-") {
				if entry.IsDir() {
					return fs.SkipDir
				}
				return nil
			}
			info, e := entry.Info()
			if e != nil {
				return e
			}
			if !entry.IsDir() && !info.Mode().IsRegular() {
				return fmt.Errorf("backup refuses symlink or special file %q; move it aside or export manually", name)
			}
			count++
			size += info.Size()
			if size > backupLimit || count > 20000 {
				return errors.New("backup exceeds 256 MiB or 20000 entries; archive large practice artifacts separately")
			}
			if entry.IsDir() {
				return archive.WriteHeader(&tar.Header{Name: path.Join(part.name, name), Typeflag: tar.TypeDir, Mode: 0700})
			}
			f, e := root.OpenFile(name, os.O_RDONLY|syscall.O_NOFOLLOW|syscall.O_NONBLOCK, 0)
			if e != nil {
				return e
			}
			defer f.Close()
			opened, e := f.Stat()
			if e != nil {
				return e
			}
			if !opened.Mode().IsRegular() || !os.SameFile(info, opened) || opened.Size() != info.Size() {
				return fmt.Errorf("file changed during backup: %q; close other sessions and retry", name)
			}
			if e = archive.WriteHeader(&tar.Header{Name: path.Join(part.name, name), Mode: int64(info.Mode().Perm()), Size: info.Size()}); e != nil {
				return e
			}
			_, e = io.CopyN(archive, f, info.Size())
			return e
		})
		root.Close()
		if e != nil {
			return e
		}
	}
	return nil
}

// ImportBackup is deliberately non-merging. Validate into temporary directories
// before publishing into empty destinations. Existing learner data is never replaced.
func (s Store) ImportBackup(source string) error {
	dataAbs, e := filepath.Abs(s.DataDir)
	if e != nil {
		return e
	}
	stateAbs, e := filepath.Abs(s.StateDir)
	if e != nil {
		return e
	}
	if dataAbs == stateAbs || strings.HasPrefix(dataAbs, stateAbs+"/") || strings.HasPrefix(stateAbs, dataAbs+"/") {
		return errors.New("data and state destinations must be separate directories")
	}
	for _, dir := range []string{dataAbs, stateAbs} {
		entries, e := os.ReadDir(dir)
		if e != nil && !errors.Is(e, fs.ErrNotExist) {
			return e
		}
		if len(entries) > 0 {
			return fmt.Errorf("destination is not empty: %s; import into new XDG data/state directories, or move your existing directories aside first", dir)
		}
		if info, e := os.Lstat(dir); e == nil && !info.IsDir() {
			return fmt.Errorf("destination must be a directory: %s", dir)
		}
	}
	staged := make(map[string]string)
	for _, part := range []struct{ name, dir string }{{"data", dataAbs}, {"state", stateAbs}} {
		if e := os.MkdirAll(filepath.Dir(part.dir), 0700); e != nil {
			return e
		}
		dir, e := os.MkdirTemp(filepath.Dir(part.dir), ".cli-silly-import-")
		if e != nil {
			return e
		}
		staged[part.name] = dir
		defer os.RemoveAll(dir)
	}
	f, e := os.Open(source)
	if e != nil {
		return e
	}
	defer f.Close()
	z, e := gzip.NewReader(f)
	if e != nil {
		return e
	}
	defer z.Close()
	r := tar.NewReader(io.LimitReader(z, backupLimit+(32<<20)))
	seen := make(map[string]bool)
	valid := false
	var size int64
	for {
		h, e := r.Next()
		if e == io.EOF {
			break
		}
		if e != nil {
			return e
		}
		if seen[h.Name] || !fs.ValidPath(h.Name) || strings.Contains(h.Name, "\\") || (h.Typeflag != tar.TypeReg && h.Typeflag != tar.TypeDir) || h.Size < 0 || (h.Typeflag == tar.TypeDir && h.Size != 0) {
			return fmt.Errorf("invalid backup entry %q", h.Name)
		}
		seen[h.Name] = true
		if len(seen) > 20001 {
			return errors.New("backup contains too many entries")
		}
		size += h.Size
		if size > backupLimit {
			return errors.New("backup exceeds 256 MiB")
		}
		if h.Name == "backup-version" {
			if h.Typeflag != tar.TypeReg || h.Size != int64(len(backupManifest)) {
				return errors.New("unsupported backup version")
			}
			body, e := io.ReadAll(r)
			if e != nil {
				return e
			}
			valid = string(body) == backupManifest
			if !valid {
				return errors.New("unsupported backup version")
			}
			continue
		}
		part, name, ok := strings.Cut(h.Name, "/")
		dir, known := staged[part]
		if !ok || !known || !fs.ValidPath(name) || name == "." || strings.HasPrefix(name, "locks/") || name == "locks" {
			return fmt.Errorf("invalid backup path %q", h.Name)
		}
		target := filepath.Join(dir, filepath.FromSlash(name))
		if h.Typeflag == tar.TypeDir {
			if e := os.MkdirAll(target, 0700); e != nil {
				return e
			}
			continue
		}
		if e := os.MkdirAll(filepath.Dir(target), 0700); e != nil {
			return e
		}
		out, e := os.OpenFile(target, os.O_WRONLY|os.O_CREATE|os.O_EXCL, os.FileMode(h.Mode)&0777)
		if e != nil {
			return e
		}
		_, copyErr := io.CopyN(out, r, h.Size)
		if e := errors.Join(copyErr, out.Close()); e != nil {
			return e
		}
	}
	if !valid {
		return errors.New("missing backup version")
	}
	// Drain the gzip stream so a damaged checksum cannot publish a partial import.
	if n, e := io.Copy(io.Discard, io.LimitReader(z, (32<<20)+1)); e != nil {
		return e
	} else if n > 32<<20 {
		return errors.New("excess trailing backup data")
	}
	// Rename cannot replace a populated directory. Recheck empty targets immediately
	// before publication, leaving staged files unpublished on any collision.
	for _, dir := range []string{dataAbs, stateAbs} {
		if e := syscall.Rmdir(dir); e != nil && !errors.Is(e, fs.ErrNotExist) {
			return fmt.Errorf("destination changed during import: %w", e)
		}
	}
	if e := os.Rename(staged["data"], dataAbs); e != nil {
		return e
	}
	if e := os.Rename(staged["state"], stateAbs); e != nil {
		rollback := os.Rename(dataAbs, staged["data"])
		return errors.Join(e, rollback)
	}
	return nil
}
