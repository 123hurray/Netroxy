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

package main

import (
	"sync"
	"time"

	"github.com/123hurray/netroxy/config"
	"github.com/123hurray/netroxy/server"
	"github.com/123hurray/netroxy/utils/logger"
	"github.com/123hurray/netroxy/utils/network"
	"github.com/123hurray/netroxy/utils/security"
	"github.com/123hurray/netroxy/web"
)

func main() {
	logger.Start(logger.LOG_LEVEL_DEBUG, "")
	conf := new(server.ServerConfig)
	err := config.Parse("server_config.json", conf)
	if err != nil {
		logger.Fatal(err)
	}
	plainServer, err := network.NewPlainServer("Netroxy_main", conf.Ip, conf.Port)
	if err != nil {
		logger.Fatal(err)
	}
	plainNetroxyServer := server.NewServer(conf, "PlainServer-"+security.GenerateUID(8), false)
	go func() {
		ticker := time.NewTicker(time.Duration(conf.Timeout/3) * time.Second)
		for {
			select {
			case <-ticker.C:
				plainNetroxyServer.Supervise()
			}
		}
	}()
	tlsServer, err := network.NewTLSServer("Netroxy_main_TLS", conf.Ip, conf.TLS.Port, conf.TLS.Ca, conf.TLS.Key)
	if err != nil {
		logger.Fatal(err)
	}
	tlsNetroxyServer := server.NewServer(conf, "TLSServer-"+security.GenerateUID(8), true)
	go func() {
		ticker := time.NewTicker(time.Duration(conf.Timeout/3) * time.Second)
		for {
			select {
			case <-ticker.C:
				tlsNetroxyServer.Supervise()
			}
		}
	}()
	webServer := web.NewNetroxyWebServer([]web.ServerModel{tlsNetroxyServer, plainNetroxyServer}, &conf.Web)
	var wg sync.WaitGroup
	wg.Add(3)
	go func() {
		go tlsNetroxyServer.WebDemon()
		tlsServer.Serve(tlsNetroxyServer)
		wg.Done()
	}()
	go func() {
		go plainNetroxyServer.WebDemon()
		plainServer.Serve(plainNetroxyServer)
		wg.Done()
	}()
	go func() {
		webServer.Serve()
		wg.Done()
	}()
	wg.Wait()
}
