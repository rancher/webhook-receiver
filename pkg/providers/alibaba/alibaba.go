package alibaba

import (
	"errors"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/dysmsapi"

	"github.com/rancher/receiver/pkg/providers"
)

const (
	scheme = "https"
	regionID = "region_id"
	accessKeyID = "access_key_id"
	accessKeySecret = "access_key_secret"
	templateCode = "template_code"
	Name = "alibaba"
)


type sender struct {
	client *dysmsapi.Client
}

func (s *sender) Send(msg string, receiver providers.Receiver) error {
	// TODO

	return nil

	request := dysmsapi.CreateSendSmsRequest()
	request.Scheme = scheme
	//request.PhoneNumbers = strings.Join(receiver.To, ",")
	//request.TemplateCode = receiver.TemplateCode

	_, err := s.client.SendSms(request)
	return err
}

func New(opt map[string]string) (providers.Sender, error) {
	if err := validate(opt); err != nil {
		return nil, err
	}

	client, err := dysmsapi.NewClientWithAccessKey(opt[regionID], opt[accessKeyID], opt[accessKeySecret])
	if err != nil {
		return nil, err
	}
	return &sender{client:  client}, nil
}

func validate(opt map[string]string) error {
	if _, exists := opt[regionID]; !exists {
		return errors.New("region_id can't be empty")
	}
	if _, exists := opt[accessKeyID]; !exists {
		return errors.New("access_key_id can't be empty")
	}
	if _, exists := opt[accessKeySecret]; !exists {
		return errors.New("access_key_secret can't be empty")
	}
	if _, exists := opt[templateCode]; !exists {
		return errors.New("template_code can't be empty")
	}

	return nil
}
