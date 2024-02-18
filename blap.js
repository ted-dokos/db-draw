import { Subject } from 'rxjs';

const go = new Go();
const obs = new Subject();
obs.subscribe(db=>{
    let elem = document.getElementById("count");
    elem.innerText = `{${db.pos.x}, ${db.pos.y}}`
});
//go.importObject.howdy = {JsDo: (db) => obs.next(db)};
go.importObject.howdy = {JsDo: () => {
    let db = window.callback();
    obs.next(db);
}};
WebAssembly.instantiateStreaming(fetch("main.wasm"),
    go.importObject).then((result) => {
    go.run(result.instance);
});
//const numbers = interval(1000);
//numbers.subscribe(x => console.log('Next: ', x));
