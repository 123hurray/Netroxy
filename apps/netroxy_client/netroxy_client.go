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
	"time"

	"github.com/123hurray/netroxy/client"
	"github.com/123hurray/netroxy/config"
	"github.com/123hurray/netroxy/utils/logger"
)

func main() {
	logger.Start(logger.LOG_LEVEL_DEBUG, "")
	conf := new(client.ClientConfig)
	err := config.Parse("client_config.json", conf)
	if err != nil {
		logger.Fatal(err)
	}
	for {
		cli := client.NewClient(conf)
		err = cli.Login()
		if err != nil {
			logger.Warn("Failed to connect server.", err)
			time.Sleep(3 * time.Second)
			continue
		}
		for _, i := range conf.Connections {
			cli.Connect(&i)
		}
		cli.Wait()
		logger.Warn("Connection to server closed. Reconnecting...")
	}
}
