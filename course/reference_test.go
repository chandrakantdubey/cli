package course

import (
	"os/exec"
	"testing"
)

func TestReferenceShellQuotesParameters(t *testing.T) {
	references, err := References()
	if err != nil || len(references) < 2 {
		t.Fatalf("load references: %d, %v", len(references), err)
	}
	entry := Reference{ID: "sample", Parameters: []Parameter{{Name: "target"}}, Body: "printf '%s' {{target}}"}
	value := "it's one argument; $HOME"
	output, err := entry.Format([]string{"target=" + value})
	if err != nil {
		t.Fatal(err)
	}
	actual, err := exec.Command("bash", "-c", output).Output()
	if err != nil || string(actual) != value {
		t.Fatalf("rendered command = %q (%v)", actual, err)
	}
	for _, args := range [][]string{{"invalid=yes"}, {"target=a", "target=b"}, {"target=bad\nline"}} {
		if _, err := entry.Format(args); err == nil {
			t.Fatalf("invalid input accepted: %q", args)
		}
	}
}
