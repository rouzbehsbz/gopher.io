package engineio

type Middleware func()

type Transport string

const (
	TransportPolling      = "polling"
	TransportWebsocket    = "websocket"
	TransportWebtransport = "webtransport"
)

type ServerOptions struct {
	PingTimeout       int
	PingInterval      int
	UpgradeTimeout    int
	MaxHttpBufferSize int
	AllowRequest      func()
	Transports        []Transport
	AllowUpgrades     bool
}

type Server struct {
	ClientsCount int

	clients     map[string]*Socket
	middlewares []Middleware
}

func NewServer(options ServerOptions) *Server {
	return &Server{
		clients:      make(map[string]*Socket),
		ClientsCount: 0,
	}
}
