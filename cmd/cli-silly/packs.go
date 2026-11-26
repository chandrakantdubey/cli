package main

import (
	"errors"
	"fmt"
	"path/filepath"

	"cli-silly/course"
	"cli-silly/internal/lab"
)

func courseLibrary(store lab.Store) course.Library {
	return course.Library{PackDir: filepath.Join(store.DataDir, "courses")}
}

func runCourse(store lab.Store, library course.Library, args []string) error {
	if len(args) == 1 && args[0] == "list" {
		packs, err := library.Packs()
		if err != nil {
			return err
		}
		for _, pack := range packs {
			fmt.Printf("%s  %s  (version %d)\n", pack.ID, pack.Title, pack.Version)
		}
		if len(packs) == 0 {
			fmt.Println("The built-in course is ready. No additional course packs installed.")
		}
		return nil
	}
	if len(args) == 2 && args[0] == "use" {
		id := args[1]
		if id == "builtin" {
			return store.SelectCourse("")
		}
		packs, err := library.Packs()
		if err != nil {
			return err
		}
		for _, pack := range packs {
			if pack.ID == id {
				if err := store.SelectCourse(id); err != nil {
					return err
				}
				fmt.Println("Selected", pack.Title)
				return nil
			}
		}
		return fmt.Errorf("course %q is not installed", id)
	}
	if len(args) == 2 && args[0] == "install" {
		pack, err := library.Install(args[1])
		if err != nil {
			return err
		}
		fmt.Printf("Installed %s (version %d). Open cli-silly or run cli-silly learn %s/%s.\n",
			pack.Title, pack.Version, pack.ID, pack.Lessons[0])
		return nil
	}
	return errors.New("use cli-silly course list, cli-silly course install PATH, or cli-silly course use ID")
}
