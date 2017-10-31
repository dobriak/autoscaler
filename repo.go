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
