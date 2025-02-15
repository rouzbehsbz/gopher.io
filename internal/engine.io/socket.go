package engineio

import (
	"crypto/rand"
	"encoding/base64"
	"net/http"
	"strings"
)

type Socket struct {
	Sid       string
	Transport Transporter

	w http.ResponseWriter
	r *http.Request
}

func NewSocket(w http.ResponseWriter, r *http.Request, transport Transporter) (*Socket, error) {
	s := &Socket{
		Transport: transport,
		w:         w,
		r:         r,
	}

	sid, err := s.generateSid()

	if err != nil {
		return nil, err
	}

	s.Sid = sid

	return s, nil
}

func (s *Socket) generateSid() (string, error) {
	bytes := make([]byte, 15)

	_, err := rand.Read(bytes)

	if err != err {
		return "", nil
	}

	sid := strings.TrimRight(base64.URLEncoding.EncodeToString(bytes), "=")

	return sid, nil
}

func (s *Socket) Handle() {
	s.Transport.Handle(s.w, s.r)
}

func (s *Socket) Send(packet Packet) {
	s.Transport.Send(s.w, s.r, packet)
}
