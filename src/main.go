//go:build js && wasm
// +build js,wasm

package main

import (
	"fmt"
	"math/rand"
	"syscall/js"
	"time"
)

//go:wasmimport howdy JsDo
func JsDo()

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

type Shape uint

const (
	circle Shape = iota
	square
	triangle
)

type ShapeState uint

const (
	solid ShapeState = iota
	hstripe
	vstripe
	absent
)

type Database struct {
	pos  Position
	data map[Shape]ShapeState
}

func (d Database) tojs() js.Value {
	getState := func(key Shape, mp map[Shape]ShapeState) ShapeState {
		v, present := mp[key]
		if !present {
			return absent
		}
		return v
	}
	return js.ValueOf(map[string]interface{}{
		"pos": d.pos.tojs(),
		"data": js.ValueOf(map[string]interface{}{
			"circle":   js.ValueOf(uint(getState(circle, d.data))),
			"square":   js.ValueOf(uint(getState(square, d.data))),
			"triangle": js.ValueOf(uint(getState(triangle, d.data))),
		}),
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

type Op uint

const (
	read Op = iota
	write
)

type Transaction struct {
	progress float64
	shape    Shape
	style    ShapeState
	ty       Op
}

func maybe_transaction(t *Transaction) js.Value {
	if t == nil {
		return js.Null()
	}
	return t.tojs()
}
func (t Transaction) tojs() js.Value {
	return js.ValueOf(map[string]interface{}{
		"progress": js.ValueOf(t.progress),
		"shape":    js.ValueOf(uint(t.shape)),
		"style":    js.ValueOf(uint(t.style)),
		"type":     js.ValueOf(uint(t.ty)),
	})
}

type Channel struct {
	ep1         Endpoint
	ep2         Endpoint
	travel_time float64
	outgoing    *Transaction
	incoming    *Transaction
}

func (c *Channel) sendWrite(outgoing bool, shape Shape, style ShapeState) {
	t := Transaction{progress: 0.0, shape: shape, style: style, ty: write}
	if outgoing {
		c.outgoing = &t
	} else {
		c.incoming = &t
	}
}
func (c Channel) tojs() js.Value {
	return js.ValueOf(map[string]interface{}{
		"ep1":         c.ep1.tojs(),
		"ep2":         c.ep2.tojs(),
		"travel_time": js.ValueOf(c.travel_time),
		"outgoing":    maybe_transaction(c.outgoing),
		"incoming":    maybe_transaction(c.incoming),
	})
}

type SimulationState struct {
	tick      uint
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

func receive(s *SimulationState, t *Transaction, e *Endpoint) {
	if e.ty == 'd' {
		s.databases[e.idx].data[t.shape] = t.style
	}
}

// func updateChannelDirection(s *SimulationState, t *Transaction, e *Endpoint, travel_time float64) {
// 	if t == nil {
// 		return
// 	}
// 	t.progress += TIME_PER_TICK.Seconds() / travel_time
// 	if t.progress < 1.0 {
// 		return
// 	}
// 	if t.ty == write {
// 		receive(s, t, e)
// 	}
// }

func update(s *SimulationState) {
	s.tick++
	style := rand.Int() % 3
	if s.tick == 200 {

		s.channels[0].sendWrite(true, triangle, ShapeState(style))
	}
	if s.tick == 600 {
		s.channels[0].sendWrite(true, square, ShapeState(style))
	}
	if s.tick == 1000 {
		s.channels[0].sendWrite(true, circle, ShapeState(style))
	}
	if s.tick == 1400 {
		s.channels[0].sendWrite(true, square, ShapeState(style))
	}
	if s.tick > 1400 && (s.tick-200)%400 == 0 {
		shape := rand.Int() % 3
		s.channels[0].sendWrite(true, Shape(shape), ShapeState(style))
	}
	for i := 0; i < len(s.channels); i++ {
		ch := &s.channels[i]
		updateDir := func(dir *Transaction, ep *Endpoint) *Transaction {
			if dir == nil {
				return nil
			}
			dir.progress += TIME_PER_TICK.Seconds() / ch.travel_time
			if dir.progress < 1.0 {
				return dir
			}
			if dir.ty == write {
				receive(s, dir, ep)
			}
			return nil
		}
		ch.outgoing = updateDir(ch.outgoing, &ch.ep2)
		ch.incoming = updateDir(ch.incoming, &ch.ep1)
	}
}

func make_intro_sim() SimulationState {
	d := Database{pos: Position{x: 0.5, y: 0.0}, data: make(map[Shape]ShapeState)}
	client := Client{pos: Position{x: -0.5, y: 0.0}}
	ch := Channel{ep1: Endpoint{ty: 'c', idx: 0}, ep2: Endpoint{ty: 'd', idx: 0}, travel_time: 2.0}
	return SimulationState{databases: []Database{d}, channels: []Channel{ch}, clients: []Client{client}}
}

const TICKS_PER_SECOND = 100.0
const TIME_PER_TICK = time.Second / TICKS_PER_SECOND

func main() {
	var time_at_prev_tick = time.Now()
	sims := []SimulationState{make_intro_sim(), make_intro_sim()}
	current_sim_idx := 0

	storeInJs := func(this js.Value, args []js.Value) any {
		return sims[current_sim_idx].tojs()
	}
	setIdx := func(this js.Value, args []js.Value) any {
		current_sim_idx = args[0].Int()
		fmt.Println(current_sim_idx)
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
