package lab

import (
	"os"
	"path/filepath"
	"testing"

	"cli-silly/course"
)

func TestWriteTemplateScopesFilesAndPreservesExistingWork(t *testing.T) {
	lesson, err := course.Load(course.First)
	if err != nil {
		t.Fatal(err)
	}
	store := Store{DataDir: t.TempDir(), StateDir: t.TempDir()}
	path, err := store.Open(lesson)
	if err != nil {
		t.Fatal(err)
	}
	external := t.TempDir()
	if err := os.Symlink(external, filepath.Join(path, "out")); err != nil {
		t.Fatal(err)
	}
	for _, name := range []string{"../escape", "out/file", Marker, "/tmp/escape"} {
		if err := store.WriteTemplate(lesson, path, name, "example"); err == nil {
			t.Fatalf("successful write to invalid path %q", name)
		}
	}
	if err := store.WriteTemplate(lesson, path, "bin/my starter.sh", "first\n"); err != nil {
		t.Fatal(err)
	}
	if err := store.WriteTemplate(lesson, path, "bin/my starter.sh", "second\n"); err == nil {
		t.Fatal("template generator overwrote learner work")
	}
	data, err := os.ReadFile(filepath.Join(path, "bin/my starter.sh"))
	if err != nil || string(data) != "first\n" {
		t.Fatalf("read generated template: %q (%v)", data, err)
	}
}
