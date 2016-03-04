package server

import (
	"net"
	"time"

	"github.com/123hurray/netroxy/common"
)

type ClientConn struct {
	expireTime time.Time
	name       string
	mappings   map[int]*common.Mapping
	conn       net.Conn
	timeout    int
}

func NewClientConn(conn net.Conn, name string, timeout int) *ClientConn {
	cli := new(ClientConn)
	cli.conn = conn
	cli.name = name
	cli.mappings = make(map[int]*common.Mapping)
	cli.timeout = timeout
	cli.expireTime = time.Now().Add(time.Duration(timeout) * time.Second)
	return cli
}

func (self *ClientConn) AddMapping(mapping *common.Mapping) {
	self.mappings[mapping.RemotePort] = mapping
}

func (self *ClientConn) UpdateExpireTime() {
	self.expireTime = time.Now().Add(time.Duration(self.timeout) * time.Second)
}
