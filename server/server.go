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
	"time"

	"github.com/123hurray/netroxy/common"

	"github.com/123hurray/netroxy/utils/logger"
	"github.com/123hurray/netroxy/utils/security"

	"github.com/123hurray/netroxy/utils/network"
)

const defaultBufferSize = 16 * 1024

type Server struct {
	config  *ServerConfig
	clients map[string]*ClientConn
}

func NewServer(config *ServerConfig) *Server {
	handler := new(Server)
	handler.clients = make(map[string]*ClientConn)
	handler.config = config
	return handler
}

func (self *Server) Supervise() {
	now := time.Now()
	for _, cli := range self.clients {
		if now.Sub(cli.expireTime) > time.Duration(0) {
			cli.conn.Close()
		}
	}
}

type ClientReader struct {
	common.ProtocolReader
}

func (self *Server) Handle(conn net.Conn) {
	clientReader := ClientReader{}
	clientReader.SetReader(bufio.NewReaderSize(conn, defaultBufferSize))
	freeFlag := true
	var client *ClientConn
	var token string
	defer func() {
		if freeFlag == false {
			return
		}
		conn.Close()
		if token == "" {
			return
		}
		delete(self.clients, token)
		if client.handlers == nil {
			return
		}
		for i, _ := range client.handlers {
			client.GetHandler(i).Free()
		}
		client.handlers = nil
		logger.Info("Ports closed.")
	}()
	for {
		line, err := clientReader.GetString()
		if err != nil {
			logger.Warn("Connnection closed.", err)
			return
		}
		logger.Debug(line)
		switch {
		case line == "ATH":
			name, err := clientReader.GetString()
			if err != nil {
				logger.Warn("Parameters error, receive ", name, ". Error:", err)
				return
			}
			username, err := clientReader.GetString()
			if err != nil {
				logger.Warn("Parameters error, receive ", username, ". Error:", err)
				return
			}
			password, err := clientReader.GetString()
			if err != nil {
				logger.Warn("Parameters error, receive ", password, ". Error:", err)
				return
			}
			if username == self.config.Username && password == self.config.Password {
				// Auth passed
				token = security.GenerateUID(16)
				client = NewClientConn(conn, name, token, self.config.Timeout)
				self.clients[token] = client
				conn.Write([]byte("ARS\ntrue\n" + strconv.Itoa(self.config.Timeout) + "\n" + token + "\n"))
				logger.Debug("Client", name, "Auth OK.")
			} else {
				conn.Write([]byte("ARS\nfalse\n"))
				logger.Warn("Auth failed. Username or password error.")
				return
			}
		case line == "SRQ":
			if token == "" {
				logger.Warn("Token not found.")
				return
			}
			client.UpdateExpireTime()
			conn.Write([]byte("SRS\n"))
		case line == "MAP":
			if token == "" {
				logger.Warn("Token not found.")
				return
			}
			port, err := clientReader.GetInt()
			if err != nil {
				logger.Warn("Parameters error:", err)
				return
			}
			mapAddress, err := clientReader.GetString()
			if err != nil {
				logger.Warn("Parameters error:", err)
				return
			}
			isOpen, err := clientReader.GetBool()
			if err != nil {
				logger.Warn("Parameters error:", err)
				return
			}
			s, err := network.NewPlainServer("Netroxy_"+strconv.Itoa(port), "0.0.0.0", port)
			if err != nil {
				logger.Warn("Cannot Listen", port, ". Error:", err)
				conn.Write([]byte("MRS\n" + strconv.Itoa(port) + "\nfalse\n"))
				break
			}
			cliHost, cliPortStr, _ := net.SplitHostPort(mapAddress)
			cliPort, _ := strconv.Atoi(cliPortStr)
			mapping := common.NewMapping(cliHost, cliPort, port, isOpen)
			handlerProxy := NewProxyHandler(conn, s, mapping)
			client.AddHandler(handlerProxy)
			go s.Serve(handlerProxy)
			logger.Info("New connection " + strconv.Itoa(port) + " prepared.")
			conn.Write([]byte("MRS\n" + strconv.Itoa(port) + "\ntrue\n"))

		case line == "TRS":
			newToken, err := clientReader.GetString()
			if err != nil {
				logger.Warn("Illegal argument.", err)
				return
			}
			port, err := clientReader.GetInt()
			if err != nil {
				logger.Warn("Illegal argument.", err)
				return
			}
			client := self.clients[newToken]
			if client == nil {
				logger.Warn("Client not found.")
				return
			}
			proxy := client.handlers[port]
			if proxy == nil {
				logger.Warn("Port", port, "not found.")
				return
			}
			freeFlag = false
			proxy.connChan <- conn
			logger.Debug("Connection has been sent to proxy")
			return
		}

	}
}
