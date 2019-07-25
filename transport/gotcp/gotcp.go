package gotcp

import (
	"context"
	"fmt"
	"net"

	"github.com/fananchong/gotcp"
	"github.com/fananchong/v-micro/transport"
)

type socketImpl struct {
	gotcp.Session
	chanRecvMsg chan *transport.Message
}

func (socket *socketImpl) Init(ctx context.Context, conn net.Conn, derived gotcp.ISession) {
	socket.Session.Init(ctx, conn, derived)
	socket.chanRecvMsg = make(chan *transport.Message, 128)
}

func (socket *socketImpl) Recv(msg *transport.Message) error {
	tmpMsg := <-socket.chanRecvMsg
	if tmpMsg != nil {
		msg.Header = tmpMsg.Header
		msg.Body = tmpMsg.Body
		return nil
	}
	return fmt.Errorf("RECV EOF")
}

func (socket *socketImpl) Send(msg *transport.Message) error {
	data, err := marshal(msg)
	if err != nil {
		return err
	}
	if ok := socket.Session.Send(data, 0); !ok {
		return fmt.Errorf("send msg fail")
	}
	return nil
}

func (socket *socketImpl) Close() (err error) {
	socket.Session.Close()
	if socket.chanRecvMsg != nil {
		select {
		case <-socket.chanRecvMsg:
		default:
			close(socket.chanRecvMsg)
		}
	}
	return
}

func (socket *socketImpl) Local() string {
	return socket.LocalAddr()
}

func (socket *socketImpl) Remote() string {
	return socket.RemoteAddr()
}

func (socket *socketImpl) OnRecv(data []byte, flag byte) {
	msg, err := unmarshal(data)
	if err != nil {
		// TODO: 打印错误
		socket.Close()
		return
	}
	socket.chanRecvMsg <- msg
}

func (socket *socketImpl) OnClose() {
	// TODO: 打印关闭
	if socket.chanRecvMsg != nil {
		select {
		case <-socket.chanRecvMsg:
		default:
			close(socket.chanRecvMsg)
		}
	}
}

type clientImpl struct {
	socketImpl
	opts transport.DialOptions
}

type listenerImpl struct {
	gotcp.Server
	opts transport.ListenOptions
}

func (listener *listenerImpl) Addr() string {
	return listener.GetAddress()
}

func (listener *listenerImpl) Close() (err error) {
	listener.Server.Close()
	return
}

func (listener *listenerImpl) Accept(fn func(transport.Socket)) (err error) {
	listener.RegisterSessType(socketImpl{})
	listener.Server.Accept(func(sock interface{}) {
		s := sock.(*socketImpl)
		s.Verify()
		fn(s)
	})
	return
}

type transportImpl struct {
	opts transport.Options
}

func (trans *transportImpl) Init(opts ...transport.Option) (err error) {
	for _, o := range opts {
		o(&trans.opts)
	}
	if trans.opts.RecvBufSize != 0 {
		gotcp.DefaultRecvBuffSize = trans.opts.RecvBufSize
	}
	if trans.opts.SendBufSize != 0 {
		gotcp.DefaultSendBuffSize = trans.opts.SendBufSize
	}
	return
}

func (trans *transportImpl) Options() transport.Options {
	return trans.opts
}

func (trans *transportImpl) Dial(addr string, opts ...transport.DialOption) (transport.Client, error) {
	cliImpl := &clientImpl{}
	for _, o := range opts {
		o(&cliImpl.opts)
	}
	if ok := cliImpl.Connect(addr, cliImpl); ok {
		go cliImpl.Verify()
		return cliImpl, nil
	}
	return nil, fmt.Errorf("connect fail, addr:%s", addr)
}

func (trans *transportImpl) Listen(addr string, opts ...transport.ListenOption) (transport.Listener, error) {
	listenerImpl := &listenerImpl{}
	for _, o := range opts {
		o(&listenerImpl.opts)
	}
	if err := listenerImpl.Listen(addr); err != nil {
		return nil, err
	}
	return listenerImpl, nil
}

func (trans *transportImpl) String() string {
	return "gotcp"
}

// NewTransport new
func NewTransport(opts ...transport.Option) transport.Transport {
	var options transport.Options
	for _, o := range opts {
		o(&options)
	}
	return &transportImpl{
		opts: options,
	}
}
