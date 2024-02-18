//go:build js && wasm
// +build js,wasm

package main

import (
	"fmt"
	"math"
	"syscall/js"
	"time"
)

//go:wasmimport howdy JsDo
func JsDo()

//func JsDo(i int32)

type JSable interface {
	tojs() js.Value
}

type Position struct {
	x float64
	y float64
}

func (p Position) tojs() js.Value {
	return js.ValueOf(map[string]interface{}{
		"x": js.ValueOf(p.x),
		"y": js.ValueOf(p.y),
	})
}

type Database struct {
	pos Position
}

func (d Database) tojs() js.Value {
	return js.ValueOf(map[string]interface{}{
		"pos": d.pos.tojs(),
	})
}

type SimulationState struct {
	databases []Database
}

func (s SimulationState) tojs() js.Value {
	var x []interface{}
	for i := 0; i < len(s.databases); i++ {
		x = append(x, s.databases[i].tojs())
	}
	return js.ValueOf(map[string]interface{}{
		"databases": js.ValueOf(x),
	})
}

func main() {
	var i int32 = 1
	for {
		ang := float64(i) * math.Pi / 180.0
		d := Database{pos: Position{x: math.Cos(ang), y: math.Sin(ang)}}
		d2 := Database{pos: Position{x: 0.0, y: 0.0}}
		sim := SimulationState{databases: []Database{d, d2}}
		storeInJs := func(this js.Value, args []js.Value) any {
			return sim.tojs()
		}
		js.Global().Set("callback", js.FuncOf(storeInJs))
		JsDo()
		fmt.Println(i)
		i++
		time.Sleep(16 * time.Millisecond)
	}
}
