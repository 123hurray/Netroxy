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
	"errors"
	"io"
	"math/rand"
	"net"
	"os"
	"strconv"
	"sync/atomic"
	"time"

	"github.com/123hurray/netroxy/common"
	"github.com/123hurray/netroxy/utils/logger"
)

const defaultBufferSize = 16 * 1024

type Client struct {
	conn       net.Conn
	targets    map[int]*common.Mapping
	reader     *bufio.Reader
	ip         string
	port       int
	exitChan   chan bool
	expireTime int32
	timeout    int
	name       string
}

func NewClient(ip string, port int) *Client {
	client := new(Client)
	client.ip = ip
	client.port = port
	client.targets = make(map[int]*common.Mapping)
	client.exitChan = make(chan bool)
	client.expireTime = 0
	name, err := os.Hostname()
	if err != nil {
		logger.Fatal("Cannot get Hostname.")
	}
	client.name = name + "-" + strconv.Itoa(rand.Int())
	return client
}

func (self *Client) send(content string) {
	self.conn.Write([]byte(content))
}

func (self *Client) Login(username string, password string) error {
	conn, err := net.Dial("tcp", self.ip+":"+strconv.Itoa(self.port))
	if err != nil {
		return err
	}
	self.conn = conn
	self.auth(self.name, username, password)
	self.reader = bufio.NewReaderSize(conn, defaultBufferSize)
	ars, err := self.reader.ReadString('\n')
	if err != nil {
		return err
	}
	if ars != "ARS\n" {
		logger.Warn("Illegal command, expected ARS but receive", ars)
		return errors.New("Illegal command")
	}
	isOK, err := self.reader.ReadString('\n')
	if err != nil {
		return err
	}
	if isOK == "true\n" {
		timeoutStr, err := self.reader.ReadString('\n')
		if err != nil {
			return err
		}
		timeoutStr = timeoutStr[:len(timeoutStr)-1]
		timeout, err := strconv.Atoi(timeoutStr)
		if err != nil {
			logger.Warn("Illegal parameter, expected timeout but receive", timeoutStr)
			return err
		}
		self.timeout = timeout
		logger.Debug("Timeout:", timeout)
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
				self.Close()
			} else {
				self.superviseRequest()
			}
		}
	}
	<-self.exitChan
	logger.Info("Supervise stop.")
}

func (self *Client) Close() {
	err := self.conn.Close()
	if err != nil {
		close(self.exitChan)
	}
}

func (self *Client) handle() {
	defer self.Close()
	for {
		line, err := self.reader.ReadString('\n')
		if err != nil {
			return
		}
		line = line[:len(line)-1]
		switch {
		case line == "SRS":
			atomic.StoreInt32(&self.expireTime, 0)
		case line == "MRS":
			line, err = self.reader.ReadString('\n')
			line = line[:len(line)-1]
			remotePort, err := strconv.Atoi(line)
			if err != nil {
				logger.Warn("Illegal parament, expected int but receive", line, err)
				break
			}
			t := self.targets[remotePort]
			if t != nil {
				logger.Info("Mapping", self.ip+":"+line, "<->", t.Addr(), "accepted.")
			}
		case line == "TRQ":
			logger.Info("Tunnel request.")
			addr := self.ip + ":" + strconv.Itoa(self.port)
			line, err = self.reader.ReadString('\n')
			line = line[:len(line)-1]
			remotePort, err := strconv.Atoi(line)
			if err != nil {
				logger.Warn("Illegal parament, expected int but receive", line, err)
				break
			}
			t := self.targets[remotePort]
			if t != nil {
				logger.Info("New tunnel", self.ip+":"+strconv.Itoa(remotePort), "<->", t.Addr(), "Establishing...")
				conn1, err := net.Dial("tcp", addr)
				if err != nil {
					logger.Warn("Cannot connect to", addr, err)
					break
				}
				conn2, err := net.Dial("tcp", t.Addr())
				if err != nil {
					logger.Warn("Cannot connect to", t.Addr(), err)
					break
				}
				self.channelResponse(conn1, line)
				logger.Info("Dial " + t.Addr() + " OK")
				logger.Info("New tunnel", self.ip+":"+strconv.Itoa(remotePort), "<->", t.Addr(), "created.")
				go func() {
					io.Copy(conn1, conn2)
					defer conn1.Close()
				}()
				go func() {
					io.Copy(conn2, conn1)
					defer conn2.Close()
				}()
			}
		default:
			logger.Warn("Illegal command:", line)
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
	self.targets[mapConfig.RemotePort] = t
	return t, nil
}
