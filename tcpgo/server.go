package tcpgo

import (
	"fmt"
	"net"
	"sync/atomic"

	"github.com/brianliucrypto/tcpgo/constant"
	"github.com/brianliucrypto/tcpgo/iface"
)

type Server struct {
	Name      string
	Version   string
	IpVersion string
	Ip        string
	Port      uint32

	maxConnection int

	msgHandler        iface.IMessageHandler
	connectionManager iface.IConnectionManager

	onConnectionStartCallback func(iface.IConnection)
	OnConnectionStopCallback  func(iface.IConnection)
}

func NewServer(name, ipVersion, ip string, port uint32) *Server {
	return &Server{
		Name:              name,
		Version:           constant.Version,
		IpVersion:         ipVersion,
		Ip:                ip,
		Port:              port,
		maxConnection:     3,
		msgHandler:        NewMessageHandler(),
		connectionManager: NewConnectionManager(),
	}
}

func (s *Server) Start() {
	listener, err := net.Listen(s.IpVersion, fmt.Sprintf("%v:%v", s.Ip, s.Port))
	if err != nil {
		return
	}

	fmt.Printf("server is running, ip:%v, port:%v\n", s.Ip, s.Port)

	s.msgHandler.Start()

	var connID atomic.Uint32
	for {
		conn, err := listener.Accept()
		if err != nil {
			return
		}

		if s.connectionManager.GetConnectionCount() >= s.maxConnection {
			fmt.Println("too many connections")
			conn.Close()
			continue
		}

		fmt.Printf("new connection, reomte:%v,id:%v\n", conn.RemoteAddr(), connID.Load())
		newConn := NewConneciton(s, connID.Load(), conn)
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
