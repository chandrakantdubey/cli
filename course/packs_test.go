package course_test

import (
	"archive/tar"
	"compress/gzip"
	"encoding/json"
	"io/fs"
	"os"
	"path/filepath"
	"testing"

	"cli-silly/course"
	"cli-silly/internal/lab"
)

func TestCoursePackUpgradePreservesGradesAndWorkspace(t *testing.T) {
	store := lab.Store{DataDir: t.TempDir(), StateDir: t.TempDir()}
	library := course.Library{PackDir: filepath.Join(store.DataDir, "courses")}
	source := filepath.Join("..", "examples", "markdown-basics")
	pack, err := library.Install(source)
	if err != nil || pack.Version != 1 {
		t.Fatalf("install example pack: %+v (%v)", pack, err)
	}
	if _, err := library.Install(source); err == nil {
		t.Fatal("same version replaced existing installation")
	}
	id := "markdown-basics/project-heading"
	lesson, err := library.Load(id)
	if err != nil {
		t.Fatal(err)
	}
	path, err := store.Open(lesson)
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(path, "README.md"), []byte("# Team Notes\n"), 0644); err != nil {
		t.Fatal(err)
	}
	result, err := store.Check(lesson, path)
	if err != nil || !result.Complete {
		t.Fatalf("grading installed lesson: %+v (%v)", result, err)
	}
	copy := filepath.Join(t.TempDir(), "upgrade")
	if err := os.CopyFS(copy, os.DirFS(source)); err != nil {
		t.Fatal(err)
	}
	changeManifest := func(version, schema int) {
		t.Helper()
		pack.Version, pack.EngineSchema = version, schema
		data, err := json.Marshal(pack)
		if err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(filepath.Join(copy, "pack.json"), data, 0644); err != nil {
			t.Fatal(err)
		}
	}
	changeManifest(2, 1)
	if _, err := library.Install(copy); err != nil {
		t.Fatal(err)
	}
	changeManifest(3, 99)
	if _, err := library.Install(copy); err == nil {
		t.Fatal("incompatible pack upgrade succeeded")
	}
	packs, err := library.Packs()
	if err != nil || len(packs) != 1 || packs[0].Version != 2 {
		t.Fatalf("failed upgrade invalidated installed pack: %+v (%v)", packs, err)
	}
	newLesson, err := library.Load(id)
	if err != nil {
		t.Fatal(err)
	}
	newPath, err := store.Open(newLesson)
	if err != nil || newPath != path {
		t.Fatalf("prior attempt was lost: %q (%v)", newPath, err)
	}
	progress, err := store.Progress(newLesson)
	if err != nil || progress.Best != 100 {
		t.Fatalf("grade was lost: %+v (%v)", progress, err)
	}
	if _, err := library.Load(course.First); err != nil {
		t.Fatalf("built-in lesson unavailable: %v", err)
	}
}

func TestCoursePackRejectsTraversalAndSymlinks(t *testing.T) {
	root := t.TempDir()
	library := course.Library{PackDir: filepath.Join(root, "courses")}
	for _, entry := range []tar.Header{
		{Name: "../escape", Mode: 0644, Typeflag: tar.TypeReg},
		{Name: "/absolute", Mode: 0644, Typeflag: tar.TypeReg},
		{Name: "redirect", Linkname: "../escape", Typeflag: tar.TypeSymlink},
	} {
		filename := filepath.Join(t.TempDir(), "invalid.tar.gz")
		file, err := os.Create(filename)
		if err != nil {
			t.Fatal(err)
		}
		gz := gzip.NewWriter(file)
		writer := tar.NewWriter(gz)
		if err := writer.WriteHeader(&entry); err != nil {
			t.Fatal(err)
		}
		if err := writer.Close(); err != nil {
			t.Fatal(err)
		}
		if err := gz.Close(); err != nil {
			t.Fatal(err)
		}
		if err := file.Close(); err != nil {
			t.Fatal(err)
		}
		if _, err := library.Install(filename); err == nil {
			t.Fatalf("invalid tar path %q was accepted", entry.Name)
		}
	}
	if _, err := os.Stat(filepath.Join(root, "escape")); !os.IsNotExist(err) {
		t.Fatalf("pack touched parent path: %v", err)
	}
}

func TestCoursePackRejectsInvalidMetadata(t *testing.T) {
	source := filepath.Join("..", "examples", "markdown-basics")
	library := course.Library{PackDir: t.TempDir()}
	for _, manifest := range []string{
		`{"schemaVersion":1,"engineSchema":1,"version":1,"id":"shell","title":"Collision","lessons":["project-heading"]}`,
		`{"schemaVersion":1,"engineSchema":1,"version":1,"id":"markdown-basics","title":"\u001b[31mInjected","lessons":["project-heading"]}`,
		`{"schemaVersion":1,"engineSchema":1,"version":1,"id":"markdown-basics","title":"Duplicate","lessons":["project-heading","project-heading"]}`,
	} {
		directory := filepath.Join(t.TempDir(), "pack")
		if err := os.CopyFS(directory, os.DirFS(source)); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(filepath.Join(directory, "pack.json"), []byte(manifest), fs.FileMode(0644)); err != nil {
			t.Fatal(err)
		}
		if _, err := library.Install(directory); err == nil {
			t.Fatalf("invalid metadata was accepted: %s", manifest)
		}
	}
}
