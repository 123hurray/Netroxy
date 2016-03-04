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

package client

import (
	"fmt"
	"net"
	"strconv"
)

func (self *Client) auth(cliName string, username string, password string) {
	self.send("ATH\n" + cliName + "\n" + username + "\n" + password + "\n")
}
func (self *Client) superviseRequest() {
	self.send("SRQ\n")
}
func (self *Client) channelResponse(conn net.Conn, port int, token string) {
	portStr := strconv.Itoa(port)
	conn.Write([]byte("TRS\n" + token + "\n" + portStr + "\n"))
}
func (self *Client) mapRequest(remotePort int, address string, isOpen bool) {
	self.send("MAP\n" + strconv.Itoa(remotePort) + "\n" + address + "\n" + fmt.Sprintf("%t\n", isOpen))
}
