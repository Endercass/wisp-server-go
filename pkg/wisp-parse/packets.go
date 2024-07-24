package wispparse

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
	Type PacketType `json:"type"`
	// The stream ID
	StreamID uint32 `json:"stream_id"`
	// The payload of the packet
	Payload []byte `json:"payload"`
}

type StreamType uint8

const (
	StreamTypeTCP StreamType = iota + 1
	StreamTypeUDP
)

type ConnectPacket struct {
	// Whether the connection is TCP or UDP
	StreamType StreamType `json:"stream_type"`
	// The destination port: u16
	DestinationPort uint16 `json:"destination_port"`
	// The destination hostname: string
	DestinationHostname string `json:"destination_hostname"`
}

type DataPacket struct {
	// Data sent to destination server
	Data []byte `json:"data"`
}

type ContinuePacket struct {
	BufferRemaining uint32 `json:"buffer_remaining"`
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
	Reason CloseReason `json:"reason"`
}
