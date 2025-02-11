package engineio

import "github.com/google/uuid"

type Socket struct {
	Sid string
}

func NewSocket() *Socket {
	sid := uuid.New().String()

	return &Socket{
		Sid: sid,
	}
}
