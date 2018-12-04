package config

import (

	"fmt"
	goconf "github.com/zsounder/goconf"
)
var CfgServer *Config = &Config{}

func Init() {
	fmt.Println("111111")
	goconf.MustResolve(CfgServer,"config/config.toml")
}


func GetServerConfig() *Config {
	return CfgServer
}

// ------------------------------------------
func GetServerPort() string {
	return GetServerConfig().ServerPort
}

func GetCorpid() string {
	return GetServerConfig().Corpid
}

func GetCorpsecret() string {
	return GetServerConfig().Corpsecret
}

func GetAgentid() int64 {
	return GetServerConfig().Agentid
}

func GetSafe() int64 {
	return GetServerConfig().Safe
}

func GetToparty() string {
	return GetServerConfig().Toparty
}

func GetTotag() string {
	return GetServerConfig().Totag
}

func GetTouser() string {
	return GetServerConfig().Touser
}

func GetMsgtype() string {
	return GetServerConfig().Msgtype
}

// ------------------------------------------
