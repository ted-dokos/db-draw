import { Subject } from 'rxjs';

const CANVAS = document.getElementById("draw");
const wr = 2.0;
const hr = wr * CANVAS.height / CANVAS.width;

function toCanvasCoords(pos) {
    return { x: (pos.x + wr) * (CANVAS.width / (2.0 * wr)),
             y: (pos.y + hr) * (CANVAS.height / (2.0 * hr)) };
}

function drawDB(db, ctx) {
    let wr = 2.0;
    let hr = wr * ctx.canvas.height / ctx.canvas.width;
    let x = (db.pos.x + wr) * (ctx.canvas.width / (2.0 * wr));
    let y = (db.pos.y + hr) * (ctx.canvas.height / (2.0 * hr));

    let db_width = wr * 37.5;
    let db_height = wr * 50;
    let db_border = wr * 1.5;
    ctx.fillRect(x, y, db_width, db_height);
    ctx.clearRect(x + db_border, y + db_border, db_width - 2 * db_border, db_height - 2 * db_border);
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
        ctx.lineWidth = wr;
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
