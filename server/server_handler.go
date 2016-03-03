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
}

func NewHandler() *Handler {
	handler := new(Handler)
	handler.portHandlerDict = make(map[int]*ProxyHandler)
	handler.portTcpServerDict = make(map[int]*network.Server)
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
		logger.Debug(line)
		if err != nil {
			logger.Warn("Command illegal.", err)
			return
		}
		line = line[:len(line)-1]
		switch {
		case line == "Auth":
		case line == "NewConn":
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
			conn.Write([]byte("NewConnOK\n" + line + "\n"))

		case line == "NewRes":
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
