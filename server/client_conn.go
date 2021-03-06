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
	"net"
	"sync"
	"time"
)

type ClientConn struct {
	expireTime   time.Time
	name         string
	token        string
	handlers     map[int]*ProxyHandler
	conn         net.Conn
	timeout      int
	loginTime    string
	handlersLock sync.RWMutex
	clientLock   sync.RWMutex
}

func NewClientConn(conn net.Conn, name string, token string, timeout int) *ClientConn {
	cli := new(ClientConn)
	cli.conn = conn
	cli.name = name
	cli.loginTime = time.Now().Format("01-02 15:04:05")
	cli.token = token
	cli.handlers = make(map[int]*ProxyHandler)
	cli.timeout = timeout
	cli.expireTime = time.Now().Add(time.Duration(timeout) * time.Second)
	return cli
}

func (self *ClientConn) AddHandler(handler *ProxyHandler) {
	self.handlersLock.Lock()
	defer self.handlersLock.Unlock()
	self.handlers[handler.mapping.RemotePort] = handler
}

func (self *ClientConn) RemoveHandler(key int) {
	self.handlersLock.Lock()
	defer self.handlersLock.Unlock()
	delete(self.handlers, key)
}

func (self *ClientConn) GetHandler(key int) *ProxyHandler {
	self.handlersLock.RLock()
	defer self.handlersLock.RLocker()
	return self.handlers[key]
}

func (self *ClientConn) UpdateExpireTime() {
	self.expireTime = time.Now().Add(time.Duration(self.timeout) * time.Second)
}
