package engineio

import (
	"encoding/base64"
	"encoding/json"
	"errors"
)

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
	Type    PacketType
	RawData any
}

func NewPacket(_type PacketType, rawData any) *Packet {
	return &Packet{
		Type:    _type,
		RawData: rawData,
	}
}

func (p *Packet) EncodePacket(isBinarySupported bool) ([]byte, error) {
	buffer, err := p.toBuffer()

	if err != nil {
		return nil, err
	}

	if isBinarySupported {
		return buffer, nil
	}

	return append([]byte("b"), []byte(base64.StdEncoding.EncodeToString(buffer))...), nil
}

func (p *Packet) toBuffer() ([]byte, error) {
	switch v := p.RawData.(type) {
	case []byte:
		return v, nil
	case string:
		return []byte(v), nil
	default:
		b, err := json.Marshal(v)

		if err != nil {
			return nil, errors.New("can't marshal to json encoding.")
		}

		return b, nil
	}
}
