package main

import (
	"time"

	log "github.com/Sirupsen/logrus"
)

//App struct describing autoscaling app
type App struct {
	AppID        string  `json:"app_id"`
	MaxCPU       float64 `json:"max_cpu"`
	MinCPU       float64 `json:"min_cpu"`
	MaxMem       float64 `json:"max_mem"`
	MinMem       float64 `json:"min_mem"`
	Method       string  `json:"method"`
	ScaleFactor  int     `json:"scale_factor"`
	MaxInstances int     `json:"max_instances"`
	MinInstances int     `json:"min_instances"`
	WarmUp       int     `json:"warm_up"`
	CoolDown     int     `json:"cool_down"`
	Interval     int     `json:"interval"`
}

//Apps - all monitored apps
type Apps []App

type appState struct {
	warmUp   int
	coolDown int
}

//StartMonitor starts a ticker goroutine
func (a *App) StartMonitor() {
	tickers[a.AppID] = time.NewTicker(time.Second * time.Duration(a.Interval))
	go a.doMonitor()
}

//doMonitor will be storing the intermediate state of the app metrics
func (a *App) doMonitor() {
	as := appState{0, 0}
	var cpu, mem float64
	for range tickers[a.AppID].C {
		if !client.AppExists(a) {
			//fmt.Printf("%s not found in /service/marathon/v2/apps\n", a.AppID)
			log.Warningf("%s not found in /service/marathon/v2/apps\n", a.AppID)
			continue
		}
		marathonApp := client.GetMarathonApp(a.AppID)
		if marathonApp.App.Instances == 0 {
			//fmt.Printf("%s suspended, skipping monitoring cycle\n", marathonApp.App.ID)
			log.Warningf("%s suspended, skipping monitoring cycle\n", marathonApp.App.ID)
			continue
		}
		if !a.EnsureMinMaxInstances(marathonApp) {
			continue
		}
		cpu, mem = a.getCPUMem(marathonApp)
		//fmt.Printf("app:%s cpu:%f, mem:%f\n", a.AppID, cpu, mem)
		log.Infof("app:%s cpu:%f, mem:%f\n", a.AppID, cpu, mem)
		a.AutoScale(cpu, mem, &as, marathonApp)
	}
}

//StopMonitor stops the ticker associated with the given app
func (a *App) StopMonitor() {
	tickers[a.AppID].Stop()
}

func (a *App) getCPUMem(marathonApp MarathonApp) (float64, float64) {
	var (
		stats1, stats2               TaskStats
		cpu, cpu1, cpu2, cpuD, timeD float64
		mem                          float64
	)
	marathonApp.FilterNonRunningTasks()
	for _, task := range marathonApp.App.Tasks {
		//fmt.Printf("id:%s app_id:%s slave_id:%s\n", task.ID, task.AppID, task.SlaveID)
		stats1 = client.GetTaskStats(task.ID, task.SlaveID)
		//fmt.Println(stats)
		time.Sleep(time.Second * 1)
		stats2 = client.GetTaskStats(task.ID, task.SlaveID)

		cpu1 = stats1.Statistics.CpusSystemTimeSecs + stats1.Statistics.CpusUserTimeSecs
		cpu2 = stats2.Statistics.CpusSystemTimeSecs + stats2.Statistics.CpusUserTimeSecs
		cpuD = cpu2 - cpu1
		timeD = stats2.Statistics.Timestamp - stats1.Statistics.Timestamp
		cpu = cpu + (cpuD / timeD)
		mem = mem + (stats1.Statistics.MemRssBytes / stats1.Statistics.MemLimitBytes)
	}
	cpu = cpu / float64(len(marathonApp.App.Tasks)) * 100
	mem = mem / float64(len(marathonApp.App.Tasks)) * 100
	return cpu, mem
}
