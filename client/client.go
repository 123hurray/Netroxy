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

package client

import (
	"bufio"
	"crypto/tls"
	"errors"
	"io"
	"net"
	"os"
	"strconv"
	"sync"
	"sync/atomic"
	"time"

	"github.com/123hurray/netroxy/common"
	"github.com/123hurray/netroxy/utils/logger"
	"github.com/123hurray/netroxy/utils/security"
)

const defaultBufferSize = 16 * 1024

type Client struct {
	common.ProtocolReader
	conn        net.Conn
	targets     map[int]*common.Mapping
	mappingLock sync.RWMutex
	ip          string
	port        int
	config      *ClientConfig
	exitChan    chan bool
	expireTime  int32
	timeout     int
	name        string
	token       string
}

func NewClient(config *ClientConfig) *Client {
	client := new(Client)
	client.ip = config.Ip
	client.port = config.Port
	client.config = config
	client.targets = make(map[int]*common.Mapping)
	client.exitChan = make(chan bool)
	client.expireTime = 0
	name, err := os.Hostname()
	if err != nil {
		logger.Fatal("Cannot get Hostname.")
	}
	client.name = name + "-" + security.GenerateUID(8)
	return client
}

func (self *Client) send(content string) {
	self.conn.Write([]byte(content))
}

func (self *Client) Login() error {
	var conn net.Conn
	var err error
	if self.config.TLS.Enabled == true {
		tlsConfig := tls.Config{InsecureSkipVerify: !self.config.TLS.Verify}
		logger.Debug("Using TLS.")
		conn, err = tls.Dial("tcp", self.ip+":"+strconv.Itoa(self.port), &tlsConfig)
	} else {
		conn, err = net.Dial("tcp", self.ip+":"+strconv.Itoa(self.port))
	}
	if err != nil {
		return err
	}
	self.conn = conn
	logger.Debug("Client name:", self.name)
	self.auth(self.name, self.config.Username, self.config.Password)
	self.SetReader(bufio.NewReaderSize(conn, defaultBufferSize))
	ars, err := self.GetString()
	if err != nil {
		return err
	}
	if ars != "ARS" {
		logger.Warn("Illegal command, expected ARS but receive", ars)
		return errors.New("Illegal command")
	}
	isOK, err := self.GetBool()
	if err != nil {
		return err
	}
	if isOK == true {
		timeout, err := self.GetInt()
		if err != nil {
			logger.Warn("Illegal parameter")
			return err
		}
		token, err := self.GetString()
		if err != nil {
			logger.Warn("Illegal parameter")
			return err
		}
		self.timeout = timeout
		self.token = token
		logger.Info("Login to server success.")
		go self.supervise()
		go self.handle()
		return nil
	} else {
		return errors.New("Auth failed")
	}
}

func (self *Client) supervise() {
	ticker := time.NewTicker(time.Duration(self.timeout) * time.Second)
	for {
		select {
		case <-ticker.C:
			if atomic.AddInt32(&self.expireTime, 1) == 3 {
				logger.Warn("Supervise failed.")
				ticker.Stop()
				self.Close()
				return
			} else {
				self.superviseRequest()
			}
		case <-self.exitChan:
			ticker.Stop()
			logger.Info("Supervise stop.")
			return
		}
	}
}

func (self *Client) Close() {
	err := self.conn.Close()
	if err == nil {
		logger.Debug("Close self.exitChan")
		close(self.exitChan)
	}
}

func (self *Client) handle() {
	defer self.Close()
	for {
		command, err := self.GetString()
		if err != nil {
			logger.Warn("Connection closed.", err)
			return
		}
		switch {
		case command == "SRS":
			atomic.StoreInt32(&self.expireTime, 0)
		case command == "MRS":
			remotePort, err := self.GetInt()
			if err != nil {
				logger.Warn("Illegal parament.", err)
				return
			}
			isOk, err := self.GetBool()
			if err != nil {
				logger.Warn("Illegal parament.", err)
				return
			}
			self.mappingLock.RLock()
			t := self.targets[remotePort]
			self.mappingLock.RUnlock()
			if t != nil {

				if isOk == false {
					logger.Warn("Map", t.Addr(), "port failed.")
					break
				}
				logger.Info("Mapping", self.conn.RemoteAddr(), "<->", t.Addr(), "accepted.")
			}
		case command == "TRQ":
			logger.Info("Tunnel request.")
			addr := self.ip + ":" + strconv.Itoa(self.port)
			remotePort, err := self.GetInt()
			if err != nil {
				logger.Warn("Illegal parament.", err)
				return
			}
			self.mappingLock.RLock()
			t := self.targets[remotePort]
			self.mappingLock.RUnlock()
			if t != nil {
				logger.Info("New tunnel", self.ip+":"+strconv.Itoa(remotePort), "<->", t.Addr(), "Establishing...")
				var conn1 net.Conn
				if self.config.TLS.Enabled == true {
					tlsConfig := tls.Config{InsecureSkipVerify: !self.config.TLS.Verify}
					conn1, err = tls.Dial("tcp", addr, &tlsConfig)
				} else {
					conn1, err = net.Dial("tcp", addr)
				}
				if err != nil {
					logger.Warn("Cannot connect to", addr, err)
					break
				}
				conn2, err := net.Dial("tcp", t.Addr())
				if err != nil {
					logger.Warn("Cannot connect to", t.Addr(), err)
					break
				}
				self.channelResponse(conn1, remotePort, self.token)
				logger.Info("Dial " + t.Addr() + " OK")
				logger.Info("New tunnel", self.ip+":"+strconv.Itoa(remotePort), "<->", t.Addr(), "created.")
				go func() {
					io.Copy(conn1, conn2)
					logger.Debug("Proxy conn1 closed.")
					defer conn1.Close()
				}()
				go func() {
					io.Copy(conn2, conn1)
					logger.Debug("Proxy conn2 closed.")
					defer conn2.Close()
				}()
			}
		default:
			logger.Warn("Illegal command:", command)
			return
		}
	}
}

func (self *Client) Wait() {
	<-self.exitChan
}
func (self *Client) Connect(mapConfig *ConnectionConfig) (*common.Mapping, error) {
	addr := mapConfig.Ip + ":" + strconv.Itoa(mapConfig.Port)
	t := common.NewMapping(mapConfig.Ip, mapConfig.Port, mapConfig.RemotePort, mapConfig.IsOpen)
	logger.Info("Send new mapping", addr, ":", t.Addr(), "request...")
	self.mapRequest(mapConfig.RemotePort, addr, mapConfig.IsOpen)
	self.mappingLock.Lock()
	self.targets[mapConfig.RemotePort] = t
	self.mappingLock.Unlock()
	return t, nil
}
