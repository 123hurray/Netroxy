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

package web

import (
	"net/http"
	"sync"

	"github.com/123hurray/netroxy/utils/network"
)

type WebConfig struct {
	Enabled bool   `json:"enabled"`
	Ip      string `json:"ip"`
	Port    int    `json:"port"`
	Root    string `json:"root"`
	Https   struct {
		Ca  string `json:"ca"`
		Key string `json:"key"`
	}
}
type NetroxyWebServer struct {
	serverModels []ServerModel
	server       *network.WebServer
	lock         sync.RWMutex
}

func NewNetroxyWebServer(serverModels []ServerModel, conf *WebConfig) *NetroxyWebServer {
	self := NetroxyWebServer{}
	self.serverModels = serverModels
	self.server = network.NewWebServer(conf.Ip, conf.Port, conf.Root)
	return &self
}

func (self *NetroxyWebServer) Serve() {
	handlers := map[string]http.Handler{
		"/server/":   ServerHandler{self},
		"/servers/":  ServersHandler{self},
		"/client/":   ClientHandler{self},
		"/clients/":  ClientsHandler{self},
		"/mapping/":  MappingHandler{self},
		"/mappings/": MappingsHandler{self},
	}
	self.server.Serve(handlers)
}

func (self *NetroxyWebServer) GetFile(dir string) string {
	return self.server.Root + dir
}
