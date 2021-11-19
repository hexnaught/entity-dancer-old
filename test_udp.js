let packetId = 0;
const dgram = require("dgram");
const buf1 = Buffer.from("Some ");
const buf2 = Buffer.from("bytes ");
const client = dgram.createSocket("udp4");

const readline = require("readline");
const rl = readline.createInterface({
    input: process.stdin,
    output: process.stdout
});

// const join = new Uint8Array([0x00, ...Buffer.from("Hello World")]);
// const part = new Uint8Array([0x01, ...Buffer.from("Hello World")]);
// const keepAlive = new Uint8Array([0x02, ...Buffer.from("Hello World")]);
// const move = new Uint8Array([0x03, 0xff, 0x48, 0xff]);

const join = Buffer.from([0]);
const part = Buffer.from([1]);
const keepAlive = Buffer.from([2]);
const move = Buffer.from([0x03, 0xff, 0x48, 0xff], "binary");

// setTimeout(() => {
//   client.close();
// }, 30000);

// client.send(join, 8989, (err) => {});
// setTimeout(() => {}, 500)
// client.send(move, 8989, (err) => {});
client.send(keepAlive, 8989, (err) => {});
// client.send(part, 8989, (err) => {});

// client.send(join, 8989, (err) => {});
// setInterval(() => {
//   // console.info(`Sending`)
//   const buf3 = Buffer.from(`Packet ID: ${packetId++}`);
//   // client.send([buf1, buf2, buf3], 8989, (err) => { /* console.log(`Sent!`) */ });

//   client.send(move, 8989, (err) => {});
// }, 500);
// // }, 1000 / 64) // (1000ms / tick rate)

client.on("message", (buf, info) => {
  // Example info: { address: '127.0.0.1', family: 'IPv4', port: 8989, size: 23 }
  // console.log(info);
  // console.info(`Bytes: ${info.size} | Resp: ${buf.toString('binary')}`);
  console.info(`Bytes: ${info.size} | Resp: ${buf}`);
});

const toNetworkBuffer = (str) => {
  return new Uint8Array(Buffer.from(str));
};

const fromNetworkBuffer = (nwBuf) => {
  return Buffer.from(nwBuf.buffer, `binary`).toString();
};
