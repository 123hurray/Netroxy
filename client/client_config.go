package client

type ClientConfig struct {
	Ip          string `json:"ip"`
	Port        int    `json:"port"`
	Username    string `json:"username"`
	Password    string `json:"password"`
	Connections []struct {
		Ip         string `json:"ip`
		Port       int    `json:port`
		RemotePort int    `json:remotePort`
	}
}
