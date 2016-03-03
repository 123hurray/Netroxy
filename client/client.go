package client

import (
	"bufio"
	"io"
	"net"
	"strconv"

	"github.com/123hurray/netroxy/utils/logger"
)

const defaultBufferSize = 16 * 1024

type Client struct {
	conn     net.Conn
	targets  map[int]*Target
	reader   *bufio.Reader
	ip       string
	port     int
	exitChan chan bool
}

func NewClient(ip string, port int) *Client {
	client := new(Client)
	client.ip = ip
	client.port = port
	client.targets = make(map[int]*Target)
	client.exitChan = make(chan bool)
	return client
}

func (self *Client) send(content string) {
	self.conn.Write([]byte(content))
}

func (self *Client) Login() error {
	conn, err := net.Dial("tcp", self.ip+":"+strconv.Itoa(self.port))
	if err != nil {
		return err
	}
	logger.Info("Login to server success.")
	self.conn = conn
	self.send("Auth\n")
	self.reader = bufio.NewReaderSize(conn, defaultBufferSize)
	go self.handle()
	return nil
}
func (self *Client) handle() {
	defer self.conn.Close()
	defer close(self.exitChan)
	for {
		line, err := self.reader.ReadString('\n')
		if err != nil {
			return
		}
		line = line[:len(line)-1]
		switch {
		case line == "NewConnOK":
			line, err = self.reader.ReadString('\n')
			line = line[:len(line)-1]
			remotePort, err := strconv.Atoi(line)
			if err != nil {
				logger.Warn("Illegal parament, expected int:", line, err)
				break
			}
			t := self.targets[remotePort]
			if t != nil {
				logger.Info("Connection " + t.Addr() + " prepared.")
			}
		case line == "Dial":
			line, err = self.reader.ReadString('\n')
			line = line[:len(line)-1]
			remotePort, err := strconv.Atoi(line)
			if err != nil {
				break
			}
			t := self.targets[remotePort]
			if t != nil {
				addr := self.ip + ":" + strconv.Itoa(self.port)

				conn1, err := net.Dial("tcp", addr)
				if err != nil {
					break
				}
				defer conn1.Close()
				logger.Info("Dial " + addr + " OK")
				conn2, err := net.Dial("tcp", t.Addr())
				if err != nil {
					logger.Warn("Cannot connect to", t.Addr(), err)
					break
				}
				defer conn2.Close()
				conn1.Write([]byte("NewRes\n" + line + "\n"))
				logger.Info("Dial " + t.Addr() + " OK")
				logger.Info("Connection " + addr + " -> " + t.Addr() + " OK")
				exitChan := make(chan bool)
				go func() {
					io.Copy(conn1, conn2)
					_, isClose := <-exitChan
					if !isClose {
						close(exitChan)
					}
				}()
				go func() {
					io.Copy(conn2, conn1)
					_, isClose := <-exitChan
					if !isClose {
						close(exitChan)
					}
				}()
				<-exitChan
			}
		}
	}
}

func (self *Client) Wait() {
	<-self.exitChan
}
func (self *Client) Connect(ip string, port int, remotePort int) (*Target, error) {
	addr := ip + ":" + strconv.Itoa(port)
	logger.Info("Connection " + addr + " Establishing...")
	self.send("NewConn\n" + strconv.Itoa(remotePort) + "\n")
	t := NewTarget(ip, port)
	self.targets[remotePort] = t
	return t, nil
}
