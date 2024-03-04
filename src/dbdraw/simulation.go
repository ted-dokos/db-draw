package dbdraw

import (
	"syscall/js"
	"time"
)

const TICKS_PER_SECOND = 100.0
const TIME_PER_TICK = time.Second / TICKS_PER_SECOND

type JSable interface {
	ToJS() js.Value
}

type Position struct {
	x float64
	y float64
}

func (p Position) ToJS() js.Value {
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

func (d *Database) Data() *DBData {
	return &d.data
}

func (d *Database) ToJS() js.Value {
	getState := func(key Shape, mp DBData) ShapeState {
		v, present := mp[key]
		if !present {
			return absent
		}
		return v
	}
	return js.ValueOf(map[string]interface{}{
		"pos": d.pos.ToJS(),
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

func (c *Client) Data() *DBData {
	return nil
}

func (c *Client) ToJS() js.Value {
	return js.ValueOf(map[string]interface{}{
		"pos": c.pos.ToJS(),
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
	return t.ToJS()
}
func (t Transaction) ToJS() js.Value {
	return js.ValueOf(map[string]interface{}{
		"progress": js.ValueOf(t.progress),
		"shape":    js.ValueOf(uint(t.shape)),
		"style":    js.ValueOf(uint(t.style)),
		"type":     js.ValueOf(uint(t.ty)),
	})
}

type Channel struct {
	ep1        Endpoint
	ep2        Endpoint
	travelTime float64
	outgoing   *Transaction
	incoming   *Transaction
}

func (c *Channel) sendWrite(outgoing bool, shape Shape, style ShapeState) {
	t := Transaction{progress: 0.0, shape: shape, style: style, ty: write}
	if outgoing {
		c.outgoing = &t
	} else {
		c.incoming = &t
	}
}
func (c Channel) ToJS() js.Value {
	return js.ValueOf(map[string]interface{}{
		"ep1":         c.ep1.ToJS(),
		"ep2":         c.ep2.ToJS(),
		"travel_time": js.ValueOf(c.travelTime),
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

func (s SimulationState) ToJS() js.Value {
	dbs := make([]interface{}, len(s.databases))
	for i := 0; i < len(s.databases); i++ {
		dbs[i] = s.databases[i].ToJS()
	}
	clients := make([]interface{}, len(s.clients))
	for i := 0; i < len(s.clients); i++ {
		clients[i] = s.clients[i].ToJS()
	}
	channels := make([]interface{}, len(s.channels))
	for i := 0; i < len(s.channels); i++ {
		channels[i] = s.channels[i].ToJS()
	}
	return js.ValueOf(map[string]interface{}{
		"databases": js.ValueOf(dbs),
		"clients":   js.ValueOf(clients),
		"channels":  js.ValueOf(channels),
	})
}

func receive(s *SimulationState, t *Transaction, e Endpoint) {
	dbdata := e.Data()
	if dbdata != nil {
		(*dbdata)[t.shape] = t.style
	}
}

func Update(s *SimulationState) {
	s.tick++
	s.events.ProcessTick(s.tick)
	for i := 0; i < len(s.channels); i++ {
		ch := &s.channels[i]
		updateDir := func(dir *Transaction, ep Endpoint) *Transaction {
			if dir == nil {
				return nil
			}
			dir.progress += TIME_PER_TICK.Seconds() / ch.travelTime
			if dir.progress < 1.0 {
				return dir
			}
			if dir.ty == write {
				receive(s, dir, ep)
			}
			return nil
		}
		ch.outgoing = updateDir(ch.outgoing, ch.ep2)
		ch.incoming = updateDir(ch.incoming, ch.ep1)
	}
}
