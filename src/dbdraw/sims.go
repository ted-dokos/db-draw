package dbdraw

// Same order as they appear on the page.

func Sim1() SimulationState {
	dbs := []Database{{pos: Position{x: 0.5, y: 0.0}, data: map[Shape]RequestType{}}}
	clients := []Client{{pos: Position{x: -0.5, y: 0.0}}}
	chs := []Channel{{ep1: &clients[0], ep2: &dbs[0], travelTime: 2.0}}
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

func Sim2() SimulationState {
	dbs := []Database{{
		pos: Position{x: 0.0, y: 0.25},
		data: map[Shape]RequestType{
			circle:   solid,
			square:   hstripe,
			triangle: vstripe}},
	}
	clients := []Client{
		{Position{-0.5, -0.25}},
		{Position{0.5, -0.25}},
	}
	channels := []Channel{
		{ep1: &clients[0], ep2: &dbs[0], travelTime: 2.0},
		{ep1: &clients[1], ep2: &dbs[0], travelTime: 2.0},
	}
	emitters := compose_emitters(
		PeriodicEmitter{
			first_tick: 100,
			period:     300,
			emit:       ChannelEmitter{c: &channels[0], outgoing: true, sendee: randShapeWithNewStyle},
		},
		PeriodicEmitter{
			first_tick: 200,
			period:     300,
			emit:       ChannelEmitter{c: &channels[1], outgoing: true, sendee: randShapeWithNewStyle},
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
