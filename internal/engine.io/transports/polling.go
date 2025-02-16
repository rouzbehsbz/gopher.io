package transports

import (
	"io"
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
		bodyBytes, err := io.ReadAll(s.R.Body)

		if err != nil {
			http.Error(s.W, "can't read the request body.", http.StatusInternalServerError)
			return
		}

		packets, err := engineio.DecodePackets(bodyBytes)

		if err != nil {
			http.Error(s.W, "can't decode the packet.", http.StatusInternalServerError)
			return
		}

		for _, packet := range packets {
			s.ReceivingPackets <- packet
		}

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
