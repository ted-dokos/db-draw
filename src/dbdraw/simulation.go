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

type Request struct {
	reqType RequestType
	shape   Shape
	// Only used for write requests
	writeState ShapeState
}

func (r Request) ToJS() js.Value {
	v := map[string]interface{}{
		"reqType": js.ValueOf(uint(r.reqType)),
		"shape":   js.ValueOf(uint(r.shape)),
	}
	if r.reqType == write {
		v["writeState"] = js.ValueOf(uint(r.writeState))
	}
	return js.ValueOf(v)
}

type StatusCode uint

const (
	ok StatusCode = iota
	err
)

type Response struct {
	shape  Shape
	state  ShapeState
	status StatusCode
}

func (r Response) ToJS() js.Value {
	return js.ValueOf(map[string]any{
		"shape":  js.ValueOf(uint(r.shape)),
		"state":  js.ValueOf(uint(r.state)),
		"status": js.ValueOf(uint(r.status)),
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
	pos  Position
	data DBData
}

func (d *Database) Data() *DBData {
	return &d.data
}

func (d Database) ReceiveRequest(r Request) Response {
	if r.reqType == read {
		return Response{r.shape, ShapeState(d.data[r.shape]), ok}
	}
	d.data[r.shape] = r.writeState
	return Response{r.shape, ShapeState(d.data[r.shape]), ok}
}
func (d Database) ReceiveResponse(r Response) {}
func (d *Database) ToJS() js.Value {
	data := map[string]interface{}{}
	for shape, reqType := range d.data {
		data[shapeNames[shape]] = js.ValueOf(uint(reqType))
	}
	return js.ValueOf(map[string]interface{}{
		"pos":  d.pos.ToJS(),
		"data": js.ValueOf(data),
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
	progress    float64
	contentType PacketContentsType
	contents    JSable
}

func (p Packet) ToJS() js.Value {
	v := map[string]interface{}{"progress": js.ValueOf(p.progress)}
	switch p.contentType {
	case request:
		v["request"] = p.contents.ToJS()
	case response:
		v["readResponse"] = p.contents.ToJS()
	case writeResponse:
		v["writeResponse"] = p.contents.ToJS()
	}
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
		c.outgoing = &Packet{0.0, request, &r}
	} else {
		c.incoming = &Packet{0.0, request, &r}
	}
}
func (c *Channel) sendResponse(outgoing bool, r Response, ty PacketContentsType) {
	if outgoing {
		c.outgoing = &Packet{0.0, ty, &r}
	} else {
		c.incoming = &Packet{0.0, ty, &r}
	}
}
func (c Channel) ToJS() js.Value {
	maybe_packet := func(t *Packet) js.Value {
		if t == nil {
			return js.Null()
		}
		return t.ToJS()
	}
	return js.ValueOf(map[string]interface{}{
		"ep1":         c.ep1.ToJS(),
		"ep2":         c.ep2.ToJS(),
		"travel_time": js.ValueOf(c.travelTime),
		"outgoing":    maybe_packet(c.outgoing),
		"incoming":    maybe_packet(c.incoming),
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

func Update(s *SimulationState) {
	s.tick++
	s.events.ProcessTick(s.tick)
	for i := 0; i < len(s.channels); i++ {
		ch := &s.channels[i]
		updateDir := func(ch *Channel, outgoing bool) {
			ep := ch.ep1
			p := &ch.incoming
			if outgoing {
				ep = ch.ep2
				p = &ch.outgoing
			}
			if *p == nil {
				return
			}
			(*p).progress += TIME_PER_TICK.Seconds() / ch.travelTime
			if (*p).progress < 1.0 {
				return
			}
			switch (*p).contentType {
			case request:
				req := (*p).contents.(*Request)
				resp := ep.ReceiveRequest(*req)
				ty := response
				if req.reqType == write {
					ty = writeResponse
				}
				ch.sendResponse(!outgoing, resp, ty)
			case response:
				ep.ReceiveResponse(*(*p).contents.(*Response))
			}
			*p = nil
		}
		updateDir(ch, true)
		updateDir(ch, false)
	}
}
