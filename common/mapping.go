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

package common

import (
	"strconv"
	"sync"
)

type Mapping struct {
	Ip         string
	Port       int
	RemotePort int
	isOn       bool
	lock       sync.RWMutex
}

func NewMapping(ip string, port int, remotePort int, isOn bool) *Mapping {
	mapping := Mapping{}
	mapping.Ip = ip
	mapping.Port = port
	mapping.RemotePort = remotePort
	mapping.isOn = isOn
	return &mapping
}

func (self *Mapping) Addr() string {
	return self.Ip + ":" + strconv.Itoa(self.Port)
}
func (self *Mapping) TurnOn() {
	self.lock.Lock()
	defer self.lock.Unlock()
	self.isOn = true
}
func (self *Mapping) TurnOff() {
	self.lock.Lock()
	defer self.lock.Unlock()
	self.isOn = false
}

func (self *Mapping) IsOn() bool {
	self.lock.RLock()
	defer self.lock.RUnlock()
	return self.isOn
}
