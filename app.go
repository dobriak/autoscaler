package main

type App struct {
	AppId  string `json:"app_id"`
	MaxCpu int    `json:"max_cpu"`
	MinMem int    `json:"min_mem"`
	Method string `json:"method"`
}

type Apps []App
