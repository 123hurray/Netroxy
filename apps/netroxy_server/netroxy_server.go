package main

import (
	"github.com/123hurray/netroxy/config"
	"github.com/123hurray/netroxy/server"
	"github.com/123hurray/netroxy/utils/logger"
	"github.com/123hurray/netroxy/utils/network"
)

func main() {
	logger.Start(logger.LOG_LEVEL_DEBUG, "")
	conf := new(server.ServerConfig)
	err := config.Parse("server_config.json", conf)
	if err != nil {
		logger.Fatal(err)
	}
	s, err := network.NewTcpServer("Netroxy_main", conf.Ip, conf.Port)
	handler := server.NewHandler()
	s.Serve(handler)
}
