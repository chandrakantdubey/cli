package course

import (
	"encoding/json"
	"fmt"
	"regexp"
	"slices"
	"strings"
)

type Parameter struct {
	Name    string `json:"name"`
	Default string `json:"default"`
}

type Reference struct {
	ID         string      `json:"id"`
	Title      string      `json:"title"`
	Tools      []string    `json:"tools"`
	LessonID   string      `json:"lessonId"`
	Kind       string      `json:"kind"`
	Parameters []Parameter `json:"parameters"`
	Body       string      `json:"body"`
	Example    string      `json:"example"`
	Expected   string      `json:"expected"`
}

var referenceID = regexp.MustCompile(`^[a-z][a-z0-9_-]*$`)
var parameterToken = regexp.MustCompile(`\{\{([a-z][a-z0-9_-]*)\}\}`)

func References() ([]Reference, error) {
	data, err := assets.ReadFile("references.json")
	if err != nil {
		return nil, err
	}
	var references []Reference
	if err := json.Unmarshal(data, &references); err != nil {
		return nil, err
	}
	seen := make(map[string]bool)
	for _, entry := range references {
		if !referenceID.MatchString(entry.ID) || seen[entry.ID] || entry.Title == "" || entry.Body == "" ||
			(entry.Kind != "snippet" && entry.Kind != "template") {
			return nil, fmt.Errorf("invalid reference %q", entry.ID)
		}
		seen[entry.ID] = true
		parameters := make(map[string]bool)
		for _, parameter := range entry.Parameters {
			if !referenceID.MatchString(parameter.Name) || parameters[parameter.Name] {
				return nil, fmt.Errorf("invalid parameter in reference %s", entry.ID)
			}
			parameters[parameter.Name] = true
		}
		if _, err := entry.Format(nil); err != nil {
			return nil, err
		}
		if _, err := Load(entry.LessonID); err != nil {
			return nil, err
		}
	}
	return references, nil
}

func (r Reference) Matches(query string) bool {
	haystack := strings.ToLower(r.ID + " " + r.Title + " " + r.Body + " " + strings.Join(r.Tools, " "))
	return !slices.ContainsFunc(strings.Fields(strings.ToLower(query)), func(term string) bool {
		return !strings.Contains(haystack, term)
	})
}

// Format substitutes shell-quoted parameters so each value remains one argument.
func (r Reference) Format(args []string) (string, error) {
	values := make(map[string]string)
	for _, parameter := range r.Parameters {
		values[parameter.Name] = parameter.Default
	}
	explicit := make(map[string]bool)
	for _, arg := range args {
		name, value, ok := strings.Cut(arg, "=")
		if _, known := values[name]; !known || !ok || explicit[name] || strings.ContainsAny(value, "\x00\r\n") {
			return "", fmt.Errorf("invalid or duplicate reference parameter %q", name)
		}
		explicit[name] = true
		values[name] = value
	}
	var err error
	body := parameterToken.ReplaceAllStringFunc(r.Body, func(token string) string {
		name := token[2 : len(token)-2]
		value, ok := values[name]
		if !ok {
			err = fmt.Errorf("reference %s uses unknown parameter %q", r.ID, name)
		}
		return "'" + strings.ReplaceAll(value, "'", "'\\''") + "'"
	})
	return body + "\n", err
}
