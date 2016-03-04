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
	"errors"
	"net"
	"strconv"
	"time"

	"github.com/123hurray/netroxy/common"

	"github.com/123hurray/netroxy/utils/logger"

	"github.com/123hurray/netroxy/utils/network"
)

const defaultBufferSize = 16 * 1024

type Handler struct {
	portHandlerDict   map[int]*ProxyHandler
	portTcpServerDict map[int]*network.Server
	config            *ServerConfig
	password          string
	clients           map[net.Conn]*ClientConn
	reader            *bufio.Reader
}

func NewHandler(config *ServerConfig) *Handler {
	handler := new(Handler)
	handler.portHandlerDict = make(map[int]*ProxyHandler)
	handler.portTcpServerDict = make(map[int]*network.Server)
	handler.clients = make(map[net.Conn]*ClientConn)
	handler.config = config
	return handler
}

func (self *Handler) Supervise() {
	now := time.Now()
	for conn, cli := range self.clients {
		if now.Sub(cli.expireTime) > time.Duration(0) {
			conn.Close()
		}
	}
}

func (self *Handler) getString() (str string, err error) {
	str, err = self.reader.ReadString('\n')
	str = str[:len(str)-1]
	return
}
func (self *Handler) getInt() (i int, err error) {
	str, err := self.getString()
	if err != nil {
		return
	}
	i, err = strconv.Atoi(str)
	if err != nil {
		return
	}
	return
}
func (self *Handler) getBool() (b bool, err error) {
	str, err := self.getString()
	if err != nil {
		return
	}
	if str == "true" {
		b = true
	} else if str == "false" {
		b = false
	} else {
		err = errors.New("Illegal bool value. Receive:" + str)
	}
	return
}

func (self *Handler) Handle(conn net.Conn) {
	reader := bufio.NewReaderSize(conn, defaultBufferSize)
	self.reader = reader
	freeFlag := true
	defer func() {
		if !freeFlag {
			return
		}
		delete(self.clients, conn)

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
			name, err := self.getString()
			if err != nil {
				logger.Warn("Parameters error, receive ", name, ". Error:", err)
				return
			}
			username, err := self.getString()
			if err != nil {
				logger.Warn("Parameters error, receive ", username, ". Error:", err)
				return
			}
			password, err := self.getString()
			if err != nil {
				logger.Warn("Parameters error, receive ", password, ". Error:", err)
				return
			}
			if username == self.config.Username && password == self.config.Password {
				self.clients[conn] = NewClientConn(conn, name, self.config.Timeout)
				conn.Write([]byte("ARS\ntrue\n" + strconv.Itoa(self.config.Timeout) + "\n"))
				logger.Debug("Auth OK.")
			} else {
				conn.Write([]byte("ARS\nfalse\n"))
				logger.Warn("Auth failed. Username or password error.")
				return
			}
		case line == "SRQ":
			self.clients[conn].UpdateExpireTime()
			conn.Write([]byte("SRS\n"))
		case line == "MAP":
			port, err := self.getInt()
			if err != nil {
				logger.Warn("Parameters error:", err)
				return
			}
			mapAddress, err := self.getString()
			if err != nil {
				logger.Warn("Parameters error:", err)
				return
			}
			isOpen, err := self.getBool()
			if err != nil {
				logger.Warn("Parameters error:", err)
				return
			}
			s, err := network.NewTcpServer("Netroxy_"+strconv.Itoa(port), "0.0.0.0", port)
			if err != nil {
				logger.Warn("Cannot Listen", port, ". Error:", err)
				// TODO tell client
				return
			}
			cliHost, cliPortStr, _ := net.SplitHostPort(mapAddress)
			cliPort, _ := strconv.Atoi(cliPortStr)
			mapping := common.NewMapping(cliHost, cliPort, port, isOpen)
			handlerProxy := NewProxyHandler(conn, mapping)
			self.portHandlerDict[port] = handlerProxy
			self.portTcpServerDict[port] = s
			go s.Serve(handlerProxy)
			logger.Info("New connection " + strconv.Itoa(port) + " prepared.")
			conn.Write([]byte("MRS\n" + strconv.Itoa(port) + "\n"))

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
