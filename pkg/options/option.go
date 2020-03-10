package options

import (
	"errors"
	"fmt"
	"path/filepath"
	"strings"
	"sync"

	"github.com/fsnotify/fsnotify"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"gopkg.in/yaml.v2"

	"github.com/rancher/webhook-receiver/pkg/providers"
	"github.com/rancher/webhook-receiver/pkg/providers/aliyunsms"
	"github.com/rancher/webhook-receiver/pkg/providers/dingtalk"
	"github.com/rancher/webhook-receiver/pkg/providers/msteams"
)

const logLevelErr = "set log level error, support Info,Error"

var (
	mut       sync.RWMutex
	receivers map[string]providers.Receiver
	senders   map[string]providers.Sender
	state     bool
)

func Init(configPath string) {
	dir := filepath.Dir(configPath)
	name := filepath.Base(configPath)
	viperConfigName := strings.TrimRight(name, ".yaml")
	viper.AddConfigPath(dir)
	viper.SetConfigName(viperConfigName)
	viper.SetConfigType("yaml")

	updateMemoryConfig()
	viper.OnConfigChange(func(in fsnotify.Event) {
		updateMemoryConfig()
	})

	go viper.WatchConfig()
}

func GetReceiverAndSender(receiverName string) (providers.Receiver, providers.Sender, error) {
	mut.RLock()
	defer mut.RUnlock()

	receiver, exists := receivers[receiverName]
	if !exists {
		return providers.Receiver{}, nil, fmt.Errorf("error, receiver:%s is not exists\n", receiverName)
	}

	sender, exists := senders[receiver.Provider]
	if !exists {
		return providers.Receiver{}, nil, fmt.Errorf("error, provider:%s is not exists\n", receiver.Provider)
	}

	return receiver, sender, nil
}

func updateMemoryConfig() {
	log.Info("update memory config")
	if err := viper.ReadInConfig(); err != nil {
		log.Errorf("read config err:%v", err)
		setStatus(false)
		return
	}

	receiversMap := viper.GetStringMap("receivers")
	updateReceivers := make(map[string]providers.Receiver)
	for k, v := range receiversMap {
		receiver := providers.Receiver{}
		if err := convertInterfaceToStruct(v, &receiver); err != nil {
			log.Errorf("parse receiver:%s to struct err:%v", k, err)
			setStatus(false)
			return
		}
		updateReceivers[k] = receiver
	}

	providersMap := viper.GetStringMap("providers")
	updateSenders := make(map[string]providers.Sender)
	for k, v := range providersMap {
		optMap := make(map[string]string)
		if err := convertInterfaceToStruct(v, &optMap); err != nil {
			log.Errorf("parse provider:%s err:%v", k, err)
			setStatus(false)
			return
		}

		if _, exists := optMap["type"]; !exists {
			log.Errorf("no required type field defined for %s", k)
			return
		}

		creator, err := getProviderCreator(optMap["type"])
		if err != nil {
			log.Errorf("update config err:%v", err)
			setStatus(false)
			return
		}

		sender, err := creator(optMap)
		if err != nil {
			log.Errorf("update config err:%v", err)
			setStatus(false)
			return
		}
		updateSenders[k] = sender
	}

	logLevel := viper.Get("logLevel")
	ll, ok := logLevel.(string)
	if ok {
		switch ll {
		case "Info":
			log.SetLevel(log.InfoLevel)
			log.Info("set log level info")
		case "Error":
			log.SetLevel(log.ErrorLevel)
			log.Info("set log level error")
		default:
			log.SetLevel(log.InfoLevel)
		}
	} else {
		log.SetLevel(log.InfoLevel)
	}

	setStatus(true)
	mut.Lock()
	defer mut.Unlock()
	receivers = updateReceivers
	senders = updateSenders
	log.Info("update config success")
}

func getProviderCreator(name string) (providers.Creator, error) {
	switch name {
	case aliyunsms.Name:
		return aliyunsms.New, nil
	case dingtalk.Name:
		return dingtalk.New, nil
	case msteams.Name:
		return msteams.New, nil
	default:
		return nil, errors.New(fmt.Sprintf("provider %s is not support", name))
	}
}

type option struct {
	Providers map[string]map[string]string
	Receivers []providers.Receiver
}

func newOption(data []byte) (option, error) {
	opt := option{}
	if err := yaml.Unmarshal(data, &opt); err != nil {
		return option{}, err
	}

	return opt, nil
}

func convertInterfaceToStruct(inter interface{}, s interface{}) error {
	byteData, err := yaml.Marshal(inter)
	if err != nil {
		return err
	}
	if err := yaml.Unmarshal(byteData, s); err != nil {
		return err
	}

	return nil
}

func GetState() bool {
	return state
}

func setStatus(now bool) {
	state = now
}
