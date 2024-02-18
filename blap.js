import { Subject } from 'rxjs';

const go = new Go();
const obs = new Subject();
obs.subscribe(sim =>{
    let elem = document.getElementById("count");
    elem.innerText = `{${sim.databases[0].pos.x}, ${sim.databases[0].pos.y}}`
});
obs.subscribe(sim => {
    let canvas = document.getElementById("draw");
    if (canvas.getContext) {
        const ctx = canvas.getContext("2d");
        ctx.clearRect(0, 0, canvas.width, canvas.height);
        let wr = 2.0;
        let hr = 1.5;
        for (var i = 0; i < sim.databases.length; i++) {
            let db = sim.databases[i];
            let x = (db.pos.x + wr) * (canvas.width / (2.0 * wr));
            let y = (db.pos.y + hr) * (canvas.height / (2.0 * hr));
            ctx.fillRect(x, y, 25, 25);
        }
    }
});
//go.importObject.howdy = {JsDo: (db) => obs.next(db)};
go.importObject.howdy = {JsDo: () => {
    let sim = window.callback();
    obs.next(sim);
}};
WebAssembly.instantiateStreaming(fetch("main.wasm"),
    go.importObject).then((result) => {
    go.run(result.instance);
});
//const numbers = interval(1000);
//numbers.subscribe(x => console.log('Next: ', x));
