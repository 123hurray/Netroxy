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
	"github.com/123hurray/netroxy/utils/logger"
	"github.com/123hurray/netroxy/web"
)

func (self *Server) GetName() string {
	return self.name
}

func (self *Server) IsTLS() bool {
	return self.isTLS
}

func (self *Server) GetStartupTime() string {
	return self.startupTime
}

func (self *Server) GetClientNumber() int {
	self.clientsLock.RLock()
	defer self.clientsLock.RUnlock()
	return len(self.clients)
}

func (self *Server) GetMappingNumber() int {
	self.clientsLock.RLock()
	defer self.clientsLock.RUnlock()
	num := 0
	for _, cli := range self.clients {
		num += len(cli.handlers)
	}
	return num
}

func (self *Server) GetClients() (clis []web.ClientModel) {
	self.clientsLock.RLock()
	defer self.clientsLock.RUnlock()
	for _, cli := range self.clients {
		cli.clientLock.RLock()
		clis = append(clis, cli)
		cli.clientLock.RUnlock()
	}
	return
}
func (self *Server) GetClient(name string) (cli web.ClientModel) {
	self.clientsLock.RLock()
	defer self.clientsLock.RUnlock()
	if val, ok := self.clientsNameMap[name]; ok {
		return val
	} else {
		return nil
	}
}
func (self *Server) GetMappings() (mappings []web.MappingModel) {
	self.clientsLock.RLock()
	defer self.clientsLock.RUnlock()
	for _, cli := range self.clients {
		mappings = append(mappings, cli.GetMappings()...)
	}
	return
}
func (self *Server) GetMapping(port int) web.MappingModel {
	self.clientsLock.RLock()
	defer self.clientsLock.RUnlock()
	for _, cli := range self.clients {
		if mapping, ok := cli.handlers[port]; ok {
			return mapping
		}
	}
	return nil
}

func (self *Server) WebDemon() {
	for {
		select {
		case port := <-self.turnMappingOnCh:
			mapping := self.GetMapping(port)
			if mapping != nil {
				mapping.TurnOn()
				self.responseCh <- true
			} else {
				self.responseCh <- false
			}
		case port := <-self.turnMappingOffCh:
			mapping := self.GetMapping(port)
			if mapping != nil {
				mapping.TurnOff()
				self.responseCh <- true
			} else {
				self.responseCh <- false
			}
		}
	}
}

func (self *Server) TurnMappingOn(port int) bool {
	logger.Debug("Turning mapping", port, "Off")
	self.turnMappingOnCh <- port
	return <-self.responseCh
}

func (self *Server) TurnMappingOff(port int) bool {
	logger.Debug("Turning mapping", port, "Off")
	self.turnMappingOffCh <- port
	return <-self.responseCh
}
