package control

import (
	"sync"
)

const (
	SENDFAILED  = -1
	SENDSUCCESS = 0
)

//-----------------------------

//返回的响应
type Resp struct {
	State State `json:"state"`
}

type State struct {
	Rc  int    `json:"rc"`
	Msg string `json:"msg"`
}

var token requestToken
var locker sync.RWMutex

//--------------------------------

//--------------------------
var gettokenurl string = `https://qyapi.weixin.qq.com/cgi-bin/gettoken?corpid=`

var sendmessageurl = `https://qyapi.weixin.qq.com/cgi-bin/message/send?access_token=`

type responseErr struct {
	Errcode int    `json:"errcode`
	Errmsg  string `json:"errmsg"`
}

type requestToken struct {
	responseErr
	Access_token string `json:"access_token"`
	Expires_in   int    `json:"expires_in"`
}

type requestMesage struct {
	Touser  string            `json:"touser"`
	Toparty string            `json:"toparty"`
	Totag   string            `json:"totag"`
	Msgtype string            `json:"msgtype"`
	Agentid int64             `json:"agentid"`
	Text    map[string]string `json:"text"`
	Safe    int64             `json:"safe"`
}
