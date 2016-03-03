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
		cli := client.NewClient(conf.Ip, conf.Port)
		err = cli.Login()
		if err != nil {
			logger.Warn("Failed to connect server.", err)
			time.Sleep(3 * time.Second)
			continue
		}
		for _, i := range conf.Connections {
			cli.Connect(i.Ip, i.Port, i.RemotePort)
		}
		cli.Wait()
		logger.Warn("Connection to server closed. Reconnecting...")
	}
}
