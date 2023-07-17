package tcpgo

import (
	"fmt"
	"net"

	"github.com/brianliucrypto/tcpgo/constant"
	"github.com/brianliucrypto/tcpgo/iface"
)

type Server struct {
	Name      string
	Version   string
	IpVersion string
	Ip        string
	Port      uint32

	msgHandler iface.IMessageHandler
}

func NewServer(name, ipVersion, ip string, port uint32) *Server {
	return &Server{
		Name:       name,
		Version:    constant.Version,
		IpVersion:  ipVersion,
		Ip:         ip,
		Port:       port,
		msgHandler: NewMessageHandler(),
	}
}

func (s *Server) Start() error {
	listener, err := net.Listen(s.IpVersion, fmt.Sprintf("%v:%v", s.Ip, s.Port))
	if err != nil {
		return err
	}

	fmt.Printf("server is running, ip:%v, port:%v\n", s.Ip, s.Port)

	s.msgHandler.Start()

	connID := 0
	for {
		conn, err := listener.Accept()
		if err != nil {
			return err
		}

		fmt.Printf("new connection, reomte:%v,id:%v\n", conn.RemoteAddr(), connID)
		newConn := NewConneciton(uint32(connID), conn, s.msgHandler)
		newConn.Start()
		connID++
	}
}

func (s *Server) Stop() {
	s.msgHandler.Stop()
}

func (s *Server) Serve() {
	go s.Start()

	select {}
}

func (s *Server) AddRouter(msgID uint32, router iface.IRouter) {
	s.msgHandler.AddRouter(msgID, router)
}
