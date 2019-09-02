// we just support dingtalk robot now
package dingtalk

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/rancher/receiver/pkg/providers"
)

const (
	Name = "dingtalk"

	webhookURLKey = "webhook_url"
)

type sender struct {
	webhookURL string
	client     *http.Client
}

func New(opt map[string]string) (providers.Sender, error) {
	if err := validate(opt); err != nil {
		return nil, err
	}

	c := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		},
	}
	return &sender{
		webhookURL: opt[webhookURLKey],
		client:     c,
	}, nil
}

// TODO error more detail
func (s *sender) Send(msg string, receiver providers.Receiver) error {
	payload := newPayload(msg)
	req, err := http.NewRequest(http.MethodPost, s.webhookURL, bytes.NewReader(payload))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")

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
	Text    struct {
		Content string `json:"content"`
	} `json:"text"`
	At struct {
		IsAtAll bool `json:"isAtAll"`
	} `json:"at"`
}

type dingtalkResp struct {
	ErrCode int    `json:"errcode"`
	ErrMsg  string `json:"errmsg"`
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

func validate(opt map[string]string) error {
	if _, exists := opt[webhookURLKey]; !exists {
		return fmt.Errorf("%s empty", webhookURLKey)
	}

	return nil
}
