
# Introduction

[![license](http://img.shields.io/badge/license-MIT-red.svg?style=flat)](https://raw.githubusercontent.com/123hurray/tlslog/master/LICENSE)
[![travis-ci](https://api.travis-ci.org/123hurray/Netroxy.svg)](https://travis-ci.org/123hurray/Netroxy)

Netroxy is a TCP proxy that can be used to access LAN services from any other LAN, for example "Remote Desktop service".

Netroxy consists of two components: **netroxy_server** and **netroxy_client**. netroxy_server runs on a server that can be accessed by netroxy_client and user. netroxy_client runs on a client machine which has the services that user need. Once the connection from netroxy_client to netroxy_server established, user will have the ability to access services throught netroxy_server.

# Getting started

## Installation

You can build netroxy yourself or download binary release.

## Build

```shell
# Install Golang and set GOPATH first

# Build netroxy_server

go get github.com/123hurray/netroxy/apps/netroxy_server

go install src/github.com/123hurray/netroxy/apps/netroxy_server

# Build netroxy_client

go get github.com/123hurray/netroxy/apps/netroxy_client

go install src/github.com/123hurray/netroxy/apps/netroxy_client

# All things done!
```

## Usage

### Server

Modify `server_config.json` and run `netroxy_server`.

### Client

Modify `client_config.json` and run `netroxy_client`.

# Internals

Server listens on an address(IpA:PortA) to wait client connection. When a client connected, it tells server which address(IpB:PortB) it wants to map and which server port(PortC) it wants server to listen. The server then listens on the new port(PortC). 

When user connects to server's port(PortC), server sends a request to associated client and client makes connection to destination address(IpB:PortB, the real service user wants to connect) and send response to server address(IpA:PortA). Finally, tunnel established.

1. Server start:
    
    Listen on IpA:PortA

2. Client connection:

    IpD:PortD (Auth)-> IpA:PortA
	
    IpA:PortA (Auth OK)-> IpD:PortD
    
3. Mapping request:

    IpD:PortD (Mapping request, map IpB:PortB to IpA:PortC)-> IpA:PortA

4. Mapping response:
    
    Listen on IpA:PortC
	
    IpA:PortA (Mapping OK)-> IPD:PortD

5. User request:

    IpE:PortE (request)-> IpA:PortC
    
6. Tunnel request:

    IpA:PortA (tunnel request, IpB:PortB)-> IpD:PortD

7. Tunnel response:

    IpD:PortF <-> IpB:PortB
	
    IpD:PortG (tunnel response, IpB:PortB)-> IpA:PortA
    
8. Tunnel established:

    IO Copy: from (IpE:PortE <-> IpA:PortC) to (IpA:PortA <-> IpD:PortG)
	
    IO Copy: from (IpA:PortA <-> IpD:PortG) to (IpE:PortE <-> IpA:PortC)
    
    IO Copy: from (IpD:PortG <-> IpA:PortA) to (IpD:PortF <-> IpB:PortB)
	
    IO Copy: from (IpD:PortF <-> IpB:PortB) to (IpD:PortG <-> IpA:PortA)
    
9. User accesses service:

    IpE:PortE <-> IpA:PortC <-> IpD:PortG <-> IpB:PortB
    
# Protocol v0.3

## Commands

### ATH

Connect and auth to server

    ATH\n
	clientName\n
    username\n
    password\n
    
### ARS

Auth response

    ARS\n
    isOK(true or false)\n
	timeout(Present if isOk is true)\n
	token(Present if isOk is true)\n
    
### MAP

Map a local tcp address to server port

    MAP\n
    port\n
	address\n
	isOpen(true or false)\n
	
    
### MRS

Map response
    
    MRS\n
    port\n
    isOK(true or false)\n

### TRQ

Tunnel request

    TRQ\n
    port\n
    
### TRS

Tunnel response.Tunnel response is send from a new
tcp socket. So it has to use token to tell server who it is.

    TRS\n
	token\n
    port\n   

### SRQ

Keepalive request from client

	SRQ\n

### SRS

Keepalive response from server

	SRS\n

# TODO

 - [x] Use SSL/TLS in client/server connection
 - [ ] Use SSL/TLS in user/server connection
 - [ ] Server/client can specify config file name
 - [x] Fix bug: One client disconnect from server will close all server ports
 - [ ] Web interface to view all mapped ports
 - [ ] User can send mapping requests

# License

Netroxy is licensed under the MIT License. The terms of the license are as follows:

    The MIT License (MIT)

    Copyright (c) 2016 Ray Zhang

    Permission is hereby granted, free of charge, to any person obtaining a copy
    of this software and associated documentation files (the "Software"), to deal
    in the Software without restriction, including without limitation the rights
    to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
    copies of the Software, and to permit persons to whom the Software is
    furnished to do so, subject to the following conditions:

    The above copyright notice and this permission notice shall be included in all
    copies or substantial portions of the Software.

    THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
    IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
    FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
    AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
    LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
    OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
    SOFTWARE.