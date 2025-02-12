package engineio

import (
	"encoding/json"
	"errors"
	"net/http"
	"sync"
)

type Server struct {
	sockets      map[string]*Socket
	socketsCount int

	transports map[string]Transporter

	mu sync.Mutex
}

func NewServer() *Server {
	return &Server{
		sockets: make(map[string]*Socket),
	}
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
		sid = s.AddSocket(requestTransport)
	}

	socket, err := s.GetSocket(sid)

	if err != nil {
		s.ErrorResponse(w, UnknownSidErrorCode)
		return
	}

}

func (s *Server) isProtocolVersionSupported(eio string) bool {
	return eio == "4"
}

func (s *Server) getTransport(requestTransport string) (Transporter, bool) {
	serverTransport, isExists := s.transports[requestTransport]

	return serverTransport, isExists
}

func (s *Server) AddSocket(requestTransport string) string {
	s.mu.Lock()
	defer s.mu.Unlock()

	socket := NewSocket(requestTransport)

	s.sockets[socket.Sid] = socket

	return socket.Sid
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
