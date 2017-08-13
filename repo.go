package main

import "fmt"

var apps Apps

func init() {
	RepoAddApp(App{AppId: "test1", MaxCpu: 10, MinMem: 20, Method: "cpu"})
	RepoAddApp(App{AppId: "test2", MaxCpu: 20, MinMem: 10, Method: "mem"})
}

func RepoAddApp(a App) {
	if !RepoAppInApps(a.AppId) {
		apps = append(apps, a)
	}
}

func RepoAppInApps(app_id string) bool {
	for _, a := range apps {
		if a.AppId == app_id {
			return true
		}
	}
	return false
}

func RepoFindApp(app_id string) App {
	for _, a := range apps {
		if a.AppId == app_id {
			return a
		}
	}
	return App{}
}

func RepoRemoveApp(app_id string) error {
	for i, a := range apps {
		if a.AppId == app_id {
			apps = append(apps[:i], apps[i+1:]...)
			return nil
		}
	}
	return fmt.Errorf("Could not find App with id of %s to delete.", app_id)
}
