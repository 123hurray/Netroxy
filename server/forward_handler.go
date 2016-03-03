package server

import (
	"net"

	"github.com/123hurray/netroxy/utils/logger"
)

type ForwardHandler struct {
	proxy *ProxyHandler
}

func (self *ForwardHandler) Handle(conn net.Conn) {
	self.proxy.connChan <- conn
	logger.Info("Forward connection ready.")
}
