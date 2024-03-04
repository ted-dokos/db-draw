package main

import (
	"syscall/js"
	"time"
)

const TICKS_PER_SECOND = 100.0
const TIME_PER_TICK = time.Second / TICKS_PER_SECOND

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

type DBData map[Shape]ShapeState

type Endpoint interface {
	JSable
	Data() *DBData
}

type Database struct {
	pos  Position
	data DBData
}

func (d Database) Data() *DBData {
	return &d.data
}

func (d Database) tojs() js.Value {
	getState := func(key Shape, mp DBData) ShapeState {
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

func (c Client) Data() *DBData {
	return nil
}

func (c Client) tojs() js.Value {
	return js.ValueOf(map[string]interface{}{
		"pos": c.pos.tojs(),
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
	events    EventEmitter
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
	dbdata := (*e).Data()
	if dbdata != nil {
		(*dbdata)[t.shape] = t.style
	}
}

func update(s *SimulationState) {
	s.tick++
	s.events.ProcessTick(s.tick)
	// style := rand.Int() % 3
	// if s.tick == 200 {

	// 	s.channels[0].sendWrite(true, triangle, ShapeState(style))
	// }
	// if s.tick == 600 {
	// 	s.channels[0].sendWrite(true, square, ShapeState(style))
	// }
	// if s.tick == 1000 {
	// 	s.channels[0].sendWrite(true, circle, ShapeState(style))
	// }
	// if s.tick == 1400 {
	// 	s.channels[0].sendWrite(true, square, ShapeState(style))
	// }
	// if s.tick > 1400 && (s.tick-200)%400 == 0 {
	// 	shape := rand.Int() % 3
	// 	s.channels[0].sendWrite(true, Shape(shape), ShapeState(style))
	// }
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
	dbs := []Database{{pos: Position{x: 0.5, y: 0.0}, data: map[Shape]ShapeState{}}}
	clients := []Client{{pos: Position{x: -0.5, y: 0.0}}}
	chs := []Channel{{ep1: clients[0], ep2: dbs[0], travel_time: 2.0}}
	emitters := compose_emitters(
		OnceEmitter{
			tick: 100,
			emit: ChannelEmitter{c: &chs[0], outgoing: true, sendee: independent(triFunc, solidFunc)},
		},
		OnceEmitter{
			tick: 400,
			emit: ChannelEmitter{c: &chs[0], outgoing: true, sendee: independent(sqFunc, hstripeFunc)},
		},
		OnceEmitter{
			tick: 700,
			emit: ChannelEmitter{c: &chs[0], outgoing: true, sendee: independent(circFunc, vstripeFunc)},
		},
		PeriodicEmitter{
			first_tick: 1000,
			period:     300,
			emit:       ChannelEmitter{c: &chs[0], outgoing: true, sendee: randShapeWithNewStyle},
		})
	return SimulationState{
		0,
		dbs,
		clients,
		chs,
		emitters,
	}
}
