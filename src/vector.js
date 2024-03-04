export class Vector {
    constructor(x, y) {
        this.x = x;
        this.y = y;
    }
}

// Should be rotated 90 degrees clockwise.
export function getPerp(v) {
    return new Vector(-v.y, v.x);
}

export function normalize(v) {
    let length = Math.sqrt(v.x * v.x + v.y * v.y);
    return new Vector(v.x / length, v.y / length);
}

export function add(v1, v2) {
    return new Vector(v1.x + v2.x, v1.y + v2.y);
}

export function smult(scalar, v) {
    return new Vector(scalar * v.x, scalar * v.y);
}

export function distance(v1, v2) {
    let dx = v1.x - v2.x;
    let dy = v1.y - v2.y;
    return Math.sqrt(dx * dx + dy * dy);
}