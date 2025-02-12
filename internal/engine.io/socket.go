package engineio

import (
	"net/http"

	"github.com/google/uuid"
)

type Socket struct {
	Sid string

	transport Transporter
}

func NewSocket(transport Transporter) *Socket {
	sid := uuid.New().String()

	return &Socket{
		Sid:       sid,
		transport: transport,
	}
}

func (s *Socket) Handle(w http.ResponseWriter, r *http.Request) {
	s.transport.Handle(w, r)
}
