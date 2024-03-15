package dbdraw

import (
	"math/rand"
)

type EventEmitter interface {
	ProcessTick(tick uint)
}

type ChannelEmitter struct {
	c        *Channel
	outgoing bool
	sendee   func(ChannelEmitter) Request
}

func (emitter ChannelEmitter) ProcessTick(tick uint) {
	req := emitter.sendee(emitter)
	emitter.c.sendRequest(emitter.outgoing, req)
}

type OnceEmitter struct {
	tick uint
	emit EventEmitter
}

func (o OnceEmitter) ProcessTick(tick uint) {
	if tick == o.tick {
		o.emit.ProcessTick(tick)
	}
}

type PeriodicEmitter struct {
	first_tick uint
	period     uint
	emit       EventEmitter
}

func (p PeriodicEmitter) ProcessTick(tick uint) {
	if tick < p.first_tick {
		return
	}
	if (tick-p.first_tick)%p.period == 0 {
		p.emit.ProcessTick(tick)
	}
}

type ComposedEmitter struct {
	emits []EventEmitter
}

func (c ComposedEmitter) ProcessTick(tick uint) {
	for i := 0; i < len(c.emits); i++ {
		c.emits[i].ProcessTick(tick)
	}
}

func compose_emitters(emits ...EventEmitter) EventEmitter {
	return ComposedEmitter{emits}
}

func circFunc(ChannelEmitter) Shape {
	return circle
}
func sqFunc(ChannelEmitter) Shape {
	return square
}
func triFunc(ChannelEmitter) Shape {
	return triangle
}
func randShape(ChannelEmitter) Shape {
	return Shape(rand.Int() % 3)
}
func randShapeWithNewStyle(emitter ChannelEmitter) Request {
	ep := emitter.c.ep1
	if emitter.outgoing {
		ep = emitter.c.ep2
	}
	db := ep.Data()
	shape := randShape(emitter)
	if db == nil {
		return WriteRequest{
			shape:      shape,
			writeState: randState(emitter),
		}
	}
	state, ok := (*db)[shape]
	if !ok {
		return WriteRequest{
			shape:      shape,
			writeState: randState(emitter),
		}
	}
	new := (int(state) + rand.Int()%2 + 1) % 3
	return WriteRequest{
		shape:      shape,
		writeState: ShapeState(new),
	}
}

func solidFunc(ChannelEmitter) ShapeState {
	return solid
}
func hstripeFunc(ChannelEmitter) ShapeState {
	return hstripe
}
func vstripeFunc(ChannelEmitter) ShapeState {
	return vstripe
}
func randState(ChannelEmitter) ShapeState {
	return ShapeState(rand.Int() % 3)
}

func independent(shapeFunc func(ChannelEmitter) Shape, stateFunc func(ChannelEmitter) ShapeState) func(ChannelEmitter) Request {
	return func(emitter ChannelEmitter) Request {
		return WriteRequest{
			shape:      shapeFunc(emitter),
			writeState: stateFunc(emitter),
		}
	}
}

func randRead(emitter ChannelEmitter) Request {
	return ReadRequest{
		shape: randShape(emitter),
	}
}

func randReadOrWrite(emitter ChannelEmitter) Request {
	makeReadReq := rand.Int()%2 == 0
	if makeReadReq {
		return ReadRequest{
			shape: randShape(emitter),
		}
	} else {
		return WriteRequest{
			shape:      randShape(emitter),
			writeState: randState(emitter),
		}
	}
}
