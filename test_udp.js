let packetId = 0;
const dgram = require('dgram');
const buf1 = Buffer.from('Some ');
const buf2 = Buffer.from('bytes ');

const client = dgram.createSocket('udp4');

setTimeout(() => {
  client.close();
}, 30000)

setInterval(() => {
  // console.info(`Sending`)
  const buf3 = Buffer.from(`Packet ID: ${packetId++}`)
  client.send([buf1, buf2, buf3], 8989, (err) => { /* console.log(`Sent!`) */ });
}, 1000 / 64) // (1000ms / tick rate)

client.on("message", (buf, info) => {
  // Example info: { address: '127.0.0.1', family: 'IPv4', port: 8989, size: 23 }
  console.log(info)
  console.info(`Bytes: ${info.size} | Resp: ${buf}`)
})
