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
	"bufio"
	"errors"
	"strconv"
)

type ProtocolReader struct {
	reader *bufio.Reader
}

func (self *ProtocolReader) SetReader(reader *bufio.Reader) {
	self.reader = reader
}

func (self *ProtocolReader) GetString() (str string, err error) {
	str, err = self.reader.ReadString('\n')
	if err != nil {
		return
	}
	str = str[:len(str)-1]
	return
}
func (self *ProtocolReader) GetInt() (i int, err error) {
	str, err := self.GetString()
	if err != nil {
		return
	}
	i, err = strconv.Atoi(str)
	if err != nil {
		return
	}
	return
}
func (self *ProtocolReader) GetBool() (b bool, err error) {
	str, err := self.GetString()
	if err != nil {
		return
	}
	if str == "true" {
		b = true
	} else if str == "false" {
		b = false
	} else {
		err = errors.New("Illegal bool value. Receive:" + str)
	}
	return
}
