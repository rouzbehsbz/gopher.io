package engineio

type PacketType int

const (
	PacketOpenType = iota
	PacketCloseType
	PacketPingType
	PacketPongType
	PacketMessageType
	PacketUpgradeType
	PacketNoopType
)

type Packet struct {
	_type PacketType
	data  []byte
}
