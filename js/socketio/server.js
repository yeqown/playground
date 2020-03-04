const express = require('express')
const app = express()
const http = require('http').Server(app)
const io = require('socket.io')(http, {
    path: "/socket.io",
    transports: ["polling", "websocket"],
});

const port = process.env.PORT || 3000;

var yunfanNsp = io.of("/yunfan")
var socketsPool = {}

// yunfan nsp 在连接时处理链接到房间
yunfanNsp.on("connection", socket => {
    console.log("connected, ", socket.id);
    socket.on("disconnect", () => {
        console.log("disconnect");
    })

    socket.on("join", data => {
        socket.join(data.room) // room or rooms
        console.log(socket.id, " joined ", data.room);
    })

    socket.on("disconnect", () => {
        console.log("disc");
    })

    socketsPool[socket.id] = socket
    // console.log(socketsPool);
})

yunfanNsp.on("disconnect", socket => {
    console.log("nsp disconnected, ", socket.id);
})

// 单发
app.get("/to_one", (req, res) => {
    let socketId = req.query["sid"]
    if (socketId === '') {
        res.send("invalid sid")
        return
    }
    let msg = req.query["message"] || "empty message"
    socketId = yunfanNsp.name + "#" + socketId
    console.log("socketId: ", socketId);

    let socket = socketsPool[socketId]
    if (socket) {
        socket.send(msg)
        res.send({ code: 0, message: "ok" })
    } else {
        res.send({ code: 1, message: "socket undefined" })
    }
})

// 房间内发送
app.get("/broadcast_room", (req, res) => {
    let roomId = req.query["room_id"]
    if (roomId === '') {
        res.send("invalid room_id")
        return
    }
    console.log("roomId: ", roomId);
    let msg = req.query["message"] || "empty message"
    yunfanNsp.in(roomId).send(msg)
    res.send({ code: 0, message: "ok" })
})

// namespace下 全局发送
app.get("/broadcast", (req, res) => {
    let msg = req.query["message"] || "empty message"
    yunfanNsp.send(msg)
    res.send({ code: 0, message: "ok" })
})


http.listen(port, () => console.log('listening on port ' + port));