package main

//App struct describing autoscaling app
type App struct {
	AppID    string `json:"app_id"`
	MaxCPU   int    `json:"max_cpu"`
	MinMem   int    `json:"min_mem"`
	Method   string `json:"method"`
	Interval int    `json:"interval"`
}

//Apps - all monitored apps
type Apps []App
