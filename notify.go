package notify

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/spf13/viper"
)

var corpID, corpSecret string

const (
	sendURL     = `https://qyapi.weixin.qq.com/cgi-bin/message/send?access_token=%s`
	getTokenURL = `https://qyapi.weixin.qq.com/cgi-bin/gettoken?corpid=%s&corpsecret=%s`
)

func init() {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	err := viper.ReadInConfig() // Find and read the config file
	if err != nil {             // Handle errors reading the config file
		panic(fmt.Errorf("fatal error config file: %w", err))
	}
	log.Printf("reading configi file success: %s ", viper.GetString("name"))
	corpID = viper.GetString("corpID")
	corpSecret = viper.GetString("corpSecret")
}

type WeNotifier struct {
	token string
	msg   string
}

// new  instance
func NewWeNotifier(msg string) *WeNotifier {
	return &WeNotifier{msg: msg}
}

var defaultNotifier = NewWeNotifier("ping")

// wrapper
func Send() error {
	return defaultNotifier.Send()
}

// send the msg
func (wn *WeNotifier) Send() error {

	//if  not inited  , then  fresh token
	if wn.token == "" {
		wn.freshToken()
	}
	if err := wn.send(); err != nil {
		if err == errExpiredToken {
			wn.freshToken()
			wn.send()
		}
		return err
	}
	return nil
}

// doing actural send thing
func (wn *WeNotifier) send() error {
	fmt.Println("like sent")
	return nil
	u := fmt.Sprintf(sendURL, wn.token)
	s := `
	{
		"touser" : "@all",
		"msgtype"    : "text",
		"agentid" : 1000002,
		"text" : {
			"content" : "%s"
		},
		"safe":0
	 }
`
	s = fmt.Sprintf(s, wn.msg)
	var body = bytes.NewBufferString(s)
	resp, err := http.Post(u, "application/json", body)

	if err != nil {
		return err
	}
	defer resp.Body.Close()
	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil
	}
	var call = struct {
		ErrCode int    `json:"errcode"`
		ErrMsg  string `json:"errmsg"`
		MsgID   string `json:"msgid"`
	}{}

	err = json.Unmarshal(b, &call)
	if call.ErrCode == codeExpiredToken {
		return errExpiredToken
	}
	if err != nil {
		return err
	}
	return nil
}

// fresh the token
// call if  token is empty  ,  unvalid
func (wn *WeNotifier) freshToken() error {
	u := fmt.Sprintf(getTokenURL, corpID, corpSecret)
	resp, err := http.Get(u)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	var call = struct {
		ErrCode     int    `json:"errcode"`
		ErrMsg      string `json:"errmsg"`
		AccessToken string `json:"access_token"`
		ExpiresIn   int    `json:"expires_in"`
	}{}
	if err = json.Unmarshal(b, &call); err != nil {
		return err

	}
	wn.token = call.AccessToken
	return nil
}

var codeExpiredToken = 42001
var errExpiredToken error
