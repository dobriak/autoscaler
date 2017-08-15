package main

import (
	"fmt"
	"time"
)

var apps Apps
var tickers map[string]*time.Ticker

func init() {
	tickers = make(map[string]*time.Ticker)
	RepoAddApp(App{"test1", 10, 20, "cpu", 5})
	RepoAddApp(App{"test2", 20, 10, "mem", 7})
}

//RepoAddApp adds an App to the repo
func RepoAddApp(a App) {
	if !RepoAppInApps(a.AppID) {
		apps = append(apps, a)
		//start ticker
		tickers[a.AppID] = time.NewTicker(time.Second * time.Duration(a.Interval))
		go func(myapp App) {
			for t := range tickers[myapp.AppID].C {
				fmt.Printf("ticker:\t%s\t", myapp.AppID)
				fmt.Println(t)
			}
		}(a)
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
