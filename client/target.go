package client

import (
	"strconv"
)

type Target struct {
	Ip   string
	Port int
}

func NewTarget(ip string, port int) *Target {
	return &Target{ip, port}
}

func (self *Target) Addr() string {
	return self.Ip + ":" + strconv.Itoa(self.Port)
}
