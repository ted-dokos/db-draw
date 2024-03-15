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

func (p Position) scale(s float64) Position {
	p.x *= s
	p.y *= s
	return p
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

var shapeNames = [...]string{"circle", "square", "triangle"}

type RequestType uint
type ShapeState uint

const (
	solid ShapeState = iota
	hstripe
	vstripe
)
const (
	read RequestType = iota
	write
)

type Request interface {
	JSable
	implementsRequest()
}

type ReadRequest struct {
	shape Shape
}

func (r ReadRequest) implementsRequest() {}
func (r ReadRequest) ToJS() js.Value {
	return js.ValueOf(map[string]any{
		"shape": js.ValueOf(uint(r.shape)),
	})
}

type WriteRequest struct {
	shape      Shape
	writeState ShapeState
}

func (w WriteRequest) implementsRequest() {}
func (w WriteRequest) ToJS() js.Value {
	return js.ValueOf(map[string]any{
		"shape":      js.ValueOf(uint(w.shape)),
		"writeState": js.ValueOf(uint(w.writeState)),
	})
}

type StatusCode uint

const (
	ok StatusCode = iota
	err
)

type Response interface {
	JSable
	implementsResponse()
}

type ReadResponse struct {
	status StatusCode
	shape  Shape
	state  ShapeState
}

func (r ReadResponse) implementsResponse() {}
func (r ReadResponse) ToJS() js.Value {
	return js.ValueOf(map[string]any{
		"shape":  js.ValueOf(uint(r.shape)),
		"state":  js.ValueOf(uint(r.state)),
		"status": js.ValueOf(uint(r.status)),
	})
}

type WriteResponse struct {
	status StatusCode
	shape  Shape
	state  ShapeState
}

func (w WriteResponse) implementsResponse() {}
func (w WriteResponse) ToJS() js.Value {
	return js.ValueOf(map[string]any{
		"shape":  js.ValueOf(uint(w.shape)),
		"state":  js.ValueOf(uint(w.state)),
		"status": js.ValueOf(uint(w.status)),
	})
}

type Endpoint interface {
	JSable
	Data() *DBData
	ReceiveRequest(r Request) Response
	ReceiveResponse(r Response)
}

type DBData map[Shape]ShapeState

type Database struct {
	pos    Position
	data   DBData
	leader bool
}

func (d *Database) Data() *DBData {
	return &d.data
}
func (d Database) ReceiveRequest(r Request) Response {
	switch r.(type) {
	case ReadRequest:
		read := r.(ReadRequest)
		return ReadResponse{ok, read.shape, ShapeState(d.data[read.shape])}
	case WriteRequest:
		write := r.(WriteRequest)
		d.data[write.shape] = write.writeState
		return WriteResponse{ok, write.shape, ShapeState(d.data[write.shape])}
	default:
		panic("Unknown request type")
	}
}
func (d Database) ReceiveResponse(r Response) {}
func (d *Database) ToJS() js.Value {
	data := map[string]interface{}{}
	for shape, reqType := range d.data {
		data[shapeNames[shape]] = js.ValueOf(uint(reqType))
	}
	return js.ValueOf(map[string]interface{}{
		"pos":    d.pos.ToJS(),
		"data":   js.ValueOf(data),
		"leader": js.ValueOf(d.leader),
	})
}

type Client struct {
	pos Position
}

func (c *Client) Data() *DBData {
	return nil
}
func (c Client) ReceiveRequest(r Request) Response {
	panic("Client::ReceiveRequest is not implemented.")
}
func (c Client) ReceiveResponse(r Response) {}
func (c *Client) ToJS() js.Value {
	return js.ValueOf(map[string]interface{}{
		"pos": c.pos.ToJS(),
	})
}

type PacketContentsType uint

const (
	request PacketContentsType = iota
	response
	writeResponse
)

type Packet struct {
	progress float64
	contents JSable
}

func (p Packet) ToJS() js.Value {
	v := map[string]interface{}{"progress": js.ValueOf(p.progress)}
	key := ""
	switch p.contents.(type) {
	case ReadRequest:
		key = "readRequest"
	case WriteRequest:
		key = "writeRequest"
	case ReadResponse:
		key = "readResponse"
	case WriteResponse:
		key = "writeResponse"
	default:
		panic("Unexpected packet contents.")
	}
	v[key] = p.contents.ToJS()
	return js.ValueOf(v)
}

type Channel struct {
	ep1        Endpoint
	ep2        Endpoint
	travelTime float64
	outgoing   *Packet
	incoming   *Packet
	invisible  bool
}

func (c *Channel) sendRequest(outgoing bool, r Request) {
	if outgoing {
		c.outgoing = &Packet{0.0, r}
	} else {
		c.incoming = &Packet{0.0, r}
	}
}
func (c *Channel) sendResponse(outgoing bool, r Response) {
	if outgoing {
		c.outgoing = &Packet{0.0, r}
	} else {
		c.incoming = &Packet{0.0, r}
	}
}
func (c Channel) ToJS() js.Value {
	maybe_packet := func(p *Packet) js.Value {
		if p == nil {
			return js.Null()
		}
		return p.ToJS()
	}
	return js.ValueOf(map[string]interface{}{
		"ep1":         c.ep1.ToJS(),
		"ep2":         c.ep2.ToJS(),
		"travel_time": js.ValueOf(c.travelTime),
		"outgoing":    maybe_packet(c.outgoing),
		"incoming":    maybe_packet(c.incoming),
	})
}

type Simulation struct {
	tick      uint
	databases []Database
	clients   []Client
	channels  []Channel
	events    EventEmitter
	active    bool
}

func (s *Simulation) Activate() {
	s.active = true
}
func (s *Simulation) Deactivate() {
	s.active = false
}
func (s Simulation) ToJS() js.Value {
	dbs := make([]any, len(s.databases))
	for i := 0; i < len(s.databases); i++ {
		dbs[i] = s.databases[i].ToJS()
	}
	clients := make([]any, len(s.clients))
	for i := 0; i < len(s.clients); i++ {
		clients[i] = s.clients[i].ToJS()
	}
	channels := make([]any, len(s.channels))
	for i := 0; i < len(s.channels); i++ {
		channels[i] = s.channels[i].ToJS()
	}
	return js.ValueOf(map[string]any{
		"databases": js.ValueOf(dbs),
		"clients":   js.ValueOf(clients),
		"channels":  js.ValueOf(channels),
		"active":    js.ValueOf(s.active),
		"tick":      js.ValueOf(s.tick),
	})
}
func Update(s *Simulation) {
	if !s.active {
		return
	}
	s.tick++
	s.events.ProcessTick(s.tick)
	updateDir := func(ch *Channel, outgoing bool) Response {
		ep := ch.ep1
		p := &ch.incoming
		if outgoing {
			ep = ch.ep2
			p = &ch.outgoing
		}
		if *p == nil {
			return nil
		}
		(*p).progress += TIME_PER_TICK.Seconds() / ch.travelTime
		if (*p).progress < 1.0 {
			return nil
		}
		var resp Response
		switch (*p).contents.(type) {
		case Request:
			resp = ep.ReceiveRequest((*p).contents.(Request))
		case Response:
			ep.ReceiveResponse((*p).contents.(Response))
		}
		*p = nil
		return resp
	}

	for i := 0; i < len(s.channels); i++ {
		ch := &s.channels[i]

		outResp := updateDir(ch, true)
		inResp := updateDir(ch, false)
		if outResp != nil {
			ch.sendResponse(false, outResp)
		}
		if inResp != nil {
			ch.sendResponse(true, inResp)
		}
	}
}
