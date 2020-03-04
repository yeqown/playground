const socket = require("socket.io-client")("http://localhost:3000/yunfan", {
    path: "/socket.io"
})

// console.log(process.env);


var roomId = process.env.room
if (roomId === '') {
    throw Error("no roomId message")
}

socket.on('connect', function () {
    socket.emit("auth", { clientId: 'clientId to parse client info' })
    socket.emit("join", { room: roomId })
});

// socket.on('event', function (data) {
//     console.log(data);
// });
// socket.on('disconnect', function () { });

console.log("namespace: ", socket.nsp, "socketId: ", socket.id);
socket.on("message", data => {
    console.log(data);
})

// socket.emit("message", { you: "ok", id: 120389 }, () => {
//     console.log("recv server ack");
// })

socket.on("disconnect", () => {
    console.log("server disconnect");
    socket.close()
})
