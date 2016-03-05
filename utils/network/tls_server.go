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
	"crypto/tls"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"io/ioutil"
	"net"
	"strconv"

	"github.com/123hurray/netroxy/utils/logger"
)

type tlsServer struct {
	ip     string
	port   int
	name   string
	socket net.Listener
}

func getServerConfig(caFile string, keyFile string) (*tls.Config, error) {
	ca_b, _ := ioutil.ReadFile(caFile)
	block, _ := pem.Decode(ca_b)
	if block == nil {
		return nil, errors.New("CA file not found.")
	}
	ca, _ := x509.ParseCertificate(block.Bytes)
	priv_b, _ := ioutil.ReadFile(keyFile)
	if priv_b == nil {
		return nil, errors.New("Key file not found.")
	}
	privBlock, _ := pem.Decode(priv_b)
	priv, _ := x509.ParsePKCS1PrivateKey(privBlock.Bytes)
	cert := tls.Certificate{
		Certificate: [][]byte{block.Bytes},
		PrivateKey:  priv,
	}
	pool := x509.NewCertPool()
	pool.AddCert(ca)
	config := tls.Config{
		ClientAuth:   tls.NoClientCert,
		Certificates: []tls.Certificate{cert},
		ClientCAs:    pool,
	}
	return &config, nil
}

// return a new TLS server
func NewTLSServer(name string, ip string, port int, caFile string, keyFile string) (TCPServer, error) {
	serverConfig, err := getServerConfig(caFile, keyFile)
	if err != nil {
		logger.Error(err)
		return nil, err
	}
	socket, err := tls.Listen("tcp", net.JoinHostPort(ip, strconv.Itoa(port)), serverConfig)
	if err != nil {
		logger.Error(err)
		return nil, err
	}
	server := tlsServer{ip, port, name, socket}
	return &server, err
}

// Start TLS server with a handler
func (self *tlsServer) Serve(handler Handler) {
	logger.Info(self.name, "Listening", self.ip, ":", self.port)

	for {
		conn, err := self.socket.Accept()
		if err != nil {
			logger.Error(err)
			return
		}
		go handler.Handle(conn)
	}
}

// Close TLS server
func (self *tlsServer) Close() {
	self.socket.Close()
}
