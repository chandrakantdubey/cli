package main

import (
	"errors"
	"fmt"
	"slices"
	"strings"

	"cli-silly/course"
	"cli-silly/internal/lab"
)

func runReference(command string, args []string) error {
	references, err := course.References()
	if err != nil {
		return err
	}
	if command == "snippets" {
		for _, entry := range references {
			if entry.Matches(strings.Join(args, " ")) {
				fmt.Printf("%s  %s  (%s)\n", entry.ID, entry.Title, entry.Kind)
			}
		}
		return nil
	}
	if len(args) == 0 {
		return errors.New("provide a reference ID from cli-silly snippets")
	}
	id, args := args[0], args[1:]
	index := slices.IndexFunc(references, func(entry course.Reference) bool { return entry.ID == id })
	if index < 0 {
		return fmt.Errorf("unknown reference %q", id)
	}
	entry := references[index]
	if command == "snippet" {
		body, err := entry.Format(args)
		if err == nil {
			fmt.Print(body)
		}
		return err
	}
	if entry.Kind != "template" || len(args) == 0 {
		return errors.New("use cli-silly template ID PATH with an authored template and a workspace-relative path")
	}
	body, err := entry.Format(args[1:])
	if err != nil {
		return err
	}
	store, err := localStore()
	if err != nil {
		return err
	}
	dir, id, err := lab.FindWorkspace()
	if err != nil {
		return err
	}
	lesson, err := courseLibrary(store).Load(id)
	if err != nil {
		return err
	}
	if err := store.WriteTemplate(lesson, dir, args[0], body); err != nil {
		return err
	}
	fmt.Println("Created", args[0])
	return nil
}
