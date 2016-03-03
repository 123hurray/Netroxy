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
	"bufio"
	"net"
	"strconv"

	"github.com/123hurray/netroxy/utils/logger"

	"github.com/123hurray/netroxy/utils/network"
)

const defaultBufferSize = 16 * 1024

type Handler struct {
	portHandlerDict   map[int]*ProxyHandler
	portTcpServerDict map[int]*network.Server
	username          string
	password          string
}

func NewHandler(username string, password string) *Handler {
	handler := new(Handler)
	handler.portHandlerDict = make(map[int]*ProxyHandler)
	handler.portTcpServerDict = make(map[int]*network.Server)
	handler.username = username
	handler.password = password
	return handler
}

func (self *Handler) Handle(conn net.Conn) {
	reader := bufio.NewReaderSize(conn, defaultBufferSize)
	freeFlag := true
	defer func() {
		if !freeFlag {
			return
		}
		conn.Close()
		for _, s := range self.portTcpServerDict {
			s.Close()
		}
		logger.Info("Ports closed.")
	}()
	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			logger.Warn("Connnection closed")
			return
		}
		line = line[:len(line)-1]
		logger.Debug(line)
		switch {
		case line == "ATH":
			username, err := reader.ReadString('\n')
			username = username[:len(username)-1]
			if err != nil {
				logger.Warn("Parameters error, receive ", line, ". Error:", err)
				return
			}
			password, err := reader.ReadString('\n')
			password = password[:len(password)-1]
			if err != nil {
				logger.Warn("Parameters error, receive ", line, ". Error:", err)
				return
			}
			if username == self.username && password == self.password {
				conn.Write([]byte("ARS\ntrue\n"))
				logger.Debug("Auth OK.")
			} else {
				conn.Write([]byte("ARS\nfalse\n"))
				logger.Warn("Auth failed. Username or password error.")
				return
			}
		case line == "MAP":
			line, err = reader.ReadString('\n')
			if err != nil {
				logger.Warn("Parameters error, receive ", line, ". Error:", err)
				return
			}
			line = line[:len(line)-1]
			port, err := strconv.Atoi(line)
			if err != nil {
				logger.Warn("Parameters error, receive ", line, ". Error:", err)
				return
			}
			s, err := network.NewTcpServer("Netroxy_"+strconv.Itoa(port), "0.0.0.0", port)
			if err != nil {
				logger.Warn("Cannot Listen", port, ". Error:", err)
				// TODO tell client
				return
			}
			handlerProxy := NewProxyHandler(conn, line)
			self.portHandlerDict[port] = handlerProxy
			self.portTcpServerDict[port] = s
			go s.Serve(handlerProxy)
			logger.Info("New connection " + line + " prepared.")
			conn.Write([]byte("MRS\n" + line + "\n"))

		case line == "TRS":
			line, err = reader.ReadString('\n')
			if err != nil {
				logger.Warn("Illegal argument.", err)
				return
			}
			line = line[:len(line)-1]
			port, err := strconv.Atoi(line)
			if err != nil {
				logger.Warn("Illegal argument.", err)
				return
			}
			proxy := self.portHandlerDict[port]
			if proxy == nil {
				logger.Warn("Port", line, "not found.")
				return
			}
			proxy.connChan <- conn
			freeFlag = false
			return
		}

	}
}
