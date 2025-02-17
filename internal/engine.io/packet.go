package engineio

import (
	"bytes"
	"encoding/json"
	"errors"
)

type PacketType string

const SeperatorCharacter = 0x1e

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

func (p *Packet) Encode() ([]byte, error) {
	buffer, err := p.rawDataToBuffer()

	if err != nil {
		return nil, err
	}

	return append([]byte(p.Type), buffer...), nil
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

func BufferToRawDate(buffer []byte) any {
	var v any

	if err := json.Unmarshal(buffer, &v); err == nil {
		return v
	}

	if str := string(buffer); str != "" {
		return str
	}

	return buffer
}

func EncodePackets(packets []Packet) ([]byte, error) {
	var bytes []byte

	for i, packet := range packets {
		encodedPacket, err := packet.Encode()

		if err != nil {
			return nil, err
		}

		bytes = append(bytes, encodedPacket...)

		if i < len(packets)-1 {
			bytes = append(bytes, SeperatorCharacter)
		}
	}

	return bytes, nil
}

func DecodePacket(rawData []byte) (Packet, error) {
	if len(rawData) == 0 {
		return Packet{}, errors.New("empty packet data")
	}

	packetType := PacketType(rawData[0])
	data := rawData[1:]

	return Packet{
		Type:    packetType,
		RawData: BufferToRawDate(data),
	}, nil
}

func DecodePackets(rawData []byte) ([]Packet, error) {
	var packets []Packet

	for _, packetData := range bytes.Split(rawData, []byte{SeperatorCharacter}) {
		packet, err := DecodePacket(packetData)

		if err != nil {
			return nil, err
		}

		packets = append(packets, packet)
	}

	return packets, nil
}
