package transports

import (
	"net/http"

	engineio "github.com/rouzbehsbz/gopher.io/internal/engine.io"
)

type PollingTransport struct {
}

func NewPollingTransport() *PollingTransport {
	return &PollingTransport{}
}

func (p *PollingTransport) Name() string {
	return engineio.TransportPollingType
}

func (p *PollingTransport) Handle(w http.ResponseWriter, r *http.Request) {

}

func (p *PollingTransport) Send(w http.ResponseWriter, r *http.Request, packet engineio.Packet) {
	encodedPacket, err := packet.Encode()

	if err != nil {
		http.Error(w, "can't parse the packet.", 400)
		return
	}

	_, err = w.Write(encodedPacket)

	if err != nil {
		http.Error(w, "Failed to write response", http.StatusInternalServerError)
		return
	}
}
