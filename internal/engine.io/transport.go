package engineio

type TransportType string

const (
	TransportPollingType   = "polling"
	TransportWebsocketType = "websocket"
)

type Transporter interface {
}
