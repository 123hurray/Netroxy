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

type ClientsHandler struct {
	webServer *NetroxyWebServer
}

func (self ClientsHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	webServer := self.webServer
	webServer.lock.RLock()
	var clients []ClientModel
	for _, i := range webServer.serverModels {
		clients = append(clients, i.GetClients()...)
	}
	webServer.lock.RUnlock()
	t := template.New("clients.html")
	_, err := t.ParseFiles(self.webServer.GetFile("templates/clients.html"), self.webServer.GetFile("templates/header.html"), self.webServer.GetFile("templates/clients_table.html"))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		logger.Error(err)
		return
	}
	err = t.Execute(w, clients)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		logger.Error(err)
	}
}
