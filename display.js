import { Subject } from 'rxjs';

const CANVAS = document.getElementById("draw");
// Using the coordinate system from the golang side:
// treat the canvas as the space [-max_x, max_x] X [-max_y, max_y].
const GO_MAX_X = 2.0;
const GO_MAX_Y = GO_MAX_X * CANVAS.height / CANVAS.width;

function toCanvasCoords(pos) {
    return { x: (pos.x + GO_MAX_X) * (CANVAS.width / (2.0 * GO_MAX_X)),
             y: (pos.y + GO_MAX_Y) * (CANVAS.height / (2.0 * GO_MAX_Y)) };
}

function drawDB(db, ctx) {
    let x = (db.pos.x + GO_MAX_X) * (ctx.canvas.width / (2.0 * GO_MAX_X));
    let y = (db.pos.y + GO_MAX_Y) * (ctx.canvas.height / (2.0 * GO_MAX_Y));

    const DB_WIDTH = GO_MAX_X * 37.5 * 1000.0 / CANVAS.width;
    const DB_HEIGHT = GO_MAX_X * 50 * 1000.0 / CANVAS.width;
    const DB_BORDER = GO_MAX_X * 1.5 * 1000.0 / CANVAS.width;
    ctx.fillRect(x - DB_WIDTH / 2.0, y - DB_HEIGHT / 2.0, DB_WIDTH, DB_HEIGHT);
    ctx.clearRect(x - DB_WIDTH / 2.0 + DB_BORDER,
                  y - DB_HEIGHT / 2.0 + DB_BORDER,
                  DB_WIDTH - 2 * DB_BORDER,
                  DB_HEIGHT - 2 * DB_BORDER);
}

function getPosFromEndpoint(sim, ep) {
    if (ep.type === 'd') {
        return sim.databases[ep.index].pos;
    } else {
        return sim.clients[ep.index].pos;
    }
}

function drawChannels(sim, ctx) {
    ctx.beginPath();
    for (let i = 0; i < sim.channels.length; i++) {
        let ch = sim.channels[i];
        let pos1 = toCanvasCoords(getPosFromEndpoint(sim, ch.ep1));
        let pos2 = toCanvasCoords(getPosFromEndpoint(sim, ch.ep2));
        ctx.moveTo(pos1.x, pos1.y);
        ctx.lineTo(pos2.x, pos2.y);
        ctx.lineWidth = 3000.0 / CANVAS.width;
        ctx.stroke();
    }
}

const go = new Go();
const obs = new Subject();
obs.subscribe(sim => {
    let canvas = document.getElementById("draw");
    if (canvas.getContext) {
        const ctx = canvas.getContext("2d");
        ctx.clearRect(0, 0, canvas.width, canvas.height);
        let wr = 2.0;
        let hr = 1.5;
        for (let i = 0; i < sim.databases.length; i++) {
            let db = sim.databases[i];
            drawDB(db, ctx);
        }
        drawChannels(sim, ctx);
    }
});
go.importObject.howdy = {
    JsDo: () => {
        let sim = window.callback();
        obs.next(sim);
    }
};
WebAssembly.instantiateStreaming(fetch("main.wasm"),
    go.importObject).then((result) => {
        go.run(result.instance);
    });
