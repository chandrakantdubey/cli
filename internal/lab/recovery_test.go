package lab

import (
	"archive/tar"
	"compress/gzip"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"cli-silly/course"
)

func TestBackupRoundTripAndCollision(t *testing.T) {
	l, err := course.Load(course.First)
	if err != nil {
		t.Fatal(err)
	}
	s := Store{DataDir: t.TempDir(), StateDir: t.TempDir()}
	dir, err := s.Open(l)
	if err != nil {
		t.Fatal(err)
	}
	if err := os.Chmod(filepath.Join(dir, "bin/health-check"), 0744); err != nil {
		t.Fatal(err)
	}
	if _, err := s.Check(l, dir); err != nil {
		t.Fatal(err)
	}
	archive := filepath.Join(t.TempDir(), "work.tar.gz")
	if err := os.Mkdir(filepath.Join(dir, "empty-folder"), 0700); err != nil {
		t.Fatal(err)
	}
	if err := s.ExportBackup(archive); err != nil {
		t.Fatal(err)
	}
	if err := s.ExportBackup(archive); err == nil {
		t.Fatal("overwrote backup")
	}
	target := Store{DataDir: filepath.Join(t.TempDir(), "data"), StateDir: filepath.Join(t.TempDir(), "state")}
	if err := target.ImportBackup(archive); err != nil {
		t.Fatal(err)
	}
	if info, err := os.Stat(filepath.Join(target.Path(l), "empty-folder")); err != nil || !info.IsDir() {
		t.Fatalf("empty directory lost: %v", err)
	}
	p, err := target.Progress(l)
	if err != nil || p.Best != 40 {
		t.Fatalf("lost score %+v %v", p, err)
	}
	info, err := os.Stat(filepath.Join(target.Path(l), "bin/health-check"))
	if err != nil || info.Mode().Perm() != 0744 {
		t.Fatalf("lost permissions %v %v", info, err)
	}
	if err := target.ImportBackup(archive); err == nil {
		t.Fatal("overwrote learner data")
	}
	if err := os.Symlink("/etc/passwd", filepath.Join(dir, "link")); err != nil {
		t.Fatal(err)
	}
	if err := s.ExportBackup(filepath.Join(t.TempDir(), "unsafe.tar.gz")); err == nil {
		t.Fatal("followed symlink")
	}
}

func TestImportRejectsEscapingPathsAndLinks(t *testing.T) {
	for _, name := range []string{"data/../../escape", "/absolute", "other/file", "data/link"} {
		t.Run(name, func(t *testing.T) {
			file, err := os.Create(filepath.Join(t.TempDir(), "bad.tar.gz"))
			if err != nil {
				t.Fatal(err)
			}
			z := gzip.NewWriter(file)
			tw := tar.NewWriter(z)
			h := &tar.Header{Name: name, Mode: 0600, Size: 0}
			if name == "data/link" {
				h.Typeflag = tar.TypeSymlink
				h.Linkname = "/tmp"
			}
			if err := tw.WriteHeader(h); err != nil {
				t.Fatal(err)
			}
			tw.Close()
			z.Close()
			file.Close()
			s := Store{DataDir: filepath.Join(t.TempDir(), "data"), StateDir: filepath.Join(t.TempDir(), "state")}
			if err := s.ImportBackup(file.Name()); err == nil {
				t.Fatal("accepted invalid archive")
			}
			if _, err := os.Stat(s.DataDir); !os.IsNotExist(err) {
				t.Fatal("published invalid archive")
			}
		})
	}
}

func TestRestorePreservesBothAttemptsAndInspection(t *testing.T) {
	l, err := course.Load(course.First)
	if err != nil {
		t.Fatal(err)
	}
	s := Store{DataDir: t.TempDir(), StateDir: t.TempDir()}
	dir, err := s.Open(l)
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir, "my work"), []byte("old\n"), 0600); err != nil {
		t.Fatal(err)
	}
	archive, err := s.Reset(l)
	if err != nil {
		t.Fatal(err)
	}
	paths, err := s.Archives(l)
	if err != nil || len(paths) != 1 || paths[0] != archive {
		t.Fatalf("%v %v", paths, err)
	}
	if err := os.WriteFile(filepath.Join(dir, "fresh"), []byte("new\n"), 0600); err != nil {
		t.Fatal(err)
	}
	backup, err := s.Restore(l, filepath.Base(filepath.Dir(archive)))
	if err != nil {
		t.Fatal(err)
	}
	if _, err := os.Stat(filepath.Join(backup, "fresh")); err != nil {
		t.Fatal(err)
	}
	body, err := s.Inspect(l)
	if err != nil || !strings.Contains(body, `"my work"`) {
		t.Fatalf("%s %v", body, err)
	}
	if _, err := s.Restore(l, "../escape"); err == nil {
		t.Fatal("accepted escaping archive")
	}
}

func TestInspectionRejectsRedirectedGitDirectory(t *testing.T) {
	l, err := course.Load("git/selective-staging")
	if err != nil {
		t.Fatal(err)
	}
	s := Store{DataDir: t.TempDir(), StateDir: t.TempDir()}
	dir, err := s.Open(l)
	if err != nil {
		t.Fatal(err)
	}
	body, err := s.Inspect(l)
	if err != nil || !strings.Contains(body, "Git history") {
		t.Fatalf("%s %v", body, err)
	}
	gitDir := filepath.Join(dir, ".git")
	saved := filepath.Join(t.TempDir(), "saved-git")
	if err := os.Rename(gitDir, saved); err != nil {
		t.Fatal(err)
	}
	if err := os.Symlink(saved, gitDir); err != nil {
		t.Fatal(err)
	}
	if _, err := s.Inspect(l); err == nil || !strings.Contains(err.Error(), "own .git directory") {
		t.Fatalf("inspected another repository: %v", err)
	}
}
