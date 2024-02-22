//go:build js && wasm
// +build js,wasm

package main

import (
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

type TransactionShape uint

const (
	square TransactionShape = iota
	triangle
	circle
)

type TransactionStyle uint

const (
	solid TransactionStyle = iota
	stripe
	stipple
)

type Transaction struct {
	progress float64
	shape    TransactionShape
	style    TransactionStyle
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
	})
}

type Channel struct {
	ep1         Endpoint
	ep2         Endpoint
	travel_time float64
	outgoing    *Transaction
	incoming    *Transaction
}

func (c *Channel) send(outgoing bool, shape TransactionShape, style TransactionStyle) {
	t := Transaction{progress: 0.0, shape: shape, style: style}
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

func update(tick uint, s SimulationState) {
	// ang := float64(tick) * math.Pi / 180.0
	// s.databases[0].pos = Position{x: math.Cos(ang), y: math.Sin(ang)}
	if tick == 200 {
		s.channels[0].send(true, triangle, stripe)
	}
	if tick == 600 {
		s.channels[0].send(true, triangle, solid)
	}
	if tick == 1000 {
		s.channels[0].send(true, circle, solid)
	}
	for i := 0; i < len(s.channels); i++ {
		ch := &s.channels[i]
		if ch.outgoing != nil {
			ch.outgoing.progress += TIME_PER_TICK.Seconds() / ch.travel_time
			if ch.outgoing.progress >= 1.0 {
				ch.outgoing = nil
			}
		}
		if ch.incoming != nil {
			ch.incoming.progress += TIME_PER_TICK.Seconds() / ch.travel_time
			if ch.incoming.progress >= 1.0 {
				ch.incoming = nil
			}
		}
	}
}

const TICKS_PER_SECOND = 100.0
const TIME_PER_TICK = time.Second / TICKS_PER_SECOND

func main() {
	var tick uint = 0
	var time_at_prev_tick = time.Now()
	d := Database{pos: Position{x: 1.0, y: 0.0}}
	//d2 := Database{pos: Position{x: 0.0, y: 0.0}}
	client := Client{pos: Position{x: -0.5, y: 0.0}}
	//ch := Channel{ep1: Endpoint{ty: 'd', idx: 0}, ep2: Endpoint{ty: 'd', idx: 1}}
	ch2 := Channel{ep1: Endpoint{ty: 'c', idx: 0}, ep2: Endpoint{ty: 'd', idx: 0}, travel_time: 3.0}
	sim := SimulationState{databases: []Database{d /*, d2*/}, channels: []Channel{ /*ch,*/ ch2}, clients: []Client{client}}

	storeInJs := func(this js.Value, args []js.Value) any {
		return sim.tojs()
	}
	js.Global().Set("callback", js.FuncOf(storeInJs))

	for {
		time_to_next_tick_truncated := time_at_prev_tick.Add(TIME_PER_TICK).Sub(time.Now()).Truncate(time.Millisecond)
		if time_to_next_tick_truncated > time.Microsecond {
			time.Sleep(time_to_next_tick_truncated)
		}

		time_after_sleep := time.Now()
		for time_after_sleep.Sub(time_at_prev_tick) > TIME_PER_TICK {
			tick++
			update(tick, sim)
			JsDo()
			time_at_prev_tick = time_at_prev_tick.Add(TIME_PER_TICK)
		}
	}
}
