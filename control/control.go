package control

import (
	"bytes"
	"encoding/json"
	"errors"

	"github.com/peaceflash/WeChatAlarmServer_Fix/config"

	"time"

	"io/ioutil"
	"net/http"

	l4g "github.com/alecthomas/log4go"
)

func Init() {
	getToken()
	//定时更新token
	go func() {
		tick := time.NewTicker(time.Minute * 1)
		for {
			select {
			case <-tick.C:
				{
					getToken()
				}
			}
		}
	}()
}

func MessageNotify(message string) (resp Resp) {
	var msg requestMesage
	//发送消息
	msg.Touser = config.GetTouser()
	msg.Totag = config.GetTotag()
	msg.Toparty = config.GetToparty()
	msg.Safe = config.GetSafe()
	msg.Msgtype = config.GetMsgtype()
	msg.Text = map[string]string{"content": message}
	msg.Agentid = config.GetAgentid()

	msgbuf, err := json.Marshal(msg)
	if err != nil {
		l4g.Error("json err:", err)
		resp.State.Rc = SENDFAILED
		resp.State.Msg = "json msg failed"
		return
	}
	err = sendMessage(msgbuf)
	if err != nil {
		l4g.Error("sendmessage faile err:", err)
		resp.State.Rc = SENDFAILED
		resp.State.Msg = "sendMessage failed"
		return
	}
	resp.State.Rc = SENDSUCCESS
	resp.State.Msg = "sendMessage success"

	return
}

//--------------------------------------

func getToken() (err error) {
	corpid := config.GetCorpid()
	corpsecret := config.GetCorpsecret()
	var url string = gettokenurl + corpid + "&corpsecret=" + corpsecret
	resp, err := http.Get(url)
	if err != nil {
		l4g.Error("get token failed:", err)
	}

	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		l4g.Error("get token failed:", err)
	}
	body, _ := ioutil.ReadAll(resp.Body)

	locker.Lock()
	defer locker.Unlock()
	err = json.Unmarshal([]byte(body), &token)
	if err != nil {
		l4g.Error("json token failed:", err)
	}
	return
}

func sendMessage(msg []byte) error {

	body := bytes.NewBuffer(msg)
	locker.RLock()
	accesstoken := token.Access_token
	locker.RUnlock()

	resp, err := http.Post(sendmessageurl+accesstoken, "application/json", body)
	if resp.StatusCode != 200 {
		return errors.New("request error")
	}
	buf, _ := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	var e responseErr
	err = json.Unmarshal(buf, &e)
	if err != nil {
		return err
	}
	if e.Errcode != 0 && e.Errmsg != "ok" {
		return errors.New(string(buf))
	}
	return nil
}
