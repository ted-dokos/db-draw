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

func main() {
	var i int32 = 1
	for {
		ang := float64(i) * math.Pi / 180.0
		d := Database{pos: Position{x: math.Cos(ang), y: math.Sin(ang)}}
		storeInJs := func(this js.Value, args []js.Value) any {
			return d.tojs()
		}
		js.Global().Set("callback", js.FuncOf(storeInJs))
		JsDo()
		fmt.Println(i)
		i++
		time.Sleep(time.Second)
	}
}
