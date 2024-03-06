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
	sendee   func(ChannelEmitter) (Shape, RequestType)
}

func (emitter ChannelEmitter) ProcessTick(tick uint) {
	shape, req := emitter.sendee(emitter)
	emitter.c.send(emitter.outgoing, shape, req)
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
func randShapeWithNewStyle(emitter ChannelEmitter) (Shape, RequestType) {
	ep := emitter.c.ep1
	if emitter.outgoing {
		ep = emitter.c.ep2
	}
	db := ep.Data()
	shape := randShape(emitter)
	if db == nil {
		return shape, randWrite(emitter)
	}
	style, ok := (*db)[shape]
	if !ok {
		return shape, randWrite(emitter)
	}
	new := (int(style) + rand.Int()%2 + 1) % 3
	return shape, RequestType(new)
}

func solidFunc(ChannelEmitter) RequestType {
	return solid
}
func hstripeFunc(ChannelEmitter) RequestType {
	return hstripe
}
func vstripeFunc(ChannelEmitter) RequestType {
	return vstripe
}
func randWrite(ChannelEmitter) RequestType {
	return RequestType(rand.Int() % 3)
}

func independent(shapeFunc func(ChannelEmitter) Shape, styleFunc func(ChannelEmitter) RequestType) func(ChannelEmitter) (Shape, RequestType) {
	return func(emitter ChannelEmitter) (Shape, RequestType) {
		return shapeFunc(emitter), styleFunc(emitter)
	}
}

func randRead(emitter ChannelEmitter) (Shape, RequestType) {
	return randShape(emitter), read
}
