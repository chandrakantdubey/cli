package course

import "testing"

func TestProjectsAndFoundationsHavePractice(t *testing.T) {
	projects, err := Projects()
	if err != nil || len(projects) < 5 {
		t.Fatalf("projects %v: %v", projects, err)
	}
	entries, err := Entries()
	if err != nil {
		t.Fatal(err)
	}
	for _, e := range entries {
		if e.Workshop != "Foundations" {
			continue
		}
		found := false
		for _, variant := range entries {
			if variant.Lesson.ReviewOf == e.Lesson.ID {
				found = true
			}
		}
		if !found {
			t.Errorf("no follow-up for %s", e.Lesson.ID)
		}
	}
}
