const id = Math.random();
console.log(id)
let source = new EventSource(`http://172.16.2.10:9000/notify?id=${id}`);
// let source = new EventSource(`http://172.16.2.10:9000/stream`);

source.addEventListener("open",()=>{
    console.log("OPEN:", id)
});

// source.addEventListener("message",(event)=>{
//     console.log("message:", event.data)
// });

source.addEventListener("saludar",(event)=>{
    console.log("Saludar:", event.data)
});

source.addEventListener("saltar",(event)=>{
    console.log("Saltar:", event.data)
});