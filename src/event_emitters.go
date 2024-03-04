package main

import (
	"math/rand"
)

type EventEmitter interface {
	ProcessTick(tick uint)
}

type ChannelEmitter struct {
	c        *Channel
	outgoing bool
	sendee   func(ChannelEmitter) (Shape, ShapeState)
}

func (emitter ChannelEmitter) ProcessTick(tick uint) {
	shape, state := emitter.sendee(emitter)
	emitter.c.sendWrite(emitter.outgoing, shape, state)
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
func randShapeWithNewStyle(emitter ChannelEmitter) (Shape, ShapeState) {
	ep := emitter.c.ep1
	if emitter.outgoing {
		ep = emitter.c.ep2
	}
	db := ep.Data()
	shape := randShape(emitter)
	if db == nil {
		return shape, randStyle(emitter)
	}
	style, ok := (*db)[shape]
	if !ok {
		return shape, randStyle(emitter)
	}
	new := (int(style) + rand.Int()%2 + 1) % 3
	return shape, ShapeState(new)
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
func randStyle(ChannelEmitter) ShapeState {
	return ShapeState(rand.Int() % 3)
}

func independent(shapeFunc func(ChannelEmitter) Shape, styleFunc func(ChannelEmitter) ShapeState) func(ChannelEmitter) (Shape, ShapeState) {
	return func(emitter ChannelEmitter) (Shape, ShapeState) {
		return shapeFunc(emitter), styleFunc(emitter)
	}
}
