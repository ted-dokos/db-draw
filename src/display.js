import * as vector from './vector.js';
import { createSolid, createHorizStripes, createVertStripes, createQuestionMark } from './patterns.js';
import { unitCircle, unitSquare, unitTriangle } from './shapes.js';

const CANVAS_WIDTH = 1920 - 3 * 160;
const CANVAS_HEIGHT = 1080 - 3 * 90;

let canvases = document.getElementsByTagName("canvas");
for (let i = 0; i < canvases.length; i++) {
    let canvas = canvases[i];
    canvas.setAttribute('width', CANVAS_WIDTH.toString());
    if (canvas.getAttribute('height') === null) {
        canvas.setAttribute('height', CANVAS_HEIGHT.toString());
    }
}

// Using the coordinate system from the golang side:
// treat the canvas as the space [-max_x, max_x] X [-max_y, max_y].
const REFERENCE_WIDTH = 1024.0;
const GO_MAX_X = 2.0;

const CLIENT_RADIUS = 25 * CANVAS_WIDTH / REFERENCE_WIDTH;
const CLIENT_BORDER = 3.0 * CANVAS_WIDTH / REFERENCE_WIDTH;
const FONT_SIZE_EM = (6.4 / 9.0) * CANVAS_WIDTH / REFERENCE_WIDTH;

function toCanvasCoords(pos, canvas) {
    const GO_MAX_Y = GO_MAX_X * canvas.height / CANVAS_WIDTH;
    return new vector.Vector(
        (pos.x + GO_MAX_X) * (CANVAS_WIDTH / (2.0 * GO_MAX_X)),
        (-pos.y + GO_MAX_Y) * (canvas.height / (2.0 * GO_MAX_Y)));
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

function getStylePattern(context, request) {
    switch (request) {
        case 0: // write request, solid
            return createSolid(context, "#FFC107");;
        case 1: // write request, horiz stripes
            return createHorizStripes(context);;
        case 2: // write request, stripes
            return createVertStripes(context);;
        default: // read request
            return createQuestionMark(context);
    }
}

function drawShapeWithStyle(context, pos, scale, shape_path, style_pattern) {
    context.save();
    let transform = { a: scale, d: scale, e: pos.x, f: pos.y };
    let pattern_transform = { a: 1 / scale, d: 1 / scale, e: -0.4, f: -0.7 };
    context.setTransform(transform);
    style_pattern.setTransform(pattern_transform);
    context.fillStyle = style_pattern;
    context.fill(shape_path);
    context.lineWidth = 2.0 / scale;
    context.stroke(shape_path);
    context.restore();
}

const DB_WIDTH = GO_MAX_X * 37.5 * CANVAS_WIDTH / REFERENCE_WIDTH;
const DB_HEIGHT = GO_MAX_X * 50 * CANVAS_WIDTH / REFERENCE_WIDTH;
const DB_BORDER = GO_MAX_X * 1.5 * CANVAS_WIDTH / REFERENCE_WIDTH;

function drawText(text, ctx, pos, scale) {
    ctx.save();
    ctx.font = (scale * FONT_SIZE_EM).toString() + "em Arial";
    ctx.textAlign = "center";
    ctx.textBaseline = "middle";
    ctx.fillText(text, pos.x, pos.y);
    ctx.restore();
}

function drawDB(db, ctx) {
    ctx.save();
    let db_pos = toCanvasCoords(db.pos, ctx.canvas);
    ctx.fillRect(db_pos.x - DB_WIDTH / 2.0, db_pos.y - DB_HEIGHT / 2.0, DB_WIDTH, DB_HEIGHT);
    ctx.clearRect(db_pos.x - DB_WIDTH / 2.0 + DB_BORDER,
        db_pos.y - DB_HEIGHT / 2.0 + DB_BORDER,
        DB_WIDTH - 2 * DB_BORDER,
        DB_HEIGHT - 2 * DB_BORDER);

    for (const shape in db.data) {
        if (db.data[shape] !== ShapeState.absent) {
            let pos = vector.add(db_pos, vector.smult(DB_HEIGHT / 3 * (ShapeMap[shape] - 1), new vector.Vector(0, 1)));
            const TRANSFORM_SCALE = 12.0 * CANVAS_WIDTH / REFERENCE_WIDTH;
            drawShapeWithStyle(ctx, pos, TRANSFORM_SCALE, getShapePath(ShapeMap[shape]), getStylePattern(ctx, db.data[shape]));
        }
    }
    ctx.restore();
    if (db.leader) {
        drawText(
            "👑",
            ctx, 
            vector.add(db_pos, new vector.Vector(0, -(DB_HEIGHT/2 + DB_BORDER))),
            2.0);
    }
}

function drawClient(client, ctx) {
    let client_pos = toCanvasCoords(client.pos, ctx.canvas);
    ctx.beginPath();
    ctx.arc(client_pos.x, client_pos.y, CLIENT_RADIUS - CLIENT_BORDER / 2.0, 0, 2 * Math.PI);
    ctx.lineWidth = CLIENT_BORDER;
    ctx.stroke();
}

function getPosFromEndpoint(sim, ep) {
    return ep.pos;
}

function drawStatusCode(ctx, pos, status) {
    let emoji = status === 0 ? "✅" : "⛔";
    drawText(emoji, ctx, pos, 1.0);
}

function drawChannelPacket(packet, ctx, pos) {
    ctx.save();
    const TRANSFORM_SCALE = 12.0 * CANVAS_WIDTH / REFERENCE_WIDTH;
    const READ_REQUEST_PATTERN = 3;
    if (packet.hasOwnProperty("readRequest")) {
        drawShapeWithStyle(
            ctx, pos, TRANSFORM_SCALE, getShapePath(packet.readRequest.shape),
            getStylePattern(ctx, READ_REQUEST_PATTERN));
    } else if (packet.hasOwnProperty("writeRequest")) {
        drawShapeWithStyle(
            ctx, pos, TRANSFORM_SCALE, getShapePath(packet.writeRequest.shape),
            getStylePattern(ctx, packet.writeRequest.writeState));
    } else if (packet.hasOwnProperty("readResponse")) {
        drawShapeWithStyle(
            ctx, pos, TRANSFORM_SCALE, getShapePath(packet.readResponse.shape),
            getStylePattern(ctx, packet.readResponse.state));
    } else if (packet.hasOwnProperty("writeResponse")) {
        drawStatusCode(ctx, pos, packet.writeResponse.status);
    }
    ctx.restore();
}

function getEndpointBuffer(endpoint) {
    let isDb = endpoint.hasOwnProperty('data');
    if (isDb) {
        return 1.2 * DB_HEIGHT / 2;
    } else {
        return 1.4 * CLIENT_RADIUS;
    }
}

function drawChannels(sim, ctx) {
    ctx.beginPath();
    for (let i = 0; i < sim.channels.length; i++) {
        let ch = sim.channels[i];
        let pos1 = toCanvasCoords(getPosFromEndpoint(sim, ch.ep1), ctx.canvas);
        let pos2 = toCanvasCoords(getPosFromEndpoint(sim, ch.ep2), ctx.canvas);
        let vec = vector.normalize({ x: pos2.x - pos1.x, y: pos2.y - pos1.y });
        let buffer1 = getEndpointBuffer(ch.ep1);
        let pos3 = vector.add(pos1, vector.smult(buffer1, vec));
        let buffer2 = getEndpointBuffer(ch.ep2);
        let pos4 = vector.add(pos2, vector.smult(-buffer2, vec));
        ctx.moveTo(pos3.x, pos3.y);
        ctx.lineTo(pos4.x, pos4.y);
        ctx.lineWidth = 3.0 * CANVAS_WIDTH / REFERENCE_WIDTH;
        ctx.stroke();
        if (ch.outgoing !== null) {
            let dist = vector.distance(pos3, pos4);
            let linePos = vector.add(pos3, vector.smult(dist * ch.outgoing.progress, vec));
            const PERP_DIST = -20.0 * CANVAS_WIDTH / REFERENCE_WIDTH;
            let perp = vector.getPerp(vec);
            if (perp.y < 0) {
                perp = vector.smult(-1, perp);
            }
            let packetPos = vector.add(linePos, vector.smult(PERP_DIST, perp));
            drawChannelPacket(ch.outgoing, ctx, packetPos);
        }
        if (ch.incoming !== null) {
            let dist = vector.distance(pos3, pos4);
            let linePos = vector.add(pos4, vector.smult(-dist * ch.incoming.progress, vec));
            const PERP_DIST = -20.0 * CANVAS_WIDTH / REFERENCE_WIDTH;
            let perp = vector.getPerp(vec);
            if (perp.y >= 0) {
                perp = vector.smult(-1, perp);
            }
            let packetPos = vector.add(linePos, vector.smult(PERP_DIST, perp));
            drawChannelPacket(ch.incoming, ctx, packetPos);
        }
    }
}

const go = new Go();
go.importObject.howdy = {
    JsDo: () => {
        let sims = window.getSims();
        sims.map((sim, idx) => {
            if (!sim.active) {
                return;
            }
            let canvas = canvases[idx];
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
    }
};
WebAssembly.instantiateStreaming(fetch("bin/main.wasm"),
    go.importObject).then((result) => {
        go.run(result.instance);
    });

function setSimIndex(idx) {
    if (!window.setSimIndex) {
        setTimeout(() => { setSimIndex(idx); }, 10);
        return;
    }
    current_sim_running = idx;
    window.setSimIndex(idx);
}

function observeCanvases() {
    for (let i = 0; i < canvases.length; i++) {
        let canvas = canvases[i];
        let intersection_observer = new IntersectionObserver(
            (entries) => {
                let entry = entries[0];
                if (entry.intersectionRatio === 1.0 && window.activateSim) {
                    window.activateSim(i);
                    return;
                } else if (window.disableSim) {
                    window.disableSim(i);
                    return;
                }
            },
            /*options=*/ {
                root: document.querySelector("#scrollArea"),
                rootMargin: "0px",
                threshold: [1.0, 0.95],
            });
        intersection_observer.observe(canvas);
    }
}

function drawInitialSimStates() {
    if (!window.setSimIndex) {
        setTimeout(drawInitialSimStates, 10);
        return;
    }
    for (let i = 0; i < canvases.length; i++) {
        setSimIndex(i);
        go.importObject.howdy.JsDo();
    }
    setSimIndex(-1);
}
drawInitialSimStates();
observeCanvases();