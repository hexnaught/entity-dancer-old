let packetId = 0;
const dgram = require("dgram");

const client = dgram.createSocket("udp4");

// setTimeout(() => {
//   client.close();
// }, 30000);

const keepAliveMsg = new Uint8Array([2, ...Buffer.from("Hello World")]);

// console.log(keepAliveMsg);

client.send(keepAliveMsg, 8989, (err) => {
  /* console.log(`Sent!`) */
});
setInterval(() => {
  client.send(keepAliveMsg, 8989, (err) => {
    /* console.log(`Sent!`) */
  });
}, 10000);

client.connect(8989, `localhost`, () => {
  `Connection Made`;
});

client.on("message", (buf, info) => {
  // Example info: { address: '127.0.0.1', family: 'IPv4', port: 8989, size: 23 }
  // console.log(info);
  console.info(`Bytes: ${info.size} | Resp: ${Buffer.from(buf, 'hex')}`);

  /**
     The Below Will Read:

     type S struct {
       PType int8
       X     uint16
       Y     uint16
       Z     uint16
     }
   */
  console.log(buf.readInt8(0))
  console.log(buf.readUInt16LE(1))
  console.log(buf.readUInt16LE(3))
  console.log(buf.readUInt16LE(5))
});
