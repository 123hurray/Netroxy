/*
 * The MIT License (MIT)
 *
 * Copyright (c) 2016 Ray Zhang
 *
 * Permission is hereby granted, free of charge, to any person obtaining a copy
 * of this software and associated documentation files (the "Software"), to deal
 * in the Software without restriction, including without limitation the rights
 * to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
 * copies of the Software, and to permit persons to whom the Software is
 * furnished to do so, subject to the following conditions:
 *
 * The above copyright notice and this permission notice shall be included in all
 * copies or substantial portions of the Software.
 *
 * THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
 * IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
 * FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
 * AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
 * LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
 * OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
 * SOFTWARE.
 */

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
