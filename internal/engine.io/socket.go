package engineio

import (
	"crypto/rand"
	"encoding/base64"
	"net/http"
	"strings"
	"sync"
)

type Socket struct {
	Sid       string
	Transport Transporter
	W         http.ResponseWriter
	R         *http.Request

	packets []Packet
	mu      sync.Mutex
}

func NewSocket(w http.ResponseWriter, r *http.Request, transport Transporter) (*Socket, error) {
	s := &Socket{
		Transport: transport,
		W:         w,
		R:         r,
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
	s.Transport.Handle(s)
}

func (s *Socket) Send(packet Packet) {
	s.mu.Lock()

	defer s.mu.Unlock()

	s.packets = append(s.packets, packet)
}

func (s *Socket) Packets() []Packet {
	s.mu.Lock()

	defer s.mu.Unlock()

	packets := s.packets
	s.packets = []Packet{}

	return packets
}
