package options

import (
	"errors"
	"fmt"
	"io/ioutil"
	"sync"

	"github.com/fsnotify/fsnotify"
	"gopkg.in/yaml.v2"

	"github.com/rancher/receiver/pkg/providers"
	"github.com/rancher/receiver/pkg/providers/alibaba"
	"github.com/rancher/receiver/pkg/providers/dingtalk"
	log "github.com/sirupsen/logrus"
)

var (
	mut       sync.RWMutex
	receivers map[string]providers.Receiver
	senders   map[string]providers.Sender
)

// when occur error, it will panic directly
func Init(configPath string) {
	updateMemoryConfig(configPath)
	go syncMemoryConfig(configPath)
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

func syncMemoryConfig(configPath string) {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal("new config fs watcher err:", err)
	}
	if err := watcher.Add(configPath); err != nil {
		log.Fatalf("watch file:%s err:%v", configPath, err)
	}

	go func() {
		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					break
				}
				if event.Op&fsnotify.Write == fsnotify.Write {
					updateMemoryConfig(configPath)
				}
			case err, _ := <-watcher.Errors:
				log.Errorf("fs watcher event err:%v", err)
			}
		}
	}()
}

func updateMemoryConfig(configPath string) {
	opt, err := newOption(configPath)
	if err != nil {
		// log error, and not update
		log.Fatalf("update config on path:%s, err:%v, will not update config", configPath, err)
		return
	}

	updateReceivers := make(map[string]providers.Receiver)
	for _, v := range opt.Receivers {
		updateReceivers[v.Name] = v
	}

	updateSenders := make(map[string]providers.Sender)
	for k, v := range opt.Providers {
		creator, err := getProviderCreator(k)
		if err != nil {
			log.Errorf("update config err:%v", err)
			return
		}
		sender, err := creator(v)
		if err != nil {
			log.Errorf("update config err:%v", err)
			return
		}
		updateSenders[k] = sender
	}

	// replace
	mut.Lock()
	defer mut.Unlock()
	receivers = updateReceivers
	senders = updateSenders
	log.Info("update config sucess")
}

func getProviderCreator(name string) (providers.Creator, error) {
	switch name {
	case alibaba.Name:
		return alibaba.New, nil
	case dingtalk.Name:
		return dingtalk.New, nil
	default:
		return nil, errors.New(fmt.Sprintf("provider %s is not support", name))
	}
}

// for yaml parse
type option struct {
	Providers map[string]map[string]string
	Receivers []providers.Receiver
}

func newOption(configPath string) (option, error) {
	data, err := ioutil.ReadFile(configPath)
	if err != nil {
		return option{}, err
	}
	opt := option{}
	if err := yaml.Unmarshal(data, &opt); err != nil {
		return option{}, err
	}

	return opt, nil
}
