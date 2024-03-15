//go:build js && wasm
// +build js,wasm

package main

import (
	"syscall/js"
	"time"

	"main/dbdraw"
)

//go:wasmimport howdy JsDo
func JsDo()

func main() {
	var time_at_prev_tick = time.Now()
	sims := []dbdraw.Simulation{
		dbdraw.Sim1(),
		dbdraw.Sim2(),
		dbdraw.Sim3(),
		dbdraw.Sim4(),
		dbdraw.Sim5(),
		dbdraw.Sim6(),
	}

	storeInJs := func(this js.Value, args []js.Value) any {
		jsSims := []interface{}{}
		for _, sim := range sims {
			jsSims = append(jsSims, sim.ToJS())
		}
		return js.ValueOf(jsSims)
	}
	activateSim := func(this js.Value, args []js.Value) any {
		idx := args[0].Int()
		sims[idx].Activate()
		return js.Undefined()
	}
	disableSim := func(this js.Value, args []js.Value) any {
		idx := args[0].Int()
		sims[idx].Deactivate()
		return js.Undefined()
	}
	js.Global().Set("getSims", js.FuncOf(storeInJs))
	js.Global().Set("activateSim", js.FuncOf(activateSim))
	js.Global().Set("disableSim", js.FuncOf(disableSim))

	for {
		time_to_next_tick_truncated := time_at_prev_tick.Add(dbdraw.TIME_PER_TICK).Sub(time.Now()).Truncate(time.Millisecond)
		if time_to_next_tick_truncated > time.Microsecond {
			time.Sleep(time_to_next_tick_truncated)
		}

		time_after_sleep := time.Now()
		for time_after_sleep.Sub(time_at_prev_tick) > dbdraw.TIME_PER_TICK {
			time_at_prev_tick = time_at_prev_tick.Add(dbdraw.TIME_PER_TICK)
			for i := range sims {
				dbdraw.Update(&sims[i])
			}
			JsDo()
		}
	}
}
