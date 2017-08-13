package main

import "fmt"

var apps Apps

func init() {
	RepoAddApp(App{AppID: "test1", MaxCPU: 10, MinMem: 20, Method: "cpu"})
	RepoAddApp(App{AppID: "test2", MaxCPU: 20, MinMem: 10, Method: "mem"})
}

//RepoAddApp adds an App to the repo
func RepoAddApp(a App) {
	if !RepoAppInApps(a.AppID) {
		apps = append(apps, a)
	}
}

//RepoAppInApps finds if an app is present in the apps list
func RepoAppInApps(appID string) bool {
	for _, a := range apps {
		if a.AppID == appID {
			return true
		}
	}
	return false
}

//RepoFindApp returns an App object based on app ID
func RepoFindApp(appID string) App {
	for _, a := range apps {
		if a.AppID == appID {
			return a
		}
	}
	return App{}
}

//RepoRemoveApp re-slices the apps list to remove an app by its ID
func RepoRemoveApp(appID string) error {
	for i, a := range apps {
		if a.AppID == appID {
			apps = append(apps[:i], apps[i+1:]...)
			return nil
		}
	}
	return fmt.Errorf("could not find App with id of %s to delete", appID)
}
