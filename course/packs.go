package course

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path"
	"path/filepath"
	"strconv"
	"strings"
	"syscall"
)

type Pack struct {
	SchemaVersion int      `json:"schemaVersion"`
	EngineSchema  int      `json:"engineSchema"`
	ID            string   `json:"id"`
	Title         string   `json:"title"`
	Version       int      `json:"version"`
	Lessons       []string `json:"lessons"`
}

type Library struct {
	PackDir string
}

func (l Library) Load(id string) (Lesson, error) {
	if lesson, err := Load(id); err == nil {
		return lesson, nil
	}
	packID, slug, ok := strings.Cut(id, "/")
	if !ok || !referenceID.MatchString(packID) || !referenceID.MatchString(slug) || l.PackDir == "" {
		return Lesson{}, fmt.Errorf("unknown lesson %q (try cli-silly list)", id)
	}
	pack, entries, err := l.installedPack(packID)
	if err == nil && pack.Version > 0 {
		for _, entry := range entries {
			if entry.Lesson.ID == id {
				return entry.Lesson, nil
			}
		}
	}
	if err != nil && !errors.Is(err, fs.ErrNotExist) {
		return Lesson{}, err
	}
	return Lesson{}, fmt.Errorf("unknown lesson %q (try cli-silly list)", id)
}

func (l Library) Entries() ([]Entry, error) {
	entries, err := Entries()
	if err != nil || l.PackDir == "" {
		return entries, err
	}
	items, err := os.ReadDir(l.PackDir)
	if errors.Is(err, fs.ErrNotExist) {
		return entries, nil
	}
	if err != nil {
		return nil, err
	}
	for _, item := range items {
		if !item.IsDir() || !referenceID.MatchString(item.Name()) {
			continue
		}
		_, installed, err := l.installedPack(item.Name())
		if errors.Is(err, fs.ErrNotExist) {
			continue
		}
		if err != nil {
			return nil, err
		}
		entries = append(entries, installed...)
	}
	return entries, nil
}

func (l Library) Packs() ([]Pack, error) {
	var packs []Pack
	if l.PackDir == "" {
		return packs, nil
	}
	items, err := os.ReadDir(l.PackDir)
	if errors.Is(err, fs.ErrNotExist) {
		return packs, nil
	}
	if err != nil {
		return nil, err
	}
	for _, item := range items {
		if !item.IsDir() || !referenceID.MatchString(item.Name()) {
			continue
		}
		pack, _, err := l.installedPack(item.Name())
		if errors.Is(err, fs.ErrNotExist) {
			continue
		}
		if err != nil {
			return nil, err
		}
		packs = append(packs, pack)
	}
	return packs, nil
}

func (l Library) installedPack(id string) (Pack, []Entry, error) {
	dir := filepath.Join(l.PackDir, id)
	data, err := os.ReadFile(filepath.Join(dir, "current"))
	if err != nil {
		return Pack{}, nil, err
	}
	version, err := strconv.Atoi(strings.TrimSpace(string(data)))
	if err != nil || version < 1 {
		return Pack{}, nil, fmt.Errorf("invalid installed version for %s", id)
	}
	pack, entries, err := validatePack(os.DirFS(filepath.Join(dir, strconv.Itoa(version))))
	if err != nil {
		return Pack{}, nil, fmt.Errorf("read installed course %s: %w", id, err)
	}
	if pack.ID != id || pack.Version != version {
		return Pack{}, nil, fmt.Errorf("installed course identity does not match %s", id)
	}
	return pack, entries, nil
}

func validatePack(source fs.FS) (Pack, []Entry, error) {
	data, err := fs.ReadFile(source, "pack.json")
	if err != nil {
		return Pack{}, nil, fmt.Errorf("read pack.json: %w", err)
	}
	var pack Pack
	if err := validatePackText(data, true); err != nil {
		return Pack{}, nil, err
	}
	if err := json.Unmarshal(data, &pack); err != nil {
		return Pack{}, nil, fmt.Errorf("parse pack.json: %w", err)
	}
	if pack.SchemaVersion != 1 || pack.EngineSchema != 1 || pack.Version < 1 ||
		(!referenceID.MatchString(pack.ID) || pack.ID == "builtin") || pack.Title == "" || len(pack.Lessons) == 0 {
		return Pack{}, nil, errors.New("invalid or incompatible course pack manifest")
	}
	builtIn, err := Catalog()
	if err != nil {
		return Pack{}, nil, err
	}
	for _, workshop := range builtIn {
		if workshop.ID == pack.ID {
			return Pack{}, nil, fmt.Errorf("course ID %q is reserved by a built-in workshop", pack.ID)
		}
	}
	var entries []Entry
	seen := make(map[string]bool)
	for _, slug := range pack.Lessons {
		if !referenceID.MatchString(slug) || seen[slug] {
			return Pack{}, nil, fmt.Errorf("invalid or duplicate lesson slug %q", slug)
		}
		seen[slug] = true
		dir := path.Join("lessons", slug)
		if info, err := fs.Stat(source, path.Join(dir, "starter")); err != nil || !info.IsDir() {
			return Pack{}, nil, fmt.Errorf("lesson %s has no starter directory", slug)
		}
		for _, file := range []string{"lesson.json", "lesson.md"} {
			data, err := fs.ReadFile(source, path.Join(dir, file))
			if err != nil {
				return Pack{}, nil, err
			}
			if err := validatePackText(data, file == "lesson.json"); err != nil {
				return Pack{}, nil, fmt.Errorf("%s/%s: %w", slug, file, err)
			}
		}
		lesson, err := loadLesson(source, dir, pack.ID+"/"+slug)
		if err != nil {
			return Pack{}, nil, err
		}
		if lesson.Setup != "" || lesson.ReviewOf != "" {
			return Pack{}, nil, fmt.Errorf("installed lesson %s must use declarative local file grading", lesson.ID)
		}
		for _, check := range lesson.Checks {
			if len(check.Preserve) > 0 {
				return Pack{}, nil, errors.New("preserve guards require a future course-pack engine schema; use separate text criteria in version 1 packs")
			}
			if check.Kind != "text" && check.Kind != "absent" && check.Kind != "owner-executable" {
				return Pack{}, nil, fmt.Errorf("checker %s is unavailable in installed packs", check.Kind)
			}
		}
		entries = append(entries, Entry{PackID: pack.ID, Workshop: pack.Title + " (" + pack.ID + ")", Lesson: lesson})
	}
	return pack, entries, nil
}

func (l Library) Install(source string) (Pack, error) {
	if l.PackDir == "" {
		return Pack{}, errors.New("no course pack directory configured")
	}
	if err := os.MkdirAll(l.PackDir, 0700); err != nil {
		return Pack{}, err
	}
	stage, err := os.MkdirTemp(l.PackDir, ".install-")
	if err != nil {
		return Pack{}, err
	}
	defer os.RemoveAll(stage)
	if err := stagePack(source, stage); err != nil {
		return Pack{}, err
	}
	pack, _, err := validatePack(os.DirFS(stage))
	if err != nil {
		return Pack{}, err
	}
	lock, err := os.OpenFile(filepath.Join(l.PackDir, ".lock"), os.O_CREATE|os.O_RDWR, 0600)
	if err != nil {
		return Pack{}, err
	}
	defer lock.Close()
	for {
		err = syscall.Flock(int(lock.Fd()), syscall.LOCK_EX)
		if !errors.Is(err, syscall.EINTR) {
			break
		}
	}
	if err != nil {
		return Pack{}, err
	}
	old, _, err := l.installedPack(pack.ID)
	if err != nil && !errors.Is(err, fs.ErrNotExist) {
		return Pack{}, err
	}
	if old.Version >= pack.Version {
		return Pack{}, fmt.Errorf("course %s version %d is already installed; choose a newer version", pack.ID, old.Version)
	}
	parent := filepath.Join(l.PackDir, pack.ID)
	if err := os.MkdirAll(parent, 0700); err != nil {
		return Pack{}, err
	}
	version := strconv.Itoa(pack.Version)
	final := filepath.Join(parent, version)
	pointer, err := os.CreateTemp(parent, ".current-")
	if err != nil {
		return Pack{}, err
	}
	defer os.Remove(pointer.Name())
	_, writeErr := pointer.WriteString(version + "\n")
	if err := errors.Join(writeErr, pointer.Close()); err != nil {
		return Pack{}, err
	}
	if err := os.Rename(stage, final); err != nil {
		return Pack{}, err
	}
	if err := os.Rename(pointer.Name(), filepath.Join(parent, "current")); err != nil {
		return Pack{}, errors.Join(err, os.RemoveAll(final))
	}
	return pack, nil
}
