package wispparse

import (
	"encoding/binary"
	"fmt"
)

// Implementing https://github.com/MercuryWorkshop/wisp-protocol/blob/main/protocol.md

type PacketType uint8

const (
	PacketTypeConnect PacketType = iota + 1
	PacketTypeData
	PacketTypeContinue
	PacketTypeClose
)

type Packet struct {
	// The type of packet
	Type PacketType
	// The stream ID
	StreamID uint32
	// The payload of the packet
	Payload []byte
}

func ParsePacket(b []byte) (*Packet, error) {
	if len(b) < 5 {
		return nil, fmt.Errorf("packet is too short")
	}
	return &Packet{
		Type:     PacketType(b[0]),
		StreamID: binary.LittleEndian.Uint32(b[1:5]),
		Payload:  b[5:],
	}, nil
}

func (p *Packet) ConnectPacket() (*ConnectPacket, error) {
	if p.Type != PacketTypeConnect {
		return nil, fmt.Errorf("packet is not a connect packet")
	}
	var cp ConnectPacket
	cp.StreamType = StreamType(p.Payload[0])
	cp.DestinationPort = binary.LittleEndian.Uint16(p.Payload[1:3])
	cp.DestinationHostname = string(p.Payload[3:])
	return &cp, nil
}

func (p *Packet) DataPacket() (*DataPacket, error) {
	if p.Type != PacketTypeData {
		return nil, fmt.Errorf("packet is not a data packet")
	}
	return &DataPacket{Data: p.Payload}, nil
}

func (p *Packet) ContinuePacket() (*ContinuePacket, error) {
	if p.Type != PacketTypeContinue {
		return nil, fmt.Errorf("packet is not a continue packet")
	}
	return &ContinuePacket{BufferRemaining: binary.LittleEndian.Uint32(p.Payload)}, nil
}

func (p *Packet) ClosePacket() (*ClosePacket, error) {
	if p.Type != PacketTypeClose {
		return nil, fmt.Errorf("packet is not a close packet")
	}
	return &ClosePacket{Reason: CloseReason(p.Payload[0])}, nil
}

func (p *Packet) Marshal() []byte {
	b := make([]byte, 5+len(p.Payload))
	b[0] = byte(p.Type)
	binary.LittleEndian.PutUint32(b[1:5], p.StreamID)
	copy(b[5:], p.Payload)
	return b
}

type StreamType uint8

const (
	StreamTypeTCP StreamType = iota + 1
	StreamTypeUDP
)

type ConnectPacket struct {
	// Whether the connection is TCP or UDP
	StreamType StreamType
	// The destination port: u16
	DestinationPort uint16
	// The destination hostname: string
	DestinationHostname string
}

func BuildConnectPacket(streamType StreamType, destinationPort uint16, destinationHostname string) *ConnectPacket {
	return &ConnectPacket{
		StreamType:          streamType,
		DestinationPort:     destinationPort,
		DestinationHostname: destinationHostname,
	}
}

func (cp *ConnectPacket) Marshal() []byte {
	b := make([]byte, 3+len(cp.DestinationHostname))
	b[0] = byte(cp.StreamType)
	binary.LittleEndian.PutUint16(b[1:3], cp.DestinationPort)
	copy(b[3:], cp.DestinationHostname)
	return b
}

func (cp *ConnectPacket) ToPacket(streamID uint32) *Packet {
	return &Packet{
		Type:     PacketTypeConnect,
		StreamID: streamID,
		Payload:  cp.Marshal(),
	}
}

type DataPacket struct {
	// Data sent to destination server
	Data []byte
}

func BuildDataPacket(data []byte) *DataPacket {
	return &DataPacket{Data: data}
}

func (dp *DataPacket) Marshal() []byte {
	return dp.Data
}

func (dp *DataPacket) ToPacket(streamID uint32) *Packet {
	return &Packet{
		Type:     PacketTypeData,
		StreamID: streamID,
		Payload:  dp.Marshal(),
	}
}

type ContinuePacket struct {
	BufferRemaining uint32
}

func BuildContinuePacket(bufferRemaining uint32) *ContinuePacket {
	return &ContinuePacket{BufferRemaining: bufferRemaining}
}

func (cp *ContinuePacket) Marshal() []byte {
	b := make([]byte, 4)
	binary.LittleEndian.PutUint32(b, cp.BufferRemaining)
	return b
}

func (cp *ContinuePacket) ToPacket(streamID uint32) *Packet {
	return &Packet{
		Type:     PacketTypeContinue,
		StreamID: streamID,
		Payload:  cp.Marshal(),
	}
}

type CloseReason uint8

// Client/Server Close Reasons
const (

	// The connection was closed for an unknown reason
	CloseReasonUnknown CloseReason = iota + 1
	// The connection was closed by one side voluntarily
	CloseReasonVoluntary
	// The connection was closed by one side due to a network error
	CloseReasonNetworkError
)

// Server Close Reasons
const (
	// Stream creation failed due to invalid information
	CloseReasonInvalidStreamInfo CloseReason = iota + 0x41
	// Stream creation failed due to an unreachable destination
	CloseReasonUnreachableDestination
	// Stream creation timed out
	CloseReasonStreamCreationTimeout
	// The server refused the connection
	CloseReasonConnectionRefused
)

// Server Close Reasons caused by the proxy
const (
	// TCP data transfer timed out
	CloseReasonDataTransferTimeout CloseReason = iota + 0x47
	// The proxy blocked the connection
	CloseReasonBlocked
	// The proxy rate-limited the connection
	CloseReasonRateLimited
)

// Client Close Reasons
const (
	// The client has encountered an unexpected error and is unable to recieve any more data
	CloseReasonClientError CloseReason = iota + 0x81
)

type ClosePacket struct {
	// The reason for closing the connection
	Reason CloseReason
}

func BuildClosePacket(reason CloseReason) *ClosePacket {
	return &ClosePacket{Reason: reason}
}

func (cp *ClosePacket) Marshal() []byte {
	b := make([]byte, 1)
	b[0] = byte(cp.Reason)
	return b
}

func (cp *ClosePacket) ToPacket(streamID uint32) *Packet {
	return &Packet{
		Type:     PacketTypeClose,
		StreamID: streamID,
		Payload:  cp.Marshal(),
	}
}
