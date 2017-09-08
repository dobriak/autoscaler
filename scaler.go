package main

import (
	"fmt"
)

//ScaleSignal describes a scale proposal
type ScaleSignal struct {
	Scale scaleDirection
}

type scaleDirection struct {
	up   bool
	down bool
}

//generateSignal given cpu and mem values, return a scale proposal
func generateSignal(cpu, mem float64, a *App) ScaleSignal {
	result := ScaleSignal{}
	cpuDown := (cpu <= a.MinCPU)
	cpuUp := (cpu > a.MaxCPU)
	memDown := (mem <= a.MinMem)
	memUp := (mem > a.MinMem)
	switch method := a.Method; method {
	case "cpu":
		result.Scale.up = cpuUp
		result.Scale.down = cpuDown
	case "mem":
		result.Scale.up = memUp
		result.Scale.down = memDown
	case "and":
		result.Scale.up = cpuUp && memUp
		result.Scale.down = cpuDown && memDown
	case "or":
		result.Scale.up = cpuUp || memUp
		result.Scale.down = cpuDown || memDown
	default:
		fmt.Printf("method should be cpu|mem|and|or: %s\n", method)
		panic("Invalid scaling parameter method.")
	}
	if result.Scale.up && result.Scale.down {
		fmt.Printf("Scale up and scale down signal generated, defaulting to no operation. %+v\n", result)
		result.Scale.up = false
		result.Scale.down = false
	}

	return result
}

//AutoScale track and scale apps
func (a *App) AutoScale(cpu, mem float64, st *appState, mApp MarathonApp) {
	sig := generateSignal(cpu, mem, a)
	if !sig.Scale.down && !sig.Scale.up {
		st.coolDown = 0
		st.warmUp = 0
	} else {
		if sig.Scale.up {
			if mApp.App.Instances < a.MaxInstances {
				st.warmUp++
				if st.warmUp >= a.WarmUp {
					fmt.Printf("%s scale up triggered with %d of %d signals of %s\n",
						a.AppID, st.warmUp, a.WarmUp, a.Method)
					a.doScale(mApp, a.ScaleFactor)
					st.warmUp = 0
				} else {
					fmt.Printf("%s warming up %s(%d of %d)\n",
						a.AppID, a.Method, st.warmUp, a.WarmUp)
				}
			} else {
				fmt.Printf("%s reached max instances %d\n", a.AppID, a.MaxInstances)
			}
		}
		if sig.Scale.down {
			if mApp.App.Instances > a.MinInstances {
				st.coolDown++
				if st.coolDown >= a.CoolDown {
					fmt.Printf("%s scale down triggered with %d of %d signals of %s\n",
						a.AppID, st.coolDown, a.CoolDown, a.Method)
					a.doScale(mApp, -a.ScaleFactor)
					st.coolDown = 0
				} else {
					fmt.Printf("%s cooling down %s(%d of %d)\n",
						a.AppID, a.Method, st.coolDown, a.CoolDown)
				}
			} else {
				fmt.Printf("%s reached min instances %d\n", a.AppID, a.MinInstances)
			}
		}
	}

}

//EnsureMinMaxInstances scales up or down to get within Min-Max instances
func (a *App) EnsureMinMaxInstances(mApp MarathonApp) bool {
	diff := 0
	if mApp.App.Instances < a.MinInstances {
		diff = a.MinInstances - mApp.App.Instances
		fmt.Printf("%s will be scaled up by %d to reach minimum instances of %d\n",
			a.AppID, diff, a.MinInstances)
		a.doScale(mApp, diff)
	} else if mApp.App.Instances > a.MaxInstances {
		diff = a.MaxInstances - mApp.App.Instances
		fmt.Printf("%s will be scaled down by %d to reach maximum instances of %d\n",
			a.AppID, diff, a.MaxInstances)
		a.doScale(mApp, diff)
	}
	return diff == 0
}

func (a *App) doScale(mApp MarathonApp, instances int) {
	target := mApp.App.Instances + instances
	if target > a.MaxInstances {
		target = a.MaxInstances
	} else if target < a.MinInstances {
		target = a.MinInstances
	}
	fmt.Printf("Scaling %s to %d instances\n", a.AppID, target)
	client.ScaleMarathonApp(a.AppID, target)
}
