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
	"html/template"
	"net/http"

	"github.com/123hurray/netroxy/utils/logger"
)

type ServersHandler struct {
	webServer *NetroxyWebServer
}

func (self ServersHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	webServer := self.webServer
	servers := webServer.serverModels
	t := template.New("servers.html")
	_, err := t.ParseFiles(self.webServer.GetFile("templates/servers.html"), self.webServer.GetFile("templates/header.html"))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		logger.Error(err)
		return
	}
	err = t.Execute(w, servers)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		logger.Error(err)
	}
}
