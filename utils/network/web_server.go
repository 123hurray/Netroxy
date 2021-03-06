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

package network

import (
	"net"
	"net/http"
	"strconv"
)

type WebServer struct {
	ip    string
	port  int
	Root  string
	https bool
	ca    string
	key   string
}

func NewWebServer(ip string, port int, root string, https bool, ca string, key string) *WebServer {
	return &WebServer{ip, port, root, https, ca, key}
}

func (self *WebServer) Serve(handlers map[string]http.Handler) {
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir(self.Root+"static/"))))
	for path, handler := range handlers {
		http.Handle(path, handler)
	}
	addr := net.JoinHostPort(self.ip, strconv.Itoa(self.port))
	if self.https == true {
		http.ListenAndServeTLS(addr, self.ca, self.key, nil)
	} else {
		http.ListenAndServe(addr, nil)
	}
}
