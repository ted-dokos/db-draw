package main

import "fmt"

type EventEmitter interface {
	ProcessFrame(tick uint)
}

type ChannelEmitter struct {
	c     *Channel
	c_evt ChannelEvent
}

func (c_emt ChannelEmitter) ProcessFrame(tick uint) {
	c_emt.c.sendWrite(c_emt.c_evt.outgoing,
		c_emt.c_evt.shape,
		c_emt.c_evt.style)
}

type OnceEmitter struct {
	tick uint
	emit EventEmitter
}

func (o OnceEmitter) ProcessFrame(tick uint) {
	if tick == o.tick {
		fmt.Println("foo")
		o.emit.ProcessFrame(tick)
	}
}

type PeriodicEmitter struct {
	first_tick uint
	period     uint
	emit       EventEmitter
}

func (p PeriodicEmitter) ProcessFrame(tick uint) {
	if tick < p.first_tick {
		return
	}
	if (tick-p.first_tick)%p.period == 0 {
		p.emit.ProcessFrame(tick)
	}
}

type ComposedEmitter struct {
	emits []EventEmitter
}

func (c ComposedEmitter) ProcessFrame(tick uint) {
	for i := 0; i < len(c.emits); i++ {
		c.emits[i].ProcessFrame(tick)
	}
}

func compose_emitters(emits ...EventEmitter) EventEmitter {
	return ComposedEmitter{emits}
}
