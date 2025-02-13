package main

import (
	engineio "github.com/rouzbehsbz/gopher.io/internal/engine.io"
	"github.com/rouzbehsbz/gopher.io/internal/engine.io/transports"
)

func main() {
	polling := transports.NewPollingTransport()
	eio := engineio.NewServer(engineio.ServerOpt{
		Transports: []engineio.Transporter{polling},
	})

	eio.Listen("0.0.0.0:3000")
}
