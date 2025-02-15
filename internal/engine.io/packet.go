package engineio

import (
	"encoding/json"
	"errors"
)

type PacketType string

const (
	PacketOpenType    = "0"
	PacketCloseType   = "1"
	PacketPingType    = "2"
	PacketPongType    = "3"
	PacketMessageType = "4"
	PacketUpgradeType = "5"
	PacketNoopType    = "6"
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

func (p *Packet) Encode() ([]byte, error) {
	buffer, err := p.rawDataToBuffer()

	if err != nil {
		return nil, err
	}

	return append([]byte(p.Type), buffer...), nil
}

func DecodePacket(rawData []byte) (*Packet, error) {
	return nil, nil
}

func (p *Packet) rawDataToBuffer() ([]byte, error) {
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
