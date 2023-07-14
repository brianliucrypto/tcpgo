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

	for {
		conn, err := listener.Accept()
		if err != nil {
			return err
		}

		for {
			readBuf := make([]byte, 1024)
			cnt, err := conn.Read(readBuf)
			if err != nil {
				fmt.Println(err)
				break
			}

			cnt, err = conn.Write(readBuf[:cnt])
			if err != nil {
				fmt.Println(err)
				break
			}

			fmt.Println(string(readBuf[:cnt]))
		}
	}
}

func (s *Server) Stop() {
}

func (s *Server) Serve() {
	go s.Start()

	select {}
}
