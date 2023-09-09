const Can = require('@csllc/cs-pcan-usb');
const express = require('express');
const app = express();
const http = require('http');
const server = http.createServer(app);
const { Server } = require("socket.io");
const io = new Server(server);
const port = 5000;

app.get('/', function (req, res) {
    res.sendFile('index.html', { root: __dirname });
});

server.listen(port, () => {
    console.log(`A browser window should open, or manually visit http://localhost:${port} from this machine.`);

    var url = `http://localhost:${port}`;
    var start = (process.platform == 'darwin'? 'open': process.platform == 'win32'? 'start': 'xdg-open');
    require('child_process').exec(start + ' ' + url);
});

let can = new Can({
    canRate: 250000,
});

can.list()
    .then((ports) => {
        return can.open(ports[0].path)
            .catch((err) => {
                console.error(err);
                process.exit(1);
            });
    })
    .catch((err) => {
        console.error(err);
        can.close();
        process.exit(1);
    });

can.on('data', function (msg) {
    var buffer64 = new ArrayBuffer(8);
    var buffer16 = new ArrayBuffer(2);
    var buffer8 = new ArrayBuffer(1);

    var view64 = new DataView(buffer64);
    var view16 = new DataView(buffer16);
    var view8 = new DataView(buffer8);

    switch (msg.id) {
        case 200:
            view64.setBigUint64(0, Buffer.from(msg.buf).readBigUint64LE(0));
            io.emit("gps lat", view64.getFloat64(0));
            break;
        case 201:
            view64.setBigUint64(0, Buffer.from(msg.buf).readBigUint64LE(0));
            io.emit("gps lon", view64.getFloat64(0));
            break;
        case 202:
            view64.setBigUint64(0, Buffer.from(msg.buf).readBigUint64LE(0));
            io.emit("gps alt", view64.getFloat64(0));
            break;
        case 203:
            view64.setBigUint64(0, Buffer.from(msg.buf).readBigUint64LE(0));
            io.emit("gps speed", view64.getFloat64(0));
            break;
        case 204:
            view8.setUint8(0, msg.buf[0]);
            io.emit("gps mode", view8.getUint8(0));

            view8.setUint8(0, msg.buf[1]);
            io.emit("gps status", view8.getUint8(0));

            view16.setUint16(0, Buffer.from(msg.buf.slice(2, 4)).readUint16LE(0));
            io.emit("gps nsat", view16.getUint16(0));

            view16.setUint16(0, Buffer.from(msg.buf.slice(4, 6)).readUint16LE(0));
            io.emit("gps usat", view16.getUint16(0));

            view8.setUint8(0, msg.buf[6]);
            io.emit("gps qual", view8.getUint8(0));
            break;
        case 205:
            view16.setUint16(0, Buffer.from(msg.buf.slice(0, 2)).readUint16LE(0));
            io.emit("motion x", view16.getInt16(0));

            view16.setUint16(0, Buffer.from(msg.buf.slice(2, 4)).readUint16LE(0));
            io.emit("motion y", view16.getInt16(0));

            view16.setUint16(0, Buffer.from(msg.buf.slice(4, 6)).readUint16LE(0));
            io.emit("motion z", view16.getInt16(0));

            view8.setUint8(0, msg.buf[6]);
            io.emit("motion scale", view8.getUint8(0));
            break;
    }
});