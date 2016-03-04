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
	handlersLock sync.RWMutex
}

func NewClientConn(conn net.Conn, name string, token string, timeout int) *ClientConn {
	cli := new(ClientConn)
	cli.conn = conn
	cli.name = name
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
