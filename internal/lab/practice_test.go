package lab

import (
	"context"
	"os"
	"path/filepath"
	"sync"
	"testing"

	"cli-silly/course"
)

func TestEditorKeepsStarterPathLocal(t *testing.T) {
	t.Setenv("VISUAL", "")
	t.Setenv("EDITOR", "bash --noprofile")
	lesson, err := course.Load("shell/paths-and-quoting")
	if err != nil {
		t.Fatal(err)
	}
	store := Store{DataDir: t.TempDir(), StateDir: t.TempDir()}
	name := "Draft Files/server config.txt"
	cmd, err := store.Editor(lesson, name)
	if err != nil || cmd.Dir != store.Path(lesson) || len(cmd.Args) != 3 || cmd.Args[2] != "./"+name {
		t.Fatalf("editor: %+v (%v)", cmd, err)
	}
	absolute := filepath.Join(cmd.Dir, name)
	if err := os.Remove(absolute); err != nil {
		t.Fatal(err)
	}
	if err := os.Symlink(t.TempDir(), absolute); err != nil {
		t.Fatal(err)
	}
	if _, err := store.Editor(lesson, name); err == nil {
		t.Fatal("editor accepted a symlink")
	}
}

func TestWatchDeduplicatesAndStopsAtReset(t *testing.T) {
	lesson, err := course.Load(course.First)
	if err != nil {
		t.Fatal(err)
	}
	store := Store{DataDir: t.TempDir(), StateDir: t.TempDir()}
	path, err := store.Open(lesson)
	if err != nil {
		t.Fatal(err)
	}
	attempt, err := os.Stat(path)
	if err != nil {
		t.Fatal(err)
	}
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	var checks sync.WaitGroup
	checks.Go(func() {
		if _, err := store.Check(lesson, path); err != nil {
			t.Error(err)
		}
	})
	checks.Go(func() {
		for range 3 {
			if _, err := store.WatchCheck(ctx, lesson, attempt); err != nil {
				t.Error(err)
			}
		}
	})
	checks.Wait()
	before, err := store.Progress(lesson)
	if err != nil || before.Checks < 1 || before.Checks > 2 {
		t.Fatalf("initial check count: %+v (%v)", before, err)
	}
	if err := os.Chmod(filepath.Join(path, "bin", "health-check"), 0744); err != nil {
		t.Fatal(err)
	}
	for range 2 {
		result, err := store.WatchCheck(ctx, lesson, attempt)
		if err != nil || result.Score != 40 {
			t.Fatalf("automatic result: %+v (%v)", result, err)
		}
	}
	after, err := store.Progress(lesson)
	if err != nil || after.Checks != before.Checks+1 {
		t.Fatalf("duplicate automatic saves: %+v (%v)", after, err)
	}
	if _, err := store.Reset(lesson); err != nil {
		t.Fatal(err)
	}
	if _, err := store.WatchCheck(ctx, lesson, attempt); err == nil {
		t.Fatal("watch crossed a workspace reset")
	}
	cancel()
	if _, err := store.WatchCheck(ctx, lesson, nil); err == nil {
		t.Fatal("watch graded a cancelled session")
	}
}
