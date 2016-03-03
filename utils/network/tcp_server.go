// A TCP server wrapper, implements Handler interface to handle client connections
package network

import (
	"net"
	"strconv"

	"github.com/123hurray/netroxy/utils/logger"
)

type Server struct {
	ip     string
	port   int
	name   string
	socket *net.TCPListener
}

type Handler interface {
	Handle(net.Conn)
}

// return a new TCP server
func NewTcpServer(name string, ip string, port int) (*Server, error) {
	tcpAddr, err := net.ResolveTCPAddr("tcp4", ip+":"+strconv.Itoa(port))
	l, err := net.ListenTCP("tcp", tcpAddr)
	if err != nil {
		return nil, err
	}
	server := Server{ip, port, name, l}
	return &server, err
}

// Start TCP server with a handler
func (self *Server) Serve(handler Handler) {
	logger.Info(self.name, "Listening", self.ip, ":", self.port)
	for {
		con, err := self.socket.Accept()
		if err != nil {
			return
		}
		go handler.Handle(con)
	}
}

// Close TCP server
func (self *Server) Close() {
	self.socket.Close()
}
