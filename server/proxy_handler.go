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

package server

import (
	"io"
	"net"
	"strconv"

	"github.com/123hurray/netroxy/common"
	"github.com/123hurray/netroxy/utils/logger"
)

type ProxyHandler struct {
	mainConn net.Conn
	connChan chan net.Conn
	mapping  *common.Mapping
}

func NewProxyHandler(mainConn net.Conn, mapping *common.Mapping) *ProxyHandler {
	self := new(ProxyHandler)
	self.connChan = make(chan net.Conn)
	self.mainConn = mainConn
	self.mapping = mapping
	return self
}

func (self *ProxyHandler) Handle(conn net.Conn) {
	logger.Info("New user request", conn.LocalAddr(), "from", conn.RemoteAddr())
	if self.mapping.IsOpen == false {
		logger.Info("Reject connection.")
		conn.Close()
		return
	}

	self.mainConn.Write([]byte("TRQ\n" + strconv.Itoa(self.mapping.RemotePort) + "\n"))
	conn1 := <-self.connChan

	logger.Info("Forwarding tcp connection...")
	go func() {
		io.Copy(conn, conn1)
		conn.Close()
		logger.Info("Proxy connection1 closed.")
	}()

	io.Copy(conn1, conn)
	conn1.Close()
	logger.Info("Proxy connection2 closed.")

}

func Serve() {

}
