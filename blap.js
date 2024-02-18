import { Subject } from 'rxjs';

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

const go = new Go();
const obs = new Subject();
obs.subscribe(sim => {
    let canvas = document.getElementById("draw");
    if (canvas.getContext) {
        const ctx = canvas.getContext("2d");
        ctx.clearRect(0, 0, canvas.width, canvas.height);
        let wr = 2.0;
        let hr = 1.5;
        for (var i = 0; i < sim.databases.length; i++) {
            let db = sim.databases[i];
            drawDB(db, ctx);
        }
    }
});
go.importObject.howdy = {JsDo: () => {
    let sim = window.callback();
    obs.next(sim);
}};
WebAssembly.instantiateStreaming(fetch("main.wasm"),
    go.importObject).then((result) => {
    go.run(result.instance);
});
