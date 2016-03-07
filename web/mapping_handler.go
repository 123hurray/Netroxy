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
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/123hurray/netroxy/utils/logger"
)

type MappingHandler struct {
	webServer *NetroxyWebServer
}

func (self MappingHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	webServer := self.webServer
	action := r.FormValue("action")
	portStr := r.FormValue("port")
	if portStr == "" || action == "" {
		logger.Debug("WebPage:/mapping, illegal argument.")
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	port, err := strconv.Atoi(portStr)
	if err != nil {
		logger.Debug("WebPage:/mapping, illegal port.")
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	webServer.lock.RLock()
	defer webServer.lock.RUnlock()
	for _, i := range webServer.serverModels {
		if action == "on" {
			ok := i.TurnMappingOn(port)
			if ok {
				j, _ := json.Marshal(true)
				w.Write(j)
				return
			}
		} else {
			ok := i.TurnMappingOff(port)
			if ok {
				j, _ := json.Marshal(true)
				w.Write(j)
				return
			}
		}
	}
	j, _ := json.Marshal(false)
	w.Write(j)
}
