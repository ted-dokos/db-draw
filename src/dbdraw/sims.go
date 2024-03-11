package dbdraw

import "math"

// Same order as they appear on the page.

func Sim1() SimulationState {
	dbs := []Database{{
		pos: Position{x: 0.0, y: 0.0},
		data: DBData{
			circle:   solid,
			square:   hstripe,
			triangle: vstripe,
		},
	}}
	clientPosScale := 0.6
	clients := []Client{
		{pos: Position{x: -clientPosScale, y: 0.0}},
		{pos: Position{x: 0.5, y: math.Sqrt(3) / 2.0}.scale(clientPosScale)},
		{pos: Position{x: 0.5, y: -math.Sqrt(3) / 2.0}.scale(clientPosScale)},
	}
	chs := []Channel{
		{ep1: &clients[0], ep2: &dbs[0], travelTime: 2.0},
		{ep1: &clients[1], ep2: &dbs[0], travelTime: 2.0},
		{ep1: &clients[2], ep2: &dbs[0], travelTime: 2.0},
	}
	emitters := compose_emitters(
		PeriodicEmitter{
			first_tick: 100,
			period:     300,
			emit: ChannelEmitter{
				c:        &chs[0],
				outgoing: true,
				sendee:   randShapeWithNewStyle,
			},
		},
		PeriodicEmitter{
			first_tick: 200,
			period:     300,
			emit: ChannelEmitter{
				c:        &chs[1],
				outgoing: true,
				sendee:   randRead,
			},
		},
		PeriodicEmitter{
			first_tick: 300,
			period:     300,
			emit: ChannelEmitter{
				c:        &chs[2],
				outgoing: true,
				sendee:   randReadOrWrite,
			},
		})
	return SimulationState{
		0,
		dbs,
		clients,
		chs,
		emitters,
	}
}

func Sim2() SimulationState {
	dbPosScale := 0.4
	dbs := []Database{
		{
			pos:  Position{x: 1.0, y: 0.0}.scale(dbPosScale),
			data: DBData{circle: solid},
		},
		{
			pos:  Position{x: -0.5, y: math.Sqrt(3) / 2.0}.scale(dbPosScale),
			data: DBData{square: hstripe},
		},
		{
			pos:  Position{x: -0.5, y: -math.Sqrt(3) / 2.0}.scale(dbPosScale),
			data: DBData{triangle: vstripe},
		},
	}
	clients := []Client{}
	channels := []Channel{}
	emitters := compose_emitters()
	return SimulationState{
		0,
		dbs,
		clients,
		channels,
		emitters,
	}
}

func Sim3() SimulationState {
	dbs := []Database{
		{
			pos: Position{x: -0.3, y: 0.0},
			data: DBData{
				circle:   solid,
				square:   hstripe,
				triangle: vstripe,
			},
		},
		{
			pos: Position{x: 0.3, y: 0.0},
			data: DBData{
				circle:   solid,
				square:   hstripe,
				triangle: vstripe,
			},
		},
	}
	clients := []Client{
		{Position{-0.8, 0.3}},
		{Position{-0.8, -0.3}},
		{Position{0.8, 0.3}},
		{Position{0.8, -0.3}},
	}
	channels := []Channel{
		{ep1: &clients[0], ep2: &dbs[0], travelTime: 2.0},
		{ep1: &clients[1], ep2: &dbs[0], travelTime: 2.0},
		{ep1: &clients[2], ep2: &dbs[1], travelTime: 2.0},
		{ep1: &clients[3], ep2: &dbs[1], travelTime: 2.0},
	}
	emitters := compose_emitters(
		PeriodicEmitter{
			first_tick: 100,
			period:     300,
			emit:       ChannelEmitter{c: &channels[0], outgoing: true, sendee: randRead},
		},
		PeriodicEmitter{
			first_tick: 175,
			period:     300,
			emit:       ChannelEmitter{c: &channels[1], outgoing: true, sendee: randRead},
		},
		PeriodicEmitter{
			first_tick: 250,
			period:     300,
			emit:       ChannelEmitter{c: &channels[2], outgoing: true, sendee: randRead},
		},
		PeriodicEmitter{
			first_tick: 325,
			period:     300,
			emit:       ChannelEmitter{c: &channels[3], outgoing: true, sendee: randRead},
		},
	)
	return SimulationState{
		0,
		dbs,
		clients,
		channels,
		emitters,
	}
}

func Sim4() SimulationState {
	dbs := []Database{
		{
			pos: Position{x: -0.3, y: 0.0},
			data: DBData{
				circle:   solid,
				square:   hstripe,
				triangle: vstripe,
			},
		},
		{
			pos: Position{x: 0.3, y: 0.0},
			data: DBData{
				circle:   solid,
				square:   hstripe,
				triangle: vstripe,
			},
		},
	}
	clients := []Client{
		{Position{-0.8, 0.3}},
		{Position{-0.8, -0.3}},
		{Position{0.8, 0.3}},
		{Position{0.8, -0.3}},
	}
	channels := []Channel{
		{ep1: &clients[0], ep2: &dbs[0], travelTime: 2.0},
		{ep1: &clients[1], ep2: &dbs[0], travelTime: 2.0},
		{ep1: &clients[2], ep2: &dbs[1], travelTime: 2.0},
		{ep1: &clients[3], ep2: &dbs[1], travelTime: 2.0},
	}
	emitters := compose_emitters(
		OnceEmitter{
			tick: 100,
			emit: ChannelEmitter{c: &channels[0], outgoing: true, sendee: randShapeWithNewStyle},
		},
	)
	return SimulationState{
		0,
		dbs,
		clients,
		channels,
		emitters,
	}
}
