package server

import (
	"net"
	"time"
)

type ClientConn struct {
	expireTime time.Time
	conn       net.Conn
	timeout    int
}

func NewClientConn(conn net.Conn, timeout int) *ClientConn {
	cli := new(ClientConn)
	cli.conn = conn
	cli.timeout = timeout
	cli.expireTime = time.Now().Add(time.Duration(timeout) * time.Second)
	return cli
}

func (self *ClientConn) UpdateExpireTime() {
	self.expireTime = time.Now().Add(time.Duration(self.timeout) * time.Second)
}
