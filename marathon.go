package main

//DcosBasicAuth struct
type DcosBasicAuth struct {
	UID      string `json:"uid"`
	Password string `json:"password"`
}

//DcosAuthResponse struct
type DcosAuthResponse struct {
	Token string `json:"token"`
}

//MarathonAppInstances for scaling up and down
type MarathonAppInstances struct {
	Instances int `json:"instances"`
}

//MarathonScaleResult response you get trying to scale an app
type MarathonScaleResult struct {
	Version      string `json:"version"`
	DeploymentID string `json:"deploymentId"`
}

//MarathonApps struct
type MarathonApps struct {
	Apps []struct {
		ID        string `json:"id"`
		Instances int    `json:"instances"`
	} `json:"apps"`
}

//MarathonApp struct
type MarathonApp struct {
	App struct {
		ID        string `json:"id"`
		Instances int    `json:"instances"`
		Tasks     []struct {
			ID      string `json:"id"`
			State   string `json:"state"`
			AppID   string `json:"appId"`
			SlaveID string `json:"slaveId"`
		} `json:"tasks"`
	} `json:"app"`
}

//TaskStats struct
type TaskStats struct {
	ExecutorID  string `json:"executor_id"`
	FrameworkID string `json:"framework_id"`
	Source      string `json:"source"`
	Statistics  struct {
		CpusSystemTimeSecs float64 `json:"cpus_system_time_secs"`
		CpusUserTimeSecs   float64 `json:"cpus_user_time_secs"`
		MemRssBytes        float64 `json:"mem_rss_bytes"`
		MemLimitBytes      float64 `json:"mem_limit_bytes"`
		Timestamp          float64 `json:"timestamp"`
	} `json:"statistics"`
}

//FilterNonRunningTasks reslices the Tasks keeping only running ones
func (m *MarathonApp) FilterNonRunningTasks() {
	k := 0
	for i, task := range m.App.Tasks {
		if task.State == "TASK_RUNNING" {
			if i != k {
				m.App.Tasks[k] = task
			}
			k++
		}
	}
	m.App.Tasks = m.App.Tasks[:k]
}
