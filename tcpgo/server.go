package tcpgo

import (
	"fmt"
	"net"

	"github.com/brianliucrypto/tcpgo/constant"
)

type Server struct {
	Name      string
	Version   string
	IpVersion string
	Ip        string
	Port      uint32
}

func NewServer(name, ipVersion, ip string, port uint32) *Server {
	return &Server{
		Name:      name,
		Version:   constant.Version,
		IpVersion: ipVersion,
		Ip:        ip,
		Port:      port,
	}
}

func (s *Server) Start() error {
	listener, err := net.Listen(s.IpVersion, fmt.Sprintf("%v:%v", s.Ip, s.Port))
	if err != nil {
		return err
	}

	fmt.Printf("server is running, ip:%v, port:%v\n", s.Ip, s.Port)

	connID := 0
	for {
		conn, err := listener.Accept()
		if err != nil {
			return err
		}

		fmt.Printf("new connection, reomte:%v,id:%v\n", conn.RemoteAddr(), connID)
		newConn := NewConneciton(conn, uint32(connID))
		newConn.Start()
	}
}

func (s *Server) Stop() {
}

func (s *Server) Serve() {
	go s.Start()

	select {}
}
