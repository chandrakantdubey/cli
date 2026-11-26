package lab

import (
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	"cli-silly/course"
)

func TestShellLab(t *testing.T) {
	lesson, err := course.Load(course.First)
	if err != nil {
		t.Fatal(err)
	}
	store := Store{DataDir: t.TempDir(), StateDir: t.TempDir()}
	path, err := store.Open(lesson)
	if err != nil {
		t.Fatal(err)
	}
	check := func(want int) {
		t.Helper()
		result, err := store.Check(lesson, path)
		if err != nil || result.Score != want || result.Complete != (want == 100) {
			t.Fatalf("check = %+v, %v; want %d", result, err, want)
		}
	}
	check(0)
	if err := os.Chmod(filepath.Join(path, "bin", "health-check"), 0744); err != nil {
		t.Fatal(err)
	}
	check(40)
	cmd := exec.Command("bash", "-c", "./bin/health-check > status.txt")
	cmd.Dir = path
	if output, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("run starter script: %v\n%s", err, output)
	}
	check(100)
	if _, err := store.Open(lesson); err != nil {
		t.Fatal(err)
	}
	check(100)
	backup, err := store.Reset(lesson)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := os.Stat(filepath.Join(backup, "status.txt")); err != nil {
		t.Fatal(err)
	}
	if _, err := store.Check(lesson, backup); err == nil {
		t.Fatal("archived workspace was accepted")
	}
	check(0)
	progress, err := store.Progress(lesson)
	if err != nil || progress.Best != 100 || progress.Latest.Score != 0 || progress.Checks != 5 {
		t.Fatalf("progress = %+v, %v", progress, err)
	}
}

func TestCheckDoesNotReadExternalReport(t *testing.T) {
	lesson, err := course.Load(course.First)
	if err != nil {
		t.Fatal(err)
	}
	store := Store{DataDir: t.TempDir(), StateDir: t.TempDir()}
	path, err := store.Open(lesson)
	if err != nil {
		t.Fatal(err)
	}
	external := filepath.Join(t.TempDir(), "report")
	if err := os.WriteFile(external, []byte("service: ready\n"), 0600); err != nil {
		t.Fatal(err)
	}
	if err := os.Symlink(external, filepath.Join(path, "status.txt")); err != nil {
		t.Fatal(err)
	}
	result, err := store.Check(lesson, path)
	if err != nil || result.Score != 0 {
		t.Fatalf("check = %+v, %v; want 0", result, err)
	}
}
