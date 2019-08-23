package dingtalk

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/rancher/receiver/pkg/providers"
	"io/ioutil"
	"log"
	"net/http"
)

const (
	dingtalkErrorCode = "errcode"
	webhookURL = "https://oapi.dingtalk.com/robot/send?access_token=xxxxxxxx"
	url = "https://oapi.dingtalk.com"
	accessToken = "7200s"
)

type sender struct {
	webhookURL string
}

// TODO error more detail
func (s *sender) Send(msg string, receiver providers.Receiver) error {
	payload := newPayload(msg)
	req, err := http.NewRequest(http.MethodPost, s.webhookURL, bytes.NewReader(payload))
	if err != nil {
		return err
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	respData, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	dtr := dingtalkResp{}
	if err := json.Unmarshal(respData, &dtr); err != nil {
		return err
	}
	if dtr.ErrCode != 0 {
		return fmt.Errorf("dingtalk response errcode:%d", dtr.ErrCode)
	}

	return nil
}

type payload struct {
	MsgType string `json:"msgtype"`
	Text struct{
		Content string `json:"content"`
	} `json:"text"`
	At struct {
		IsAtAll bool `json:"isAtAll"`
	} `json:"at"`
}

type dingtalkResp struct {
	ErrCode int `json:"errcode"`
	ErrMsg string `json:"errmsg"`
}

func newPayload(msg string) []byte {
	p := payload{
		MsgType: "text",
		Text: struct {
			Content string `json:"content"`
		}{
			Content: msg,
		},
		At: struct {
			IsAtAll bool `json:"isAtAll"`
		}{
			IsAtAll: true,
		},
	}

	data, _ := json.Marshal(p)
	return data
}

func getAccessToken(appkey, appsecret string) {
	req, err := http.NewRequest(http.MethodGet, "https://oapi.dingtalk.com/gettoken",  nil)
	if err != nil {
		log.Println(err)
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Println(err)
	}
	log.Println(resp)
}
