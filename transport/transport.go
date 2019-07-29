package transport

// Message message definition
type Message struct {
	Header map[string]string
	Body   []byte
}

// Socket socket definition
type Socket interface {
	Recv(*Message) error
	Send(*Message) error
	Close() error
	Local() string
	Remote() string
}

// Client client definition
type Client interface {
	Socket
}

// Listener listener definition
type Listener interface {
	Addr() string
	Close() error
	Accept(func(Socket)) error
}

// Transport is an interface which is used for communication between
// services. It uses socket send/recv semantics and had various
// implementations {HTTP, RabbitMQ, NATS, ...}
type Transport interface {
	Init(...Option) error
	Options() Options
	Dial(addr string, opts ...DialOption) (Client, error)
	Listen(addr string, opts ...ListenOption) (Listener, error)
	String() string
}
