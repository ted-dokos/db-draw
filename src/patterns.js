let horizStripes = new OffscreenCanvas(6, 6);
let horizStripesContext = horizStripes.getContext("2d");
horizStripesContext.fillStyle = "#D81B60";
horizStripesContext.fillRect(0, 1, 6, 4);
export function createHorizStripes(context) {
    return context.createPattern(horizStripes, "repeat");
}

let vertStripes = new OffscreenCanvas(6, 6);
let vertStripesContext = vertStripes.getContext("2d");
vertStripesContext.fillStyle = "#1E88E5";
vertStripesContext.fillRect(2, 0, 5, 6);
export function createVertStripes(context) {
    return context.createPattern(vertStripes, "repeat");
}

let solid = new OffscreenCanvas(1, 1);
let solidContext = solid.getContext("2d");
export function createSolid(context, color) {
    solidContext.fillStyle = color;
    solidContext.fillRect(0, 0, 1, 1);
    return context.createPattern(solid, "repeat");
}

let qsize = 200;
let question = new OffscreenCanvas(qsize, qsize);
let questionContext = question.getContext("2d");
questionContext.font = "25px Arial";
questionContext.fillText("?", 0, 25);
export function createQuestionMark(context) {
    return context.createPattern(question, "repeat");
}