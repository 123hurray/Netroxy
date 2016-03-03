package server

type ServerConfig struct {
	Ip       string `json:"ip"`
	Port     int    `json:"port"`
	Username string `json:"username"`
	Password string `json:"password"`
}
