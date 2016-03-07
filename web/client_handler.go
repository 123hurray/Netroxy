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

type ClientHandler struct {
	webServer *NetroxyWebServer
}

func (self ClientHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	webServer := self.webServer
	name := r.FormValue("name")
	if name == "" {
		logger.Debug("WebPage:/client, name not found.")
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	webServer.lock.RLock()
	var client ClientModel
	for _, i := range webServer.serverModels {
		client = i.GetClient(name)
		if client != nil {
			break
		}
	}
	webServer.lock.RUnlock()
	if client == nil {
		logger.Debug("WebPage:/client, name not found.")
		w.WriteHeader(http.StatusNotFound)
		return
	}
	t := template.New("client.html")
	_, err := t.ParseFiles(self.webServer.GetFile("templates/client.html"), self.webServer.GetFile("templates/header.html"), self.webServer.GetFile("templates/mappings_table.html"))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	err = t.Execute(w, client)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		logger.Error(err)
	}
}
