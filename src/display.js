import { Subject } from 'rxjs';
import * as vector from './vector.js';
import { createSolid, createHorizStripes, createVertStripes } from './patterns.js';
import { unitCircle, unitSquare, unitTriangle } from './shapes.js';

const CANVAS = document.getElementById("draw");
// Using the coordinate system from the golang side:
// treat the canvas as the space [-max_x, max_x] X [-max_y, max_y].
const REFERENCE_WIDTH = 1024.0;
const GO_MAX_X = 2.0;
const GO_MAX_Y = GO_MAX_X * CANVAS.height / CANVAS.width;

const CLIENT_RADIUS = 25 * CANVAS.width / REFERENCE_WIDTH;
const CLIENT_BORDER = 3.0 * CANVAS.width / REFERENCE_WIDTH;

const SOLID = createSolid(CANVAS.getContext("2d"), "#FFC107");
const BLANK = createSolid(CANVAS.getContext("2d"), "#FFFFFF");
const HORIZ_STRIPES = createHorizStripes(CANVAS.getContext("2d"));
const VERT_STRIPES = createVertStripes(CANVAS.getContext("2d"));

function toCanvasCoords(pos) {
    return new vector.Vector(
        (pos.x + GO_MAX_X) * (CANVAS.width / (2.0 * GO_MAX_X)),
        (pos.y + GO_MAX_Y) * (CANVAS.height / (2.0 * GO_MAX_Y)));
}

const ShapeMap = {
    'circle': 0,
    'square': 1,
    'triangle': 2,
};

const ShapeState = {
    'solid': 0,
    'hstripe': 1,
    'vstripe': 2,
    'absent': 3,
};

function getShapePath(shape) {
    switch (shape) {
        case 0: // circle
            return unitCircle();
        case 1: // square
            return unitSquare();
        default: // triangle
            return unitTriangle();
    }
}

function getStylePattern(style) {
    switch (style) {
        case 0: // solid
            return SOLID;
        case 1: // horiz stripes
            return HORIZ_STRIPES;
        case 2: // vert stripes
            return VERT_STRIPES;
        default: // absent
            return BLANK;
    }
}

function drawShapeWithStyle(context, pos, scale, shape_path, style_pattern) {
    let transform = { a: scale, d: scale, e: pos.x, f: pos.y };
    let pattern_transform = { a: 1 / scale, d: 1 / scale, e: 1};
    context.setTransform(transform);
    style_pattern.setTransform(pattern_transform);
    context.fillStyle = style_pattern;
    context.fill(shape_path);
    context.lineWidth = 2.0 / scale;
    context.stroke(shape_path);
}

function drawDB(db, ctx) {
    const DB_WIDTH = GO_MAX_X * 37.5 * CANVAS.width / REFERENCE_WIDTH;
    const DB_HEIGHT = GO_MAX_X * 50 * CANVAS.width / REFERENCE_WIDTH;
    const DB_BORDER = GO_MAX_X * 1.5 * CANVAS.width / REFERENCE_WIDTH;
    let db_pos = toCanvasCoords(db.pos);
    ctx.fillRect(db_pos.x - DB_WIDTH / 2.0, db_pos.y - DB_HEIGHT / 2.0, DB_WIDTH, DB_HEIGHT);
    ctx.clearRect(db_pos.x - DB_WIDTH / 2.0 + DB_BORDER,
        db_pos.y - DB_HEIGHT / 2.0 + DB_BORDER,
        DB_WIDTH - 2 * DB_BORDER,
        DB_HEIGHT - 2 * DB_BORDER);

    ctx.save();
    for (const shape in db.data) {
        if (db.data[shape] !== ShapeState.absent) {
            let pos = vector.add(db_pos, vector.smult(DB_HEIGHT / 3 * (ShapeMap[shape] - 1), new vector.Vector(0, 1)));
            const TRANSFORM_SCALE = 12.0 * CANVAS.width / REFERENCE_WIDTH;
            drawShapeWithStyle(ctx, pos, TRANSFORM_SCALE, getShapePath(ShapeMap[shape]), getStylePattern(db.data[shape]));
        }
    }
    ctx.restore();
}

function drawClient(client, ctx) {
    let client_pos = toCanvasCoords(client.pos);
    ctx.beginPath();
    ctx.arc(client_pos.x, client_pos.y, CLIENT_RADIUS - CLIENT_BORDER / 2.0, 0, 2 * Math.PI);
    ctx.lineWidth = CLIENT_BORDER;
    ctx.stroke();
}

function getPosFromEndpoint(sim, ep) {
    if (ep.type === 'd') {
        return sim.databases[ep.index].pos;
    } else {
        return sim.clients[ep.index].pos;
    }
}

function drawChannelTransaction(transaction, ctx, pos) {
    ctx.save();
    const TRANSFORM_SCALE = 12.0 * CANVAS.width / REFERENCE_WIDTH;
    drawShapeWithStyle(ctx, pos, TRANSFORM_SCALE, getShapePath(transaction.shape), getStylePattern(transaction.style));
    ctx.restore();
}

function drawChannels(sim, ctx) {
    ctx.beginPath();
    for (let i = 0; i < sim.channels.length; i++) {
        let ch = sim.channels[i];
        let pos1 = toCanvasCoords(getPosFromEndpoint(sim, ch.ep1));
        let pos2 = toCanvasCoords(getPosFromEndpoint(sim, ch.ep2));
        let vec = vector.normalize({ x: pos2.x - pos1.x, y: pos2.y - pos1.y });
        let buffer = 1.4 * CLIENT_RADIUS;
        let pos3 = vector.add(pos1, vector.smult(buffer, vec));
        let pos4 = vector.add(pos2, vector.smult(-buffer, vec));
        ctx.moveTo(pos3.x, pos3.y);
        ctx.lineTo(pos4.x, pos4.y);
        ctx.lineWidth = 3.0 * CANVAS.width / REFERENCE_WIDTH;
        ctx.stroke();
        if (ch.outgoing !== null) {
            let dist = vector.distance(pos3, pos4);
            let linepos = vector.add(pos3, vector.smult(dist * ch.outgoing.progress, vec));
            const PERP_DIST = -20.0 * CANVAS.width / REFERENCE_WIDTH;
            let trpos = vector.add(linepos, vector.smult(PERP_DIST, vector.get_perp(vec)));
            drawChannelTransaction(ch.outgoing, ctx, trpos);
        }
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
            drawDB(sim.databases[i], ctx);
        }
        for (let i = 0; i < sim.clients.length; i++) {
            drawClient(sim.clients[i], ctx);
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
WebAssembly.instantiateStreaming(fetch("bin/main.wasm"),
    go.importObject).then((result) => {
        go.run(result.instance);
    });