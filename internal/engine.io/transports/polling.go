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
