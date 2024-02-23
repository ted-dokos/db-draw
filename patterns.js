export function createHorizStripes(context) {
    let stripes = new OffscreenCanvas(6, 6);
    let ctx = stripes.getContext("2d");
    ctx.fillStyle = "#D81B60";
    ctx.fillRect(0, 0, 6, 3);
    return context.createPattern(stripes, "repeat");
}

export function createVertStripes(context) {
    let stripes = new OffscreenCanvas(6, 6);
    let ctx = stripes.getContext("2d");
    ctx.fillStyle = "#1E88E5";
    ctx.fillRect(0, 0, 3, 6);
    return context.createPattern(stripes, "repeat");
}

export function createSolid(context) {
    let solid = new OffscreenCanvas(1, 1);
    let ctx = solid.getContext("2d");
    ctx.fillStyle = "#FFC107";
    ctx.fillRect(0, 0, 1, 1);
    return context.createPattern(solid, "repeat");
}