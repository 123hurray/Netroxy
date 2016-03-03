package server

import (
	"io"
	"net"

	"github.com/123hurray/netroxy/utils/logger"
)

type ProxyHandler struct {
	mainConn net.Conn
	connChan chan net.Conn
	portStr  string
}

func NewProxyHandler(mainConn net.Conn, portStr string) *ProxyHandler {
	self := new(ProxyHandler)
	self.connChan = make(chan net.Conn)
	self.mainConn = mainConn
	self.portStr = portStr
	return self
}

func (self *ProxyHandler) Handle(conn net.Conn) {

	logger.Info("New user request", conn.LocalAddr(), "from", conn.RemoteAddr())
	self.mainConn.Write([]byte("Dial\n" + self.portStr + "\n"))
	conn1 := <-self.connChan
	defer conn1.Close()
	defer conn.Close()
	logger.Info("Forwarding tcp connection...")
	exitChan := make(chan bool)
	go func() {
		io.Copy(conn, conn1)
		_, isClose := <-exitChan
		if !isClose {
			close(exitChan)
		}
	}()
	go func() {
		io.Copy(conn1, conn)
		_, isClose := <-exitChan
		if !isClose {
			close(exitChan)
		}
	}()
	<-exitChan
	logger.Info("Proxy connection closed.")

}

func Serve() {

}
