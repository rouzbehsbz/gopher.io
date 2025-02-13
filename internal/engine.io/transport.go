package engineio

import "net/http"

type TransportType string

const (
	TransportPollingType   = "polling"
	TransportWebsocketType = "websocket"
)

type Transporter interface {
	Name() string
	Handle(w http.ResponseWriter, r *http.Request)
	Send(w http.ResponseWriter, r *http.Request, packet Packet)
}
