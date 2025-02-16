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

func (p *PollingTransport) Handle(s *engineio.Socket) {
	switch s.R.Method {
	case http.MethodPost:
		s.W.Header().Set("Content-Type", "text/plain; charset=UTF-8")
		s.W.Write([]byte("ok"))

		break

	case http.MethodGet:
		firstPacket := <-s.SendingPackets
		packets := []engineio.Packet{firstPacket}

	CollectLoop:
		for {
			select {
			case packet := <-s.SendingPackets:
				packets = append(packets, packet)
			default:
				break CollectLoop
			}
		}

		encodedPackets, err := engineio.EncodePackets(packets)
		if err != nil {
			http.Error(s.W, "can't encode the packet.", http.StatusInternalServerError)
			return
		}

		s.W.Header().Set("Content-Type", "text/plain; charset=UTF-8")

		s.W.Write(encodedPackets)

		break
	}
}
