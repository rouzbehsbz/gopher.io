package engineio

import (
	"net/http"

	"github.com/google/uuid"
)

type Socket struct {
	Sid       string
	Transport Transporter

	w http.ResponseWriter
	r *http.Request
}

func NewSocket(w http.ResponseWriter, r *http.Request, transport Transporter) *Socket {
	sid := uuid.New().String()

	return &Socket{
		Sid:       sid,
		Transport: transport,
		w:         w,
		r:         r,
	}
}

func (s *Socket) Handle() {
	s.Transport.Handle(s.w, s.r)
}

func (s *Socket) Send(packet Packet) {
	s.Transport.Send(s.w, s.r, packet)
}
