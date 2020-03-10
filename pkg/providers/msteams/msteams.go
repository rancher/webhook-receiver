package msteams

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"

	"github.com/rancher/webhook-receiver/pkg/providers"
)

const (
	Name = "MICROSOFT_TEAMS"

	webhookURLKey = "webhook_url"
	proxyURLKey   = "proxy_url"
)

type sender struct {
	webhookURL string
	proxyURL   string

	client *http.Client
}

func New(opt map[string]string) (providers.Sender, error) {
	if err := validate(opt); err != nil {
		return nil, err
	}

	c := &http.Client{}

	return &sender{
		webhookURL: opt[webhookURLKey],
		proxyURL:   opt[proxyURLKey],
		client:     c,
	}, nil
}

func (s *sender) Send(msg string, receiver providers.Receiver) error {
	payload, err := newPayload(msg)
	if err != nil {
		return err
	}

	req, err := http.NewRequest(http.MethodPost, s.webhookURL, bytes.NewReader(payload))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	transport := &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: false,
		},
	}

	if s.proxyURL != "" {
		proxy := func(_ *http.Request) (*url.URL, error) {
			return url.Parse(s.proxyURL)
		}

		transport.Proxy = proxy
	}

	s.client.Transport = transport

	resp, err := s.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	respData, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	if string(respData) != "1" {
		return fmt.Errorf("microsoft teams response err:%s", string(respData))
	}

	return nil
}

type payload struct {
	Text string `json:"text"`
}

func newPayload(msg string) ([]byte, error) {
	p := payload{
		Text: msg,
	}

	data, err := json.Marshal(p)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal fields to JSON, %v", err)
	}

	return data, nil
}

func validate(opt map[string]string) error {
	if _, exists := opt[webhookURLKey]; !exists {
		return fmt.Errorf("%s empty", webhookURLKey)
	}

	return nil
}
