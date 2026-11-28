package lab

import (
	"fmt"
	"slices"
	"strings"
)

// Report differences without echoing learner-controlled terminal sequences or
// revealing an entire answer. The task and authored hint supply the next step.
func textDifference(actual []byte, expected string) string {
	value := string(actual)
	if len(actual) == 0 {
		return "The report is empty. Check the command's output and stderr before redirecting it."
	}
	if strings.ReplaceAll(value, "\r\n", "\n") == expected {
		return "The content matches, but uses CRLF line endings. This task expects Linux LF line endings."
	}
	if strings.TrimSuffix(value, "\n") == strings.TrimSuffix(expected, "\n") {
		return "The content matches except for the final newline. Use printf with a newline or preserve the original output."
	}
	a, b := strings.Split(strings.TrimSuffix(value, "\n"), "\n"), strings.Split(strings.TrimSuffix(expected, "\n"), "\n")
	if len(a) != len(b) {
		return fmt.Sprintf("Expected %d lines; found %d. Inspect the report for missing or extra lines.", len(b), len(a))
	}
	x, y := slices.Clone(a), slices.Clone(b)
	slices.Sort(x)
	slices.Sort(y)
	if slices.Equal(x, y) {
		return "The lines are correct but in a different order. Check the ordering required by the task."
	}
	for i := range a {
		if a[i] != b[i] {
			return fmt.Sprintf("Line %d differs from the required report. Inspect that line and the source data.", i+1)
		}
	}
	return "Inspect the report's formatting."
}
