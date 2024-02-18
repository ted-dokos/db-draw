//go:build js && wasm
// +build js,wasm

package main

import (
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

type Client struct {
	pos Position
}

func (c Client) tojs() js.Value {
	return js.ValueOf(map[string]interface{}{
		"pos": c.pos.tojs(),
	})
}

type Endpoint struct {
	ty  rune
	idx uint
}

func (e Endpoint) tojs() js.Value {
	return js.ValueOf(map[string]interface{}{
		"type":  js.ValueOf(string(e.ty)),
		"index": js.ValueOf(e.idx),
	})
}

type Channel struct {
	ep1 Endpoint
	ep2 Endpoint
}

func (c Channel) tojs() js.Value {
	return js.ValueOf(map[string]interface{}{
		"ep1": c.ep1.tojs(),
		"ep2": c.ep2.tojs(),
	})
}

type SimulationState struct {
	databases []Database
	clients   []Client
	channels  []Channel
}

func (s SimulationState) tojs() js.Value {
	dbs := make([]interface{}, len(s.databases))
	for i := 0; i < len(s.databases); i++ {
		dbs[i] = s.databases[i].tojs()
	}
	clients := make([]interface{}, len(s.clients))
	for i := 0; i < len(s.clients); i++ {
		clients[i] = s.clients[i].tojs()
	}
	channels := make([]interface{}, len(s.channels))
	for i := 0; i < len(s.channels); i++ {
		channels[i] = s.channels[i].tojs()
	}
	return js.ValueOf(map[string]interface{}{
		"databases": js.ValueOf(dbs),
		"clients":   js.ValueOf(clients),
		"channels":  js.ValueOf(channels),
	})
}

func main() {
	var i int32 = 1
	for {
		ang := float64(i) * math.Pi / 180.0
		d := Database{pos: Position{x: math.Cos(ang), y: math.Sin(ang)}}
		d2 := Database{pos: Position{x: 0.0, y: 0.0}}
		ch := Channel{ep1: Endpoint{ty: 'd', idx: 0}, ep2: Endpoint{ty: 'd', idx: 1}}
		sim := SimulationState{databases: []Database{d, d2}, channels: []Channel{ch}}
		storeInJs := func(this js.Value, args []js.Value) any {
			return sim.tojs()
		}
		js.Global().Set("callback", js.FuncOf(storeInJs))
		JsDo()
		i++
		time.Sleep(16 * time.Millisecond)
	}
}
