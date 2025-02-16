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
		break

	case http.MethodGet:
		packets := s.Packets()
		encodedPakcets, err := engineio.EncodePackets(packets)

		if err != nil {
			http.Error(s.W, "can't parse the packet.", 400)
			return
		}

		s.W.Header().Set("Content-Type", "text/plain; charset=UTF-8")

		s.W.Write(encodedPakcets)
		break
	}
}
