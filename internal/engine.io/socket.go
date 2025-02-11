package engineio

import "github.com/google/uuid"

type Socket struct {
	Sid       string
	Transport string
}

func NewSocket(transport string) *Socket {
	sid := uuid.New().String()

	return &Socket{
		Sid:       sid,
		Transport: transport,
	}
}
