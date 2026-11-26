package course

import (
	"embed"
	"encoding/json"
	"fmt"
	"io/fs"
	"path"
)

//go:embed catalog.json references.json all:workshops
var assets embed.FS

type Frame struct {
	Permissions string `json:"permissions,omitempty"`
	OutputFile  string `json:"outputFile,omitempty"`
	Command     string `json:"command"`
	Visual      string `json:"visual"`
	Explanation string `json:"explanation"`
}

type Criterion struct {
	ID       string `json:"id"`
	Kind     string `json:"kind"`
	Path     string `json:"path"`
	Expected string `json:"expected,omitempty"`
	Member   string `json:"member,omitempty"`
	Points   int    `json:"points"`
	Label    string `json:"label"`
	Feedback string `json:"feedback,omitempty"`
}

type Question struct {
	ID          string   `json:"id"`
	Prompt      string   `json:"prompt"`
	Choices     []string `json:"choices"`
	Answer      int      `json:"answer"`
	Explanation string   `json:"explanation"`
}

type Lesson struct {
	ReviewOf        string      `json:"reviewOf,omitempty"`
	Quiz            []Question  `json:"quiz,omitempty"`
	ID              string      `json:"id"`
	Title           string      `json:"title"`
	StarterRevision int         `json:"starterRevision"`
	GradingRevision int         `json:"gradingRevision"`
	Tools           []string    `json:"tools,omitempty"`
	Setup           string      `json:"setup,omitempty"`
	Task            string      `json:"task"`
	Hints           []string    `json:"hints"`
	Checks          []Criterion `json:"checks"`
	Demo            []Frame     `json:"demo"`
	Body            string      `json:"-"`
	Starter         fs.FS       `json:"-"`
}

const First = "shell/make-it-runnable"

type Workshop struct {
	ID      string   `json:"id"`
	Title   string   `json:"title"`
	Lessons []string `json:"lessons"`
}

func Catalog() ([]Workshop, error) {
	data, err := assets.ReadFile("catalog.json")
	if err != nil {
		return nil, err
	}
	var catalog struct {
		SchemaVersion int        `json:"schemaVersion"`
		Workshops     []Workshop `json:"workshops"`
	}
	if err := json.Unmarshal(data, &catalog); err != nil {
		return nil, err
	}
	if catalog.SchemaVersion != 1 {
		return nil, fmt.Errorf("unsupported catalog version %d", catalog.SchemaVersion)
	}
	return catalog.Workshops, nil
}

func Load(id string) (Lesson, error) {
	workshops, err := Catalog()
	if err != nil {
		return Lesson{}, err
	}
	found := false
	for _, workshop := range workshops {
		for _, lessonID := range workshop.Lessons {
			if lessonID == id {
				found = true
			}
		}
	}
	if !found || !fs.ValidPath(id) {
		return Lesson{}, fmt.Errorf("unknown lesson %q (try cli-silly list)", id)
	}
	dir := path.Join("workshops", id)
	data, err := assets.ReadFile(path.Join(dir, "lesson.json"))
	if err != nil {
		return Lesson{}, err
	}
	var lesson Lesson
	if err := json.Unmarshal(data, &lesson); err != nil {
		return Lesson{}, err
	}
	body, err := assets.ReadFile(path.Join(dir, "lesson.md"))
	if err != nil {
		return Lesson{}, err
	}
	lesson.Body = string(body)
	lesson.Starter, err = fs.Sub(assets, path.Join(dir, "starter"))
	if err != nil {
		return Lesson{}, err
	}

	total := 0
	seen := make(map[string]bool)
	for _, check := range lesson.Checks {
		if check.ID == "" || seen[check.ID] || check.Points < 1 || !fs.ValidPath(check.Path) || check.Path == "." {
			return Lesson{}, fmt.Errorf("invalid criterion %q in %s", check.ID, id)
		}
		switch check.Kind {
		case "owner-executable", "text", "absent", "git-count", "git-tree", "git-index", "git-branch",
			"docker-running", "docker-healthy", "docker-recreated", "http-body", "git-ancestor", "git-file", "archive-member", "script-config", "script-report":
		default:
			return Lesson{}, fmt.Errorf("unknown criterion type %q", check.Kind)
		}
		seen[check.ID] = true
		total += check.Points
	}
	if total != 100 || lesson.Title == "" || lesson.Task == "" || len(lesson.Hints) != 3 || len(lesson.Demo) == 0 || lesson.ID != id || lesson.StarterRevision < 1 || lesson.GradingRevision < 1 {
		return Lesson{}, fmt.Errorf("invalid lesson manifest for %s", id)
	}
	questions := make(map[string]bool)
	for _, question := range lesson.Quiz {
		if question.ID == "" || questions[question.ID] || question.Prompt == "" || question.Explanation == "" ||
			len(question.Choices) < 2 || len(question.Choices) > 4 || question.Answer < 1 || question.Answer > len(question.Choices) {
			return Lesson{}, fmt.Errorf("invalid quiz question %q in %s", question.ID, id)
		}
		questions[question.ID] = true
	}
	return lesson, nil
}

type Entry struct {
	Workshop string
	Lesson   Lesson
}

func Entries() ([]Entry, error) {
	workshops, err := Catalog()
	if err != nil {
		return nil, err
	}
	var entries []Entry
	seen := make(map[string]bool)
	for _, workshop := range workshops {
		for _, id := range workshop.Lessons {
			if seen[id] {
				return nil, fmt.Errorf("duplicate lesson %q", id)
			}
			lesson, err := Load(id)
			if err != nil {
				return nil, err
			}
			seen[id] = true
			entries = append(entries, Entry{Workshop: workshop.Title, Lesson: lesson})
		}
	}
	for _, entry := range entries {
		if entry.Lesson.ReviewOf == "" {
			continue
		}
		base, err := Load(entry.Lesson.ReviewOf)
		if err != nil || base.ReviewOf != "" {
			return nil, fmt.Errorf("invalid review base for %s", entry.Lesson.ID)
		}
	}
	return entries, nil
}
