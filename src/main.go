//go:build js && wasm
// +build js,wasm

package main

import (
	"syscall/js"
	"time"
)

//go:wasmimport howdy JsDo
func JsDo()

func main() {
	var time_at_prev_tick = time.Now()
	sims := []SimulationState{make_intro_sim(), make_intro_sim()}
	current_sim_idx := 0

	storeInJs := func(this js.Value, args []js.Value) any {
		return sims[current_sim_idx].tojs()
	}
	setIdx := func(this js.Value, args []js.Value) any {
		current_sim_idx = args[0].Int()
		return js.Undefined()
	}
	js.Global().Set("callback", js.FuncOf(storeInJs))
	js.Global().Set("setSimIndex", js.FuncOf(setIdx))

	for {
		time_to_next_tick_truncated := time_at_prev_tick.Add(TIME_PER_TICK).Sub(time.Now()).Truncate(time.Millisecond)
		if time_to_next_tick_truncated > time.Microsecond {
			time.Sleep(time_to_next_tick_truncated)
		}

		time_after_sleep := time.Now()
		for time_after_sleep.Sub(time_at_prev_tick) > TIME_PER_TICK {
			time_at_prev_tick = time_at_prev_tick.Add(TIME_PER_TICK)
			i := current_sim_idx
			if i < 0 {
				continue
			}
			update(&sims[i])
			JsDo()
		}
	}
}
