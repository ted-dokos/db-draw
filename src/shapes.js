export function unitTriangle() {
    let tri = new Path2D();
    tri.moveTo(0.0, -1.0);
    let x = 2 / Math.sqrt(3.0);
    tri.lineTo(x, 1.0);
    tri.lineTo(-x, 1.0);
    tri.lineTo(0.0, -1.0);
    return tri;
}

export function unitCircle() {
    let circ = new Path2D();
    circ.arc(0, 0, 1.0, 0, 2 * Math.PI);
    return circ;
}

export function unitSquare() {
    let sq = new Path2D();
    sq.moveTo(1, 1);
    sq.lineTo(-1, 1);
    sq.lineTo(-1, -1);
    sq.lineTo(1, -1);
    sq.closePath();
    return sq;
}