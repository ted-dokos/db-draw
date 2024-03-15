package dbdraw

import "math"

// Same order as they appear on the page.

func defaultDBData() DBData {
	return DBData{
		circle:   solid,
		square:   hstripe,
		triangle: vstripe,
	}
}

func Sim1() Simulation {
	dbs := []Database{{
		pos:  Position{x: 0.0, y: 0.0},
		data: defaultDBData(),
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
	return Simulation{
		tick:      0,
		databases: dbs,
		clients:   clients,
		channels:  chs,
		events:    emitters,
		active:    false,
	}
}

func Sim2() Simulation {
	dbPosScale := 0.4
	sin := math.Sqrt(3) / 2.0
	dbs := []Database{
		{
			pos:  Position{x: 1.0, y: 0.0}.scale(dbPosScale),
			data: DBData{circle: solid},
		},
		{
			pos:  Position{x: -0.5, y: sin}.scale(dbPosScale),
			data: DBData{square: hstripe},
		},
		{
			pos:  Position{x: -0.5, y: -sin}.scale(dbPosScale),
			data: DBData{triangle: vstripe},
		},
	}
	clients := []Client{
		{Position{-1.6, sin + 0.25}.scale(dbPosScale)},
		{Position{-1.6, sin - 0.25}.scale(dbPosScale)},
		{Position{-1.6, -sin}.scale(dbPosScale)},
		{Position{2.1, 0.25}.scale(dbPosScale)},
		{Position{2.1, -0.25}.scale(dbPosScale)},
	}
	channels := []Channel{
		{ep1: &clients[0], ep2: &dbs[1]},
		{ep1: &clients[1], ep2: &dbs[1]},
		{ep1: &clients[2], ep2: &dbs[2]},
		{ep1: &clients[3], ep2: &dbs[0]},
		{ep1: &clients[4], ep2: &dbs[0]},
	}
	emitters := compose_emitters()
	return Simulation{
		tick:      0,
		databases: dbs,
		clients:   clients,
		channels:  channels,
		events:    emitters,
		active:    false,
	}
}

func Sim3() Simulation {
	dbs := []Database{
		{
			pos:  Position{x: -0.3, y: 0.0},
			data: defaultDBData(),
		},
		{
			pos:  Position{x: 0.3, y: 0.0},
			data: defaultDBData(),
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
	return Simulation{
		tick:      0,
		databases: dbs,
		clients:   clients,
		channels:  channels,
		events:    emitters,
		active:    false,
	}
}

func Sim4() Simulation {
	dbs := []Database{
		{
			pos:  Position{x: -0.3, y: 0.0},
			data: defaultDBData(),
		},
		{
			pos:  Position{x: 0.3, y: 0.0},
			data: defaultDBData(),
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
			period:     500,
			emit:       ChannelEmitter{c: &channels[0], outgoing: true, sendee: randShapeWithNewStyle},
		},
	)
	return Simulation{
		tick:      0,
		databases: dbs,
		clients:   clients,
		channels:  channels,
		events:    emitters,
		active:    false,
	}
}

func Sim5() Simulation {
	dbs := []Database{
		{
			pos:  Position{x: -0.3, y: 0.0},
			data: defaultDBData(),
		},
		{
			pos:  Position{x: 0.3, y: 0.0},
			data: defaultDBData(),
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
		{ep1: &dbs[0], ep2: &dbs[1], travelTime: 1.5},
	}
	period := uint(700)
	emitters := compose_emitters(
		PeriodicEmitter{
			first_tick: 100,
			period:     period,
			emit:       ChannelEmitter{c: &channels[0], outgoing: true, sendee: independent(sqFunc, solidFunc)},
		},
		PeriodicEmitter{
			first_tick: 300,
			period:     period,
			emit:       ChannelEmitter{c: &channels[4], outgoing: true, sendee: independent(sqFunc, solidFunc)},
		},
		PeriodicEmitter{
			first_tick: 100,
			period:     period,
			emit:       ChannelEmitter{c: &channels[2], outgoing: true, sendee: independent(sqFunc, vstripeFunc)},
		},
		PeriodicEmitter{
			first_tick: 300,
			period:     period,
			emit:       ChannelEmitter{c: &channels[4], outgoing: false, sendee: independent(sqFunc, vstripeFunc)},
		},
	)
	return Simulation{
		tick:      0,
		databases: dbs,
		clients:   clients,
		channels:  channels,
		events:    emitters,
		active:    false,
	}
}

func Sim6() Simulation {
	scale := 0.4
	cos := math.Cos(math.Pi / 6.0)
	dbs := []Database{
		{
			pos:    Position{0.0, -0.1},
			data:   defaultDBData(),
			leader: true,
		},
		{
			pos:  Position{-cos * scale, -0.7},
			data: defaultDBData(),
		},
		{
			pos:  Position{cos * scale, -0.7},
			data: defaultDBData(),
		},
	}
	clients := []Client{
		{pos: Position{0.0, 0.5}},
		{pos: Position{-0.3, 0.5}},
		{pos: Position{0.3, 0.5}},
	}
	channels := []Channel{
		{ep1: &clients[0], ep2: &dbs[0], travelTime: 2.0},
		{ep1: &clients[1], ep2: &dbs[0], travelTime: 2.0},
		{ep1: &clients[2], ep2: &dbs[0], travelTime: 2.0},
		{ep1: &dbs[0], ep2: &dbs[1], travelTime: 1.5},
		{ep1: &dbs[0], ep2: &dbs[2], travelTime: 1.5},
	}
	emitters := compose_emitters()
	return Simulation{
		tick:      0,
		databases: dbs,
		clients:   clients,
		channels:  channels,
		events:    emitters,
		active:    false,
	}
}
