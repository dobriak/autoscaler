package main

import (
	"fmt"
	"time"
)

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

//StartMonitor starts a ticker goroutine
func (a *App) StartMonitor() {
	tickers[a.AppID] = time.NewTicker(time.Second * time.Duration(a.Interval))
	go doMonitoring(a)
}

func doMonitoring(a *App) {
	internal := 0
	for t := range tickers[a.AppID].C {
		fmt.Printf("ticker:\t%s,%s\n", a.AppID, t)
		internal++
		monitor(a)
		fmt.Printf("app.Monitor:%s, internal:%d\n", a.AppID, internal)
	}
}

//StopMonitor stops the ticker associated with the given app
func (a *App) StopMonitor() {
	tickers[a.AppID].Stop()
}

func monitor(a *App) {
	//get all apps in marathon
	fmt.Println("GET /service/marathon/v2/apps")
	//check if a.appid is in marathon apps, wait if not
	fmt.Printf("GET /service/marathon/v2/apps/%s\n", a.AppID)
	fmt.Println("app:instances")
	//get all tasks for appid
	fmt.Println("app:tasks")
	fmt.Printf("\tid, slaveId\n")
	//for each task X on agent Y get performance stats
	fmt.Println("GET /slave/<slaveId>/monitor/statistics.json")
	fmt.Println("<tid>==executorId, executorId:statistics")
	//  get cpu usage
	fmt.Println("statistics:cpus_system_time_secs,cpus_user_time_secs,timestamp")
	//  get mem usage
	fmt.Println("statistics:mem_rss_bytes,mem_limit_bytes")
	//calculate average mem and cpu
	fmt.Println("calculate average mem and cpu")
	//based on criteria + method, generate (scale[up|down] by factor) signal
	fmt.Printf("generate_autoscale_signal(avg_cpu, avg_mem, %s) vs (%d, %d)\n",
		a.Method, a.MaxCPU, a.MinMem)
	//put signal in a channel for sequential operation
	fmt.Println("channel <- signal")
}
