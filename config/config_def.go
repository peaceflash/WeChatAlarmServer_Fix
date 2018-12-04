package config


//配置文件相关结构定义
import (


        goconf "github.com/zsounder/goconf"
)
type Config struct {
    goconf.AutoOptions
    ServerPort string `cfg:"ServerPort" default:"3322"`
    Corpid string `cfg:"Corpid" default:"wxd42d1e7afb9d6cf8"`
    Corpsecret string `cfg:"Corpsecret" default:""`
    Agentid  int64 `cfg:"Agentid" default:"4"`
    Touser string `cfg:"Touser" default:"@all"`
    Toparty string `cfg:"Toparty" default:""`
    Totag string `cfg:"Totag" default:""`
    Safe int64 `cfg:"Safe" default:"0"`
    Msgtype string `cfg:"Msgtype" default:"text"`
}
