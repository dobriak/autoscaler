package main

import (
	"fmt"
)

//ScaleSignal describes a scale proposal
type ScaleSignal struct {
	CPU   scaleDirection
	Mem   scaleDirection
	Scale scaleDirection
}

type scaleDirection struct {
	up   bool
	down bool
}

//generateSignal given cpu and mem values, return a scale proposal
func (a *App) generateSignal(cpu, mem float64) ScaleSignal {
	result := ScaleSignal{}
	result.CPU.down = (cpu <= a.MinCPU)
	result.CPU.up = (cpu > a.MaxCPU)
	result.Mem.down = (mem <= a.MinMem)
	result.Mem.up = (mem > a.MinMem)
	switch method := a.Method; method {
	case "cpu":
		result.Scale.up = result.CPU.up
		result.Scale.down = result.CPU.down
	case "mem":
		result.Scale.up = result.Mem.up
		result.Scale.down = result.Mem.down
	case "and":
		result.Scale.up = result.CPU.up && result.Mem.up
		result.Scale.down = result.CPU.down && result.Mem.down
	case "or":
		result.Scale.up = result.CPU.up || result.Mem.up
		result.Scale.down = result.CPU.down || result.Mem.down
	default:
		fmt.Printf("method should be cpu|mem|and|or: %s\n", method)
		panic("Invalid parameter method.")
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
	sig := a.generateSignal(cpu, mem)
	if !sig.Scale.down && !sig.Scale.up {
		st.coolDown = 0
		st.warmUp = 0
	}
	if sig.Scale.up {
		if mApp.App.Instances < a.MaxInstances {
			st.warmUp++
			if st.warmUp >= a.WarmUp {
				fmt.Printf("%s scale up triggered with %d of %d signals", a.AppID, st.warmUp, a.WarmUp)
				//TODO: scale it up
				st.warmUp = 0
			} else {
				fmt.Printf("%s warming up (%d of %d)", a.AppID, st.warmUp, a.WarmUp)
			}
		} else {
			fmt.Printf("%s reached max instances %d", a.AppID, a.MaxInstances)
		}
	}
	if sig.Scale.down {
		if mApp.App.Instances > a.MinInstances {
			st.coolDown++
			if st.coolDown >= a.CoolDown {
				fmt.Printf("%s scale down triggered with %d of %d signals", a.AppID, st.coolDown, a.CoolDown)
				//TODO: scale it down
				st.coolDown = 0
			} else {
				fmt.Printf("%s cooling down (%d of %d)", a.AppID, st.coolDown, a.CoolDown)
			}
		} else {
			fmt.Printf("%s reached min instances %d", a.AppID, a.MinInstances)
		}
	}
}
