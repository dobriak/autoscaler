package main

import (
	"fmt"
	"strings"
	"time"
)

var apps Apps
var tickers map[string]*time.Ticker

func init() {
	tickers = make(map[string]*time.Ticker)
	/*
		RepoAddApp(App{"/test1", 50.0, 30.0, 70.0, 10.0, "cpu", 2, 5, 2, 3, 3, 17})
		RepoAddApp(App{"/test2", 20.5, 10.5, 45.5, 11.5, "mem", 1, 7, 2, 3, 3, 21})
	*/
}

//RepoAddApp adds an App to the repo
func RepoAddApp(a App) {
	if !RepoAppInApps(a.AppID) {
		apps = append(apps, a)
		a.StartMonitor()
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
	appID = prependSlash(appID)
	for _, a := range apps {
		if a.AppID == appID {
			return a
		}
	}
	return App{}
}

//RepoRemoveApp re-slices the apps list to remove an app by its ID
func RepoRemoveApp(appID string) error {
	appID = prependSlash(appID)
	for i, a := range apps {
		if a.AppID == appID {
			apps = append(apps[:i], apps[i+1:]...)
			//Stopping the ticker
			tickers[appID].Stop()
			return nil
		}
	}
	return fmt.Errorf("could not find App with id of %s to delete", appID)
}

//RepoRemoveAllApps cycles through the apps array and removes them all
func RepoRemoveAllApps() error {
	for _, a := range apps {
		if err := RepoRemoveApp(a.AppID); err != nil {
			return err
		}
	}
	return nil
}

func prependSlash(appID string) string {
	if strings.Index(appID, "/") != 1 {
		appID = fmt.Sprintf("/%s", appID)
	}
	return appID
}
