package protocol

type Packet struct {
	Type PacketType
	NWID string
}

type MovePacket struct {
	x, y, z uint16
}

type Vector3 struct {
	x, y, z int16
}

type Quaternion struct {
	x, y, z, w int16
}

type EntityData struct {
	NWID      []byte     // Unique identifier
	Postition Vector3    // Pos
	Rotation  Quaternion // Rot
}

func (p *Packet) FromBinary(data []byte) {
	p.Type = PacketType(data[0])

}

func (p *Packet) ToBinary() []byte {

	return []byte{}
}

// NWID SPEC
// (server_id & 255) << 56
// + (server_startup_time_in_seconds << 24)
// + incremented_variable++;
