package main

import (
	"net/http"
)

//Route struct describing a router route
type Route struct {
	Name        string
	Method      string
	Pattern     string
	HandlerFunc http.HandlerFunc
}

//Routes - array of Route
type Routes []Route

var routes = Routes{
	Route{
		"Index",
		"GET",
		"/",
		Index,
	},
	Route{
		"IndexApps",
		"GET",
		"/apps",
		IndexApps,
	},
	Route{
		"GetApp",
		"GET",
		"/apps/{appid}",
		GetApp,
	},
	Route{
		"AddApp",
		"POST",
		"/apps",
		AddApp,
	},
	Route{
		"RemoveApp",
		"DELETE",
		"/apps/{appid}",
		RemoveApp,
	},
}
