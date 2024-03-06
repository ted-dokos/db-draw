//go:build js && wasm
// +build js,wasm

package main

import (
	"dbdraw/dbdraw"
	"syscall/js"
	"time"
)

//go:wasmimport howdy JsDo
func JsDo()

func main() {
	var time_at_prev_tick = time.Now()
	sims := []dbdraw.SimulationState{dbdraw.Sim1(),
		dbdraw.Sim2(),
		dbdraw.Sim3(),
	}
	current_sim_idx := 0

	storeInJs := func(this js.Value, args []js.Value) any {
		return sims[current_sim_idx].ToJS()
	}
	setIdx := func(this js.Value, args []js.Value) any {
		current_sim_idx = args[0].Int()
		return js.Undefined()
	}
	js.Global().Set("callback", js.FuncOf(storeInJs))
	js.Global().Set("setSimIndex", js.FuncOf(setIdx))

	for {
		time_to_next_tick_truncated := time_at_prev_tick.Add(dbdraw.TIME_PER_TICK).Sub(time.Now()).Truncate(time.Millisecond)
		if time_to_next_tick_truncated > time.Microsecond {
			time.Sleep(time_to_next_tick_truncated)
		}

		time_after_sleep := time.Now()
		for time_after_sleep.Sub(time_at_prev_tick) > dbdraw.TIME_PER_TICK {
			time_at_prev_tick = time_at_prev_tick.Add(dbdraw.TIME_PER_TICK)
			i := current_sim_idx
			if i < 0 {
				continue
			}
			dbdraw.Update(&sims[i])
			JsDo()
		}
	}
}
