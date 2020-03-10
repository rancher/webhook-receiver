package aliyunsms

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"

	"github.com/rancher/webhook-receiver/pkg/providers"
)

const (
	Name = "ALIYUN_SMS"

	regionID        = "cn-hangzhou"
	accessKeyIDKey  = "access_key_id"
	signNameKey     = "sign_name"
	templateCodeKey = "template_code"
	proxyURLKey     = "proxy_url"
)

type sender struct {
	client       *sdk.Client
	signName     string
	templateCode string
	proxyURL     string
}

func (s *sender) Send(msg string, receiver providers.Receiver) error {
	if s.proxyURL != "" {
		s.client.SetHttpsProxy(s.proxyURL)
	}

	request := requests.NewCommonRequest()
	request.Method = http.MethodPost
	request.Scheme = "https"
	request.Domain = "dysmsapi.aliyuncs.com"
	request.Version = "2017-05-25"
	request.ApiName = "SendSms"
	request.QueryParams["RegionId"] = regionID
	request.QueryParams["PhoneNumbers"] = strings.Join(receiver.To, ",")
	request.QueryParams["SignName"] = s.signName
	request.QueryParams["TemplateCode"] = s.templateCode
	request.QueryParams["TemplateParam"] = fmt.Sprintf(`{"alert":"%s"}`, msg)
	request.SetContent([]byte(msg))

	res, err := s.client.ProcessCommonRequest(request)
	if err != nil {
		return fmt.Errorf("client process request err:%v", err)
	}
	if !res.IsSuccess() {
		return fmt.Errorf("send faliure, resp is :%s", res.GetHttpContentString())
	}

	return nil
}

func New(opt map[string]string) (providers.Sender, error) {
	if err := validate(opt); err != nil {
		return nil, err
	}

	client, err := sdk.NewClientWithAccessKey(regionID, opt[accessKeyIDKey], opt["access_key_secret"])
	if err != nil {
		return nil, err
	}

	return &sender{
		client:       client,
		templateCode: opt[templateCodeKey],
		signName:     opt[signNameKey],
		proxyURL:     opt[proxyURLKey],
	}, nil
}

func validate(opt map[string]string) error {
	if _, exists := opt[accessKeyIDKey]; !exists {
		return errors.New("access_key_id can't be empty")
	}
	if _, exists := opt["access_key_secret"]; !exists {
		return errors.New("access_key_secret can't be empty")
	}
	if _, exists := opt[templateCodeKey]; !exists {
		return errors.New("template_code can't be empty")
	}
	if _, exists := opt[signNameKey]; !exists {
		return errors.New("sig_name can't be empty")
	}

	return nil
}
