package engineio

import (
	"net/http"

	"github.com/google/uuid"
)

type NewSessionMessage struct {
	Sid          string   `json:"sid"`
	Upgrades     []string `json:"upgrades"`
	PingInterval int      `json:"pingInterval"`
	PingTimeout  int      `json:"pingTimeout"`
	MaxPayload   int      `json:"maxPayload"`
}

type Socket struct {
	Sid       string
	Transport Transporter
}

func NewSocket(transport Transporter) *Socket {
	sid := uuid.New().String()

	return &Socket{
		Sid:       sid,
		Transport: transport,
	}
}

func (s *Socket) Handle(w http.ResponseWriter, r *http.Request) {
	s.Transport.Handle(w, r)
}
