package tcpgo

import (
	"fmt"
	"net"
	"sync/atomic"

	"github.com/brianliucrypto/tcpgo/constant"
	"github.com/brianliucrypto/tcpgo/iface"
	"github.com/brianliucrypto/tcpgo/tlog"
)

type Server struct {
	Name      string
	Version   string
	IpVersion string
	Ip        string
	Port      uint32

	maxConnection  int
	workPoolSize   int
	writeCacheSize int

	msgHandler        iface.IMessageHandler
	connectionManager iface.IConnectionManager

	onConnectionStartCallback func(iface.IConnection)
	OnConnectionStopCallback  func(iface.IConnection)
}

func NewServer(optins ...func(*Server)) *Server {
	server := &Server{
		Name:              "tcpgo server",
		Version:           constant.Version,
		IpVersion:         "tcp4",
		Ip:                "",
		Port:              8888,
		maxConnection:     1024,
		workPoolSize:      4,
		writeCacheSize:    1024,
		connectionManager: NewConnectionManager(),
	}

	for _, opt := range optins {
		opt(server)
	}

	server.msgHandler = NewMessageHandler(uint32(server.workPoolSize))

	return server
}

func (s *Server) Start() {
	listener, err := net.Listen(s.IpVersion, fmt.Sprintf("%v:%v", s.Ip, s.Port))
	if err != nil {
		return
	}

	tlog.Info("server is running, ip:%v, port:%v", s.Ip, s.Port)

	s.msgHandler.Start()

	var connID atomic.Uint32
	for {
		conn, err := listener.Accept()
		if err != nil {
			return
		}

		if s.connectionManager.GetConnectionCount() >= s.maxConnection {
			tlog.Info("too many connections")
			conn.Close()
			continue
		}

		tlog.Info("new connection, reomte:%v,id:%v", conn.RemoteAddr(), connID.Load())
		newConn := NewConneciton(s, connID.Load(), uint32(s.writeCacheSize), conn)
		newConn.Start()
		connID.Add(1)
	}
}

func (s *Server) Stop() {
	s.msgHandler.Stop()
	s.connectionManager.ClearConnection()
}

func (s *Server) Serve() {
	go s.Start()

	select {}
}

func (s *Server) AddRouter(msgID uint32, router iface.IRouter) {
	s.msgHandler.AddRouter(msgID, router)
}

func (s *Server) GetMessageHandler() iface.IMessageHandler {
	return s.msgHandler
}

func (s *Server) GetConnectionManager() iface.IConnectionManager {
	return s.connectionManager
}

func (s *Server) SetOnConnStart(callback func(iface.IConnection)) {
	s.onConnectionStartCallback = callback
}

// set on connection callback
func (s *Server) SetOnConnStop(callback func(iface.IConnection)) {
	s.OnConnectionStopCallback = callback
}

// call on connection callback
func (s *Server) CallOnConnStart(conn iface.IConnection) {
	if s.onConnectionStartCallback != nil {
		s.onConnectionStartCallback(conn)
	}
}

// call on connection callback
func (s *Server) CallOnConnStop(conn iface.IConnection) {
	if s.OnConnectionStopCallback != nil {
		s.OnConnectionStopCallback(conn)
	}
}

func WithName(name string) func(*Server) {
	return func(s *Server) {
		s.Name = name
	}
}

func WithIpVersion(version string) func(*Server) {
	return func(s *Server) {
		s.IpVersion = version
	}
}

func WithIp(ip string) func(*Server) {
	return func(s *Server) {
		s.Ip = ip
	}
}

func WithPort(port uint32) func(*Server) {
	return func(s *Server) {
		s.Port = port
	}
}

func WithMaxConnection(maxConnection int) func(*Server) {
	return func(s *Server) {
		s.maxConnection = maxConnection
	}
}

func WithWorkPoolSize(workPoolSize int) func(*Server) {
	return func(s *Server) {
		s.workPoolSize = workPoolSize
	}
}

func WithWriteCacheSize(writeCacheSize int) func(*Server) {
	return func(s *Server) {
		s.writeCacheSize = writeCacheSize
	}
}
