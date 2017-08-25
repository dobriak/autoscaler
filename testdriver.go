package main

import (
	"fmt"
	"time"
)

func testdriver() {

	//Download dcos-ca.crt
	//err := downloadFile("dcos-ca.crt", "/ca/dcos-ca.crt")
	//if err != nil {
	//	panic(err)
	//}
	apps := client.GetAllMarathonApps()
	fmt.Println("All marathon apps:")
	fmt.Println(apps)
	app1 := client.GetMarathonApp("/test1")
	fmt.Println(app1)

	var stats1, stats2 TaskStats
	var cpu, cpu1, cpu2, cpuD, timeD float64
	var mem float64
	fmt.Println("========app1============")
	for _, task := range app1.App.Tasks {
		fmt.Printf("id:%s app_id:%s slave_id:%s\n", task.ID, task.AppID, task.SlaveID)
		fmt.Println("------stats-------")
		stats1 = client.GetTaskStats(task.ID, task.SlaveID)
		//fmt.Println(stats)
		time.Sleep(time.Second * 1)
		stats2 = client.GetTaskStats(task.ID, task.SlaveID)

		cpu1 = stats1.Statistics.CpusSystemTimeSecs + stats1.Statistics.CpusUserTimeSecs
		cpu2 = stats2.Statistics.CpusSystemTimeSecs + stats2.Statistics.CpusUserTimeSecs
		cpuD = cpu2 - cpu1
		timeD = stats2.Statistics.Timestamp - stats1.Statistics.Timestamp
		cpu = (cpuD / timeD) * 100
		mem = (stats1.Statistics.MemRssBytes / stats1.Statistics.MemLimitBytes) * 100
		fmt.Printf("cpu:%f, mem:%f\n", cpu, mem)
	}

	//app3 := client.GetMarathonApp("/folder/test3")
	//fmt.Println(app3)

	fmt.Println("Done")
}
