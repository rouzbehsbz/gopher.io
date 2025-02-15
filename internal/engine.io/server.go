package engineio

import (
	"encoding/json"
	"errors"
	"net/http"
	"sync"
)

type NewSessionMessage struct {
	Sid          string   `json:"sid"`
	Upgrades     []string `json:"upgrades"`
	PingInterval int      `json:"pingInterval"`
	PingTimeout  int      `json:"pingTimeout"`
	MaxPayload   int      `json:"maxPayload"`
}

type ServerOpt struct {
	Transports []Transporter
}

type Server struct {
	sockets      map[string]*Socket
	socketsCount int

	transports     map[string]Transporter
	transportNames []string

	pingInterval int
	pintTimeout  int
	maxPayload   int

	mu sync.Mutex
}

func NewServer(opt ServerOpt) *Server {
	transports := make(map[string]Transporter)

	var transportNames []string

	for _, transport := range opt.Transports {
		name := transport.Name()

		transports[name] = transport
		transportNames = append(transportNames, name)
	}

	return &Server{
		sockets:        make(map[string]*Socket),
		transports:     transports,
		transportNames: transportNames,
		pingInterval:   25000,
		pintTimeout:    20000,
		maxPayload:     1000000,
	}
}

func (s *Server) Listen(address string) {
	http.HandleFunc("/engine.io/", s.HandleHandshake)
	http.ListenAndServe(address, nil)
}

func (s *Server) HandleHandshake(w http.ResponseWriter, r *http.Request) {
	eio := r.URL.Query().Get("EIO")
	sid := r.URL.Query().Get("sid")
	requestTransport := r.URL.Query().Get("transport")

	if isProtocolVersionSupported := s.isProtocolVersionSupported(eio); isProtocolVersionSupported == false {
		s.ErrorResponse(w, UnsupportedProtocolVersionErrorCode)
		return
	}

	serverTransport, isTransportExists := s.getTransport(requestTransport)

	if !isTransportExists {
		s.ErrorResponse(w, UnknownTransportErrorCode)
		return
	}

	if sid == "" {
		newSid, err := s.AddSocket(w, r, serverTransport)

		if err != nil {
			s.ErrorResponse(w, BadRequestErrorCode)
			return
		}

		sid = newSid
	}

	socket, err := s.GetSocket(sid)

	if err != nil {
		s.ErrorResponse(w, UnknownSidErrorCode)
		return
	}

	socket.Handle()
}

func (s *Server) isProtocolVersionSupported(eio string) bool {
	return eio == "4"
}

func (s *Server) getTransport(requestTransport string) (Transporter, bool) {
	serverTransport, isExists := s.transports[requestTransport]

	return serverTransport, isExists
}

func (s *Server) AddSocket(w http.ResponseWriter, r *http.Request, transport Transporter) (string, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	socket, err := NewSocket(w, r, transport)

	if err != nil {
		return "", nil
	}

	s.sockets[socket.Sid] = socket

	socket.Send(Packet{
		Type: PacketOpenType,
		RawData: NewSessionMessage{
			Sid:          socket.Sid,
			Upgrades:     s.transportNames,
			PingInterval: s.pingInterval,
			PingTimeout:  s.pintTimeout,
			MaxPayload:   s.maxPayload,
		},
	})

	return socket.Sid, nil
}

func (s *Server) GetSocket(sid string) (*Socket, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	socket, isExists := s.sockets[sid]

	if !isExists {
		return nil, errors.New("session not found.")
	}

	return socket, nil
}

func (s *Server) ErrorResponse(w http.ResponseWriter, errorCode ErrorCode) {
	message := GetErrorMessage(errorCode)
	httpStatusCode := GetErrorHttpStatusCode(errorCode)

	appError := AppError{
		Code:    errorCode,
		Message: message,
	}

	bytes, _ := json.Marshal(appError)

	http.Error(w, string(bytes), httpStatusCode)
}
